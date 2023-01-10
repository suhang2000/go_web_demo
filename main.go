package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_web_demo/DAO"
	"gorm.io/gorm"
	"time"
)

func main() {
	GinServer()
	OperateDatabase()
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

func OperateDatabase() {
	db := ConnectDatabase()

	CreateTable(db)

	CRUD(db)
}

func CRUD(db *gorm.DB) {
	// 定义一个User，并初始化数据
	u := DAO.User{
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
	u = DAO.User{}
	res := db.Where("username = ?", "zhangsan").First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("not found")
		return
	}
	fmt.Println(u.Username, u.Password)

	// 更新
	db.Model(&DAO.User{}).Where("username = ?", "zhangsan").Update("password", "4321")

	// 删除
	db.Where("username = ?", "zhangsan").Delete(&DAO.User{})
}

func CreateTable(db *gorm.DB) {
	// 自动建表
	err := db.AutoMigrate(&DAO.User{})
	if err != nil {
		panic("failed to create table, error = " + err.Error())
	}
}
