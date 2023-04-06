package config

import (
	"fmt"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Feishu feishuapi.Config
	Server struct {
		Port int
	}
	Url struct {
		Url4AccessToken string
	}
	TemplateSpace struct {
		SpaceID         string
		InitNodeToken   string
		MinuteNodeToken string
	}
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		DBname   string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}

var C Config

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		logrus.Error("Failed to unmarshal config")
	}

	logrus.Info("Configuration file loaded")
}

func SetupFeishuApiClient(cli *feishuapi.AppClient) {
	cli.Conf = C.Feishu
}

func GetDatabaseLoginInfo() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		C.Mysql.User,
		C.Mysql.Password,
		C.Mysql.Host,
		C.Mysql.Port,
		C.Mysql.DBname,
	)
}
