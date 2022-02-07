package main

import (
	"context"
	"fmt"
	"go-hexa-arch/config"
	"os"
	"os/signal"
	"simple-restful-api-echo-golang-port-adapter-archcmd/api/router"
	roleService "simple-restful-api-echo-golang-port-adapter-archcore/services/role"
	userService "simple-restful-api-echo-golang-port-adapter-archcore/services/user"
	roleHandler "simple-restful-api-echo-golang-port-adapter-archhandlers/role"
	userHandler "simple-restful-api-echo-golang-port-adapter-archhandlers/user"
	customValidator "simple-restful-api-echo-golang-port-adapter-archhandlers/util/validator"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/gorm"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/mongodb"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/mysql"
	roleRepository "simple-restful-api-echo-golang-port-adapter-archrepositories/role"
	userRepository "simple-restful-api-echo-golang-port-adapter-archrepositories/user"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
	"time"

	"github.com/go-playground/validator"

	"github.com/labstack/echo/v4"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_MAIN"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func CustomValidator() *validator.Validate {
	customValidator := validator.New()
	return customValidator
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := logRotator.New(
		"resources/logs/api/%Y%m%d%H%M.log",
		logRotator.WithLinkName("resources/logs/api/.log"),
		logRotator.WithMaxAge(time.Duration(3600)*time.Second),
		logRotator.WithRotationTime(time.Duration(3600)*time.Second),
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
		log.Fatal(err)
		panic(err)
	}

	gormDb, err := gorm.InitGorm(dbCon)
	if err != nil {
		dbCon.Close()
		log.Fatal(err)
		panic(err)
	}

	mongoDb, err := mongodb.Connect(config.GetConfig().MongoDBMasterDatabase)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Print(mongoDb)

	role, err := roleRepository.New(gormDb)
	if err != nil {
		log.Error(err)
	}
	user, err := userRepository.New(gormDb)
	if err != nil {
		log.Error(err)
	}
	roleService := roleService.New(role)
	roleController := roleHandler.New(roleService)
	userService := userService.New(user)
	userController := userHandler.New(userService)

	controllers := router.Handler{
		Role: *roleController,
		User: *userController,
	}

	e := echo.New()
	e.Validator = &customValidator.BodyRequestValidator{Validator: CustomValidator()}

	router.RegisterPath(e, controllers)

	go func() {
		if err := e.Start(":2022"); err != nil {
			dbCon.Close()
			log.Fatal(err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		dbCon.Close()
		log.Fatal(err)
		panic(err)
	}
}
