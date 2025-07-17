package repositories

import (
	"appsku-golang/app/global-utils/mongodb"
	"appsku-golang/app/global-utils/redisdb"
	"fmt"
)

type IExampleRepository interface {
	Get()
}

type ExampleRepository struct {
	MongoDB mongodb.IMongoDB
	Redis   redisdb.IRedis
}

func NewExampleRepository(mongodb mongodb.IMongoDB, redisdb redisdb.IRedis) IExampleRepository {
	return &ExampleRepository{
		MongoDB: mongodb,
		Redis:   redisdb,
	}
}

func (r *ExampleRepository) Get() {
	fmt.Println("Rawr")
}
