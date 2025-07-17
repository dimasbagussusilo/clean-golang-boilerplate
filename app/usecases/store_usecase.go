package usecases

import (
	"appsku-golang/app/models"
	"appsku-golang/app/repositories"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IStoreUseCase interface {
	Insert(ctx context.Context, store *models.Store) (*models.Store, error)
	InsertWithSetting(ctx context.Context, store *models.Store, setting *models.StoreSetting) (*models.Store, *models.StoreSetting, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*models.Store, error)
	GetAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Store, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, store map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID, hardDelete bool) error
}

type StoreUseCase struct {
	StoreRepository repositories.IStoreRepository
}

func NewStoreUseCase(StoreRepository repositories.IStoreRepository) IStoreUseCase {
	return &StoreUseCase{
		StoreRepository: StoreRepository,
	}
}

func (u *StoreUseCase) Insert(ctx context.Context, store *models.Store) (*models.Store, error) {
	now := time.Now()
	store.CreatedAt = &now

	storedStore, err := u.StoreRepository.Insert(ctx, store)
	if err != nil {
		return nil, err
	}

	objectID, _ := storedStore.InsertedID.(primitive.ObjectID)
	store.ID = objectID

	return store, nil
}

func (u *StoreUseCase) GetById(ctx context.Context, id primitive.ObjectID) (*models.Store, error) {
	return u.StoreRepository.GetById(ctx, id)
}

func (u *StoreUseCase) GetAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Store, int64, error) {
	return u.StoreRepository.GetAll(ctx, filter, page, limit)
}

func (u *StoreUseCase) Update(ctx context.Context, id primitive.ObjectID, store map[string]interface{}) error {
	return u.StoreRepository.Update(ctx, id, store)
}

func (u *StoreUseCase) Delete(ctx context.Context, id primitive.ObjectID, hardDelete bool) error {
	return u.StoreRepository.Delete(ctx, id, hardDelete)
}

func (u *StoreUseCase) InsertWithSetting(ctx context.Context, store *models.Store, setting *models.StoreSetting) (*models.Store, *models.StoreSetting, error) {
	now := time.Now()
	store.CreatedAt = &now
	setting.CreatedAt = &now

	storedStore, storedSetting, err := u.StoreRepository.InsertWithSetting(ctx, store, setting)
	if err != nil {
		return nil, nil, err
	}

	storeObjectID, _ := storedStore.InsertedID.(primitive.ObjectID)
	store.ID = storeObjectID

	settingObjectID, _ := storedSetting.InsertedID.(primitive.ObjectID)
	setting.ID = settingObjectID

	return store, setting, nil
}
