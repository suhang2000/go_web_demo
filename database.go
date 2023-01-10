package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func ConnectDatabase() *gorm.DB {
	// 配置MySQL连接参数
	username := "root"    //账号
	password := "root"    //密码
	host := "127.0.0.1"   //数据库地址，可以是Ip或者域名
	port := 3306          //数据库端口
	Dbname := "crud_demo" //数据库名

	//dsn := "root:root@tcp(127.0.0.1:3306)/crud_demo?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	// 连接数据库
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 自动建表时不自动添加复数，如 user -> users
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect to database, error = " + err.Error())
	}
	return db
}

func ConnectionPool(db *gorm.DB) {
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, _ := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Second * 10)
}
