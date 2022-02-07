package main

import (
	"go-hexa-arch/config"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/gorm"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/mysql"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/role"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/user"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_CLI_MIGRATE"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := logRotator.New(
		"resources/logs/migrate/%Y%m%d%H%M.log",
		logRotator.WithLinkName("resources/logs/migrate/.log"),
		logRotator.WithMaxAge(time.Duration(86400)*time.Second),
		logRotator.WithRotationTime(time.Duration(86400)*time.Second),
	)
	if err == nil {
		logrus.SetOutput(logFile)
	} else {
		log.Error(err)
		panic(err)
	}
	defer logFile.Close()
	dbCon, err := mysql.Connect(config.GetConfig().MariaDbMasterDatabase)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	gormDb, err := gorm.InitGorm(dbCon)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	role.Migrate(*gormDb)
	user.Migrate(*gormDb)

	dbCon.Close()
}
