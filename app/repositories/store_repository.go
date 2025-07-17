package repositories

import (
	"appsku-golang/app/config"
	"appsku-golang/app/constants"
	"appsku-golang/app/global-utils/mongodb"
	"appsku-golang/app/global-utils/redisdb"
	"appsku-golang/app/models"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type IStoreRepository interface {
	Insert(ctx context.Context, store *models.Store) (*mongo.InsertOneResult, error)
	InsertWithSetting(ctx context.Context, store *models.Store, setting *models.StoreSetting) (*mongo.InsertOneResult, *mongo.InsertOneResult, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Store, error)
	GetAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Store, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, store map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID, hardDelete bool) error
}

type StoreRepository struct {
	MongoDB mongodb.IMongoDB
	Redis   redisdb.IRedis
}

func NewStoreRepository(mongodb mongodb.IMongoDB, redisdb redisdb.IRedis) IStoreRepository {
	return &StoreRepository{
		MongoDB: mongodb,
		Redis:   redisdb,
	}
}

func (r *StoreRepository) Get() {
	fmt.Println("Rawr")
}

func (r *StoreRepository) Insert(ctx context.Context, store *models.Store) (*mongo.InsertOneResult, error) {
	var (
		cfg        = config.Get()
		collection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
	)

	ID, err := collection.InsertOne(ctx, store)
	if err != nil {
		return nil, err
	}
	return ID, nil
}

func (r *StoreRepository) GetById(ctx context.Context, id primitive.ObjectID) (*models.Store, error) {
	var (
		store      *models.Store
		cfg        = config.Get()
		collection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
	)

	redisKey := fmt.Sprintf("%s_%s", constants.REDIS_KEY_STORE_BY_ID, id.Hex())
	existInRedis, err := r.Redis.Client().Exists(ctx, redisKey).Result()
	if err != nil {
		return nil, err
	}

	if existInRedis == 1 {
		storeDataInRedis, err := r.Redis.Client().Get(ctx, redisKey).Result()
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(storeDataInRedis), &store)
		if err != nil {
			return nil, err
		}

		return store, nil
	}

	filter := bson.M{"_id": id, "deleted_at": nil}

	err = collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		return nil, err
	}

	storeJson, err := json.Marshal(store)
	if err != nil {
		return nil, err
	}

	err = r.Redis.Client().Set(ctx, redisKey, string(storeJson), constants.REDIS_TTL_STORE).Err()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r *StoreRepository) GetAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Store, int64, error) {
	var (
		stores     []*models.Store
		cfg        = config.Get()
		collection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
	)

	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}

	bsonFilter["deleted_at"] = nil

	skip := (page - 1) * limit

	totalCount, err := collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &stores); err != nil {
		return nil, 0, err
	}

	return stores, totalCount, nil
}

func (r *StoreRepository) Update(ctx context.Context, id primitive.ObjectID, storeData map[string]interface{}) error {
	var (
		cfg        = config.Get()
		collection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
	)

	now := time.Now()
	storeData["updated_at"] = &now

	update := bson.M{
		"$set": storeData,
	}

	filter := bson.M{"_id": id, "deleted_at": nil}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("%s_%s", constants.REDIS_KEY_STORE_BY_ID, id.Hex())
	r.Redis.Client().Del(ctx, redisKey)

	return nil
}

func (r *StoreRepository) Delete(ctx context.Context, id primitive.ObjectID, hardDelete bool) error {
	var (
		cfg        = config.Get()
		collection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
	)

	if hardDelete {
		filter := bson.M{"_id": id}
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			return err
		}
	} else {
		now := time.Now()
		update := bson.M{
			"$set": bson.M{
				"deleted_at": &now,
			},
		}
		filter := bson.M{"_id": id, "deleted_at": nil}
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	redisKey := fmt.Sprintf("%s_%s", constants.REDIS_KEY_STORE_BY_ID, id.Hex())
	r.Redis.Client().Del(ctx, redisKey)

	return nil
}

func (r *StoreRepository) InsertWithSetting(ctx context.Context, store *models.Store, setting *models.StoreSetting) (*mongo.InsertOneResult, *mongo.InsertOneResult, error) {
	var (
		cfg               = config.Get()
		storeCollection   = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORES)
		settingCollection = r.MongoDB.Client().Database(cfg.Database.Mongo.Database).Collection(constants.MONGO_COL_STORE_SETTINGS)
	)

	session, err := r.MongoDB.Client().StartSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.EndSession(ctx)

	var storeResult *mongo.InsertOneResult
	var settingResult *mongo.InsertOneResult
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}

		storeResult, err = storeCollection.InsertOne(sessionContext, store)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		storeID, ok := storeResult.InsertedID.(primitive.ObjectID)
		if !ok {
			session.AbortTransaction(sessionContext)
			return fmt.Errorf("failed to get store ID")
		}

		setting.StoreID = storeID
		settingResult, err = settingCollection.InsertOne(sessionContext, setting)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		err = session.CommitTransaction(sessionContext)
		if err != nil {
			session.AbortTransaction(sessionContext)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return storeResult, settingResult, nil
}
