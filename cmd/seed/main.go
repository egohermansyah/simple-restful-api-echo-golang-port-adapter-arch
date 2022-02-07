package main

import (
	"go-hexa-arch/config"
	"simple-restful-api-echo-golang-port-adapter-archcmd/seed/role"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/gorm"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/mysql"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_CLI_SEED"

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
		"resources/logs/seed/%Y%m%d%H%M.log",
		logRotator.WithLinkName("resources/logs/seed/.log"),
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

	role.Seed(gormDb)

	dbCon.Close()
}
