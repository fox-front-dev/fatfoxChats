package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"

	viper "github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("fatfoxChats/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func InitMysql() {
	//自定义日志模版
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	DB, _ = gorm.Open(mysql.Open("fox:130561Cr@tcp(81.68.206.160:3306)/fatfox?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		//DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: newLogger,
	},
	)
}
func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		DB:           viper.GetInt("redis.DB"),
		Password:     viper.GetString("redis.password"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
}

const (
	PublishKey = "websocket"
)

//发布消息到redis

func Publish(c context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publicsh.....")
	err = Red.Publish(c, channel, msg).Err()
	return err
}

// 订阅消息
func Subscribe(c context.Context, channel string) (string, error) {
	sub := Red.Subscribe(c, channel)
	msg, err := sub.ReceiveMessage(c)
	fmt.Println("subscribe====>", msg.Payload)
	return msg.Payload, err
}
