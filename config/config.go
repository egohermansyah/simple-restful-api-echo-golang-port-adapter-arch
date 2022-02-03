package config

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MariaDBConfig struct {
	Driver   string `mapstructure:"Driver"`
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Name     string `mapstructure:"Name"`
}

type MongoDbConfig struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Name     string `mapstructure:"Name"`
}

type AppConfig struct {
	ProductionMode        bool          `mapstructure:"ProductionMode"`
	MariaDbMasterDatabase MariaDBConfig `mapstructure:"MariaDBMasterDatabase"`
	MongoDBMasterDatabase MongoDbConfig `mapstructure:"MongoDBMasterDatabase"`
}

func SetConfig() *AppConfig {
	var appConfig AppConfig
	viper.BindEnv("ProductionMode", "PRODUCTION_MODE")
	viper.BindEnv("MariaDbMasterDatabase.Driver", "MASTER_DATABASE_DRIVER")
	viper.BindEnv("MariaDbMasterDatabase.Host", "MASTER_DATABASE_HOST")
	viper.BindEnv("MariaDbMasterDatabase.Port", "MASTER_DATABASE_PORT")
	viper.BindEnv("MariaDbMasterDatabase.Username", "MASTER_DATABASE_USERNAME")
	viper.BindEnv("MariaDbMasterDatabase.Password", "MASTER_DATABASE_PASSWORD")
	viper.BindEnv("MariaDbMasterDatabase.Name", "MASTER_DATABASE_NAME")
	viper.BindEnv("MongoDBMasterDatabase.Host", "MASTER_MONGODB_HOST")
	viper.BindEnv("MongoDBMasterDatabase.Port", "MASTER_MONGODB_PORT")
	viper.BindEnv("MongoDBMasterDatabase.Username", "MASTER_MONGODB_USERNAME")
	viper.BindEnv("MongoDBMasterDatabase.Password", "MASTER_MONGODB_PASSWORD")
	viper.BindEnv("MongoDBMasterDatabase.Name", "MASTER_MONGODB_NAME")

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Info(err)
	}
	return &appConfig
}

var appConfig *AppConfig
var lock = &sync.Mutex{}

func GetConfig() *AppConfig {
	if appConfig != nil {
		return appConfig
	}

	lock.Lock()
	defer lock.Unlock()
	if appConfig != nil {
		return appConfig
	}

	appConfig = SetConfig()

	return appConfig
}
