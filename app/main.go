package main

import (
	"context"
	"fmt"
	"golang-vscode-setup/config"
	router "golang-vscode-setup/controller"
	roleController "golang-vscode-setup/controller/role"
	userController "golang-vscode-setup/controller/user"
	customValidator "golang-vscode-setup/controller/util/validator"
	"golang-vscode-setup/repository/mongodb"
	"golang-vscode-setup/repository/mysql"
	roleRepository "golang-vscode-setup/repository/role"
	userRepository "golang-vscode-setup/repository/user"
	roleService "golang-vscode-setup/service/role"
	userService "golang-vscode-setup/service/user"
	"golang-vscode-setup/util/gorm"
	"golang-vscode-setup/util/logger"
	"os"
	"os/signal"
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
		"resource/log/log.%Y%m%d%H%M",
		logRotator.WithLinkName("resource/log/log"),
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

	role, err := roleRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}
	user, err := userRepository.NewRepository(gormDb)
	if err != nil {
		log.Error(err)
	}
	roleService := roleService.NewService(role)
	roleController := roleController.NewController(roleService)
	userService := userService.NewService(user)
	userController := userController.NewController(userService)

	controllers := router.Controllers{
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
