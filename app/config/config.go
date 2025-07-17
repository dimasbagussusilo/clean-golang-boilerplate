package config

import (
	"appsku-golang/app/constants"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"

	"appsku-golang/app/global-utils/grpcclient"
	"appsku-golang/app/global-utils/helper"
	"appsku-golang/app/global-utils/mongodb"
	"appsku-golang/app/global-utils/redisdb"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Environment    string
	UseSSL         bool
	MainPort       int
	LogLevel       string
	PrivateSSLPath string
	PublicSSLPath  string
	EncryptKey     string
	EncryptIV      string
	WebFramework   string
	Database       struct {
		Mongo struct {
			Host        string
			Port        int
			User        string
			Password    string
			Database    string
			UsedReplica int
		}
	}
	Cache struct {
		Redis struct {
			Host     string
			Port     int
			Password string
			Database int
		}
		DefaultCacheTime int
	}
	MessageBroker struct {
		Kafka struct {
			Hosts []string
		}
	}
	AuthBasic struct {
		Username string
		Password string
	}
	Grpc struct {
		Client struct {
			ExampleService struct {
				Host string
				Port int
			}
			IsUseSSL bool
		}
		Server struct {
			Host string
			Port int
		}
	}
}

var cfg Configuration

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		errStr := fmt.Sprintf(".env not load properly %s", err.Error())
		log.Fatal(errStr)
		helper.SetSentryError(err, errStr, sentry.LevelError)
		return
	}

	// server getting config
	cfg.Environment = GetEnvString("ENVIRONMENT", "development")
	cfg.UseSSL = GetEnvBool("USE_SSL", false)
	cfg.MainPort = GetEnvInt("MAIN_PORT", 8000)
	cfg.LogLevel = GetEnvString("LOG_LEVEL", "debug")
	cfg.PrivateSSLPath = GetEnvString("PRIVATE_SSL_PATH", "")
	cfg.PublicSSLPath = GetEnvString("PUBLIC_SSL_PATH", "")
	cfg.WebFramework = GetEnvString("WEB_FRAMEWORK", "gin")

	// encryption getting config
	cfg.EncryptKey = GetEnvString("ENCRYPT_KEY", "")
	cfg.EncryptIV = GetEnvString("ENCRYPT_IV", "")

	// mongo getting config
	cfg.Database.Mongo.Host = GetEnvString("MONGO_HOST", "")
	cfg.Database.Mongo.Port = GetEnvInt("MONGO_PORT", 0)
	cfg.Database.Mongo.User = GetEnvString("MONGO_USER", "")
	cfg.Database.Mongo.Password = GetEnvString("MONGO_PASSWORD", "")
	cfg.Database.Mongo.Database = GetEnvString("MONGO_DATABASE", "")
	cfg.Database.Mongo.UsedReplica = GetEnvInt("MONGO_USED_REPLICA", 0)

	// redis getting config
	cfg.Cache.Redis.Host = GetEnvString("REDIS_HOST", "")
	cfg.Cache.Redis.Port = GetEnvInt("REDIS_PORT", 0)
	cfg.Cache.Redis.Password = GetEnvString("REDIS_PASSWORD", "")
	cfg.Cache.Redis.Database = GetEnvInt("REDIS_DATABASE", 0)
	cfg.Cache.DefaultCacheTime = GetEnvInt("DEFAULT_CACHE_TIME", 10)

	// kafka getting config
	cfg.MessageBroker.Kafka.Hosts = GetEnvSliceString("KAFKA_HOSTS", "")

	// basic auth getting config
	cfg.AuthBasic.Username = GetEnvString("AUTHBASIC_USERNAME", "")
	cfg.AuthBasic.Password = GetEnvString("AUTHBASIC_PASSWORD", "")

	// grpc getting config
	cfg.Grpc.Server.Host = GetEnvString("GRPC_SERVER_HOST", "")
	cfg.Grpc.Server.Port = GetEnvInt("GRPC_SERVER_PORT", 0)
	cfg.Grpc.Client.ExampleService.Host = GetEnvString("GRPC_STORE_SERVICE_HOST", "")
	cfg.Grpc.Client.ExampleService.Port = GetEnvInt("GRPC_STORE_SERVICE_PORT", 0)
	cfg.Grpc.Client.IsUseSSL = GetEnvBool("GRPC_CLIENT_IS_USE_SSL", false)
}

func Get() Configuration {
	return cfg
}

func GetEnvString(key string, dflt string) string {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	return value
}

func GetEnvSliceString(key string, dflt string) []string {
	values := os.Getenv(key)
	if values == "" {
		return strings.Split(dflt, ",")
	}
	return strings.Split(values, ",")
}

func GetEnvInt(key string, dflt int) int {
	value := os.Getenv(key)
	i, err := strconv.ParseInt(value, 10, 64)
	if value == "" && err != nil {
		return dflt
	}
	return int(i)
}

func GetEnvBool(key string, dflt bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return dflt
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return dflt
	}
	return b
}

func BuildGrpcClientParam(name string) grpcclient.GRPCClientParam {
	var (
		host string
		port string
	)

	switch name {
	case constants.ExampleService:
		host = cfg.Grpc.Client.ExampleService.Host
		port = strconv.Itoa(cfg.Grpc.Client.ExampleService.Port)
	}

	return grpcclient.GRPCClientParam{
		Name:      name,
		Host:      host,
		Port:      port,
		IsUseSSL:  cfg.Grpc.Client.IsUseSSL,
		ProxyPath: cfg.PublicSSLPath,
	}
}

func BuildRedisParam() redisdb.RedisParam {
	return redisdb.RedisParam{
		Host:     cfg.Cache.Redis.Host,
		Port:     cfg.Cache.Redis.Port,
		Password: cfg.Cache.Redis.Password,
		Database: cfg.Cache.Redis.Database,
	}
}

func BuildKafkaParam() []string {
	return cfg.MessageBroker.Kafka.Hosts
}

func BuildMongoDBParam() mongodb.MongoDBParam {
	return mongodb.MongoDBParam{
		Host:             cfg.Database.Mongo.Host,
		Database:         cfg.Database.Mongo.Database,
		Port:             cfg.Database.Mongo.Port,
		User:             cfg.Database.Mongo.User,
		Password:         cfg.Database.Mongo.Password,
		UsedMongoReplica: cfg.Database.Mongo.UsedReplica,
	}
}
