package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User 定义User模型，绑定users表
type User struct {
	ID         int64  // primary key
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	CreateTime int64  `gorm:"column:createtime"`
}

func (u *User) TableName() string {
	return "users"
}

func main() {
	GinServer()
	ConnectDatabase()
}

func GinServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		return
	}
}

func ConnectDatabase() {
	// 配置MySQL连接参数
	username := "root"    //账号
	password := "root"    //密码
	host := "127.0.0.1"   //数据库地址，可以是Ip或者域名
	port := 3306          //数据库端口
	Dbname := "crud_demo" //数据库名

	//dsn := "root:root@tcp(127.0.0.1:3306)/crud_demo?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database, error = " + err.Error())
	}

	CreateTable(db)

	CRUD(db)
}

func CRUD(db *gorm.DB) {
	// 定义一个User，并初始化数据
	u := User{
		Username:   "lisi",
		Password:   "1234",
		CreateTime: time.Now().Unix(),
	}

	// 插入用户数据
	if err := db.Create(&u).Error; err != nil {
		fmt.Println("insert failed", err)
		return
	}

	// 查询一条数据
	u = User{}
	res := db.Where("username = ?", "zhangsan").First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("not found")
		return
	}
	fmt.Println(u.Username, u.Password)

	// 更新
	db.Model(&User{}).Where("username = ?", "zhangsan").Update("password", "4321")

	// 删除
	db.Where("username = ?", "zhangsan").Delete(&User{})
}

func CreateTable(db *gorm.DB) {
	// 自动建表
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to create table, error = " + err.Error())
	}
}
