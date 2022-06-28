package model

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Bucket *oss.Bucket
var err error
var Config *ini.File

//连接MySQL和OSS
func Init() error {
	//读取.ini里面的数据库配置
	var iniErr error
	// Config, iniErr = ini.Load("./douyin/model/app.ini")
	//测试使用
	Config, iniErr = ini.Load("/home/xueqin/Qinxun2/src/douyin/model/app.ini")
	if iniErr != nil {
		return iniErr
	}

	ip := Config.Section("mysql").Key("ip").String()
	port := Config.Section("mysql").Key("port").String()
	user := Config.Section("mysql").Key("user").String()
	password := Config.Section("mysql").Key("password").String()
	database := Config.Section("mysql").Key("database").String()

	// 打开数据库
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, password, ip, port, database)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, //打印sql
		//SkipDefaultTransaction: true, //禁用事务
	})
	// DB.Debug()
	if err != nil {
		return err
	}

	//打开oss的Bucket
	endpoint := Config.Section("oss").Key("endpoint").String()
	accesskeyid := Config.Section("oss").Key("accessKeyId").String()
	accessKeySecret := Config.Section("oss").Key("accessKeySecret").String()
	bucket := Config.Section("oss").Key("bucket").String()
	client, err := oss.New(endpoint, accesskeyid, accessKeySecret)
	if err != nil {
		return err
	}
	Bucket, err = client.Bucket(bucket)
	if err != nil {
		return err
	}
	return nil
}
