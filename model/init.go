package model

import (
	"LarkVCBot/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var gormDb *gorm.DB
var redisClient redis.UniversalClient

func Connect(dialector gorm.Dialector) {
	var err error
	gormDb, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	if gormDb == nil {
		logrus.Fatal("DB is nil")
	}

	logrus.Info("MySQL connected")
}

func CreateTables() {
	if gormDb == nil {
		logrus.Fatal("DB is nil")
	}
	err := gormDb.AutoMigrate(&GroupSpace{})
	if err != nil {
		logrus.Fatal(err)
	}
}

func getRedisLoginURL() string {
	loginInfo := config.C.Redis

	return fmt.Sprintf("%s:%d", loginInfo.Host, loginInfo.Port)
}

func ConnectRedis() {
	addr := make([]string, 0)
	addr = append(addr, getRedisLoginURL())
	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addr,
		DB:       config.C.Redis.DB,
		Password: config.C.Redis.Password,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatal(err)
		return
	}

	logrus.Info("Redis connected")
}
