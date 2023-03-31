package main

import (
	"LarkVCBot/app"
	"LarkVCBot/config"
	"LarkVCBot/docs"
	"LarkVCBot/global"
	"LarkVCBot/model"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
)

func main() {
	config.ReadConfig()

	// log
	config.SetupLogrus()
	logrus.Info("Robot starts up")

	// feishu api client
	config.SetupFeishuApiClient(&global.FeishuClient)
	global.FeishuClient.StartTokenTimer()

	// database
	model.Connect(mysql.Open(config.GetDatabaseLoginInfo()))
	model.CreateTables()

	// robot server
	r := gin.Default()
	app.Init(r)

	// api docs by swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(":" + fmt.Sprint(config.C.Server.Port))
}
