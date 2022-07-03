package model

import (
	"douyin/util/logger"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Bucket *oss.Bucket
var err error
var Config *ini.File
var Path string

//连接MySQL和OSS
func ConfigInit() error {
	//算出绝对路径，防止service层测试时路径错误
	dir, err := os.Getwd()
	if err != nil {
		logger.Error("Getwd error, %v", err.Error())
		return errors.New("Getwd")
	}
	Path = strings.Split(dir, "/service")[0]
	//读取.ini里面的数据库配置
	Config, err = ini.Load(Path + "/model/app.ini")
	if err != nil {
		panic(err.Error())
		// logger.Error("load ini config fail: ", err)
		// return errors.New("ini")
	}
	return nil
}

func MysqlInit() error {

// 	ip := Config.Section("mysql").Key("ip").String()
	port := Config.Section("mysql").Key("port").String()
	user := Config.Section("mysql").Key("user").String()
	password := Config.Section("mysql").Key("password").String()
	database := Config.Section("mysql").Key("database").String()

	// 打开数据库
	dsn := fmt.Sprintf("%v:%v@tcp(%:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, port, database)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})
	// DB.Debug()
	if err != nil {
		logger.Error("open mysql fail ", err)
		return err
	}
	return nil
}

func OssInit() error {
	//打开oss的Bucket
	endpoint := Config.Section("oss").Key("endpoint").String()
	accesskeyid := Config.Section("oss").Key("accessKeyId").String()
	accessKeySecret := Config.Section("oss").Key("accessKeySecret").String()
	bucket := Config.Section("oss").Key("bucket").String()
	client, err := oss.New(endpoint, accesskeyid, accessKeySecret)
	if err != nil {
		logger.Error("create oss client fail ", err)
		return err
	}
	Bucket, err = client.Bucket(bucket)
	if err != nil {
		logger.Error("instance bucket fail ", err)
		return err
	}
	logger.Info("open mysql and oss success")
	return nil
}
