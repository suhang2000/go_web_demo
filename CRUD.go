package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_web_demo/DAO"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func main() {
	// 1. 数据库初始操作
	// 1.1 连接数据库
	db := ConnectDatabase()
	// 1.2 配置数据库，连接池
	ConnectionPool(db)
	// 1.3 建表
	_ = db.AutoMigrate(&DAO.Userinfo{})

	// 2. Gin Server
	r := gin.Default()

	// 3. CRUD

	// 3.1 查询
	// 3.1.1 分页查询
	/**
	@api	/user/list/pageSize=?/pageNum=?
	@param	pageSize 每页多少条数据
	@param	pageNum 第几页
	@return	分页数据
	*/
	r.GET("/user/list", func(c *gin.Context) {
		var userList []DAO.Userinfo
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))
		// 判断是否需要分页
		//if pageSize <= 0 {
		//	pageSize = -1
		//}
		//if pageNum <= 0 {
		//	pageNum = -1
		//}
		//// 处理分页
		//offset := (pageNum - 1) * pageSize
		//if pageNum == -1 && pageSize == -1 {
		//	offset = -1
		//}

		var offset int
		if pageSize > 0 && pageNum > 0 {
			offset = (pageNum - 1) * pageSize
		} else {
			offset = -1
		}

		var total int64

		// 查询数据库
		// Model(结构体).Count(总数)	查询总数
		// Offset(偏移量).Limit(数据大小).Find(查询且数据插入到哪)	返回pageNum条数据
		// 如果 Limit<=0 则不添加LIMIT，如果 Offset<0 则不添加OFFSET
		db.Model(userList).Count(&total).Limit(pageSize).Offset(offset).Find(&userList)
		// db.Limit(pageSize).Offset(offset).Find(&userList)

		if len(userList) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg":  "no data found",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": "find data successfully",
				"data": gin.H{
					"user":     userList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	// 3.1.2 普通条件查询
	/**
	@api	/user/list/:name
	@param	name
	*/
	r.GET("user/list/:name", func(c *gin.Context) {
		name := c.Param("name")
		var userList []DAO.Userinfo
		db.Where("name = ?", name).Find(&userList)

		if len(userList) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg":  "not found",
				"data": gin.H{},
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "find successfully",
				"data": userList,
			})
		}
	})

	// 3.2 添加数据
	/**
	@api	/user/add
	@param	&Userinfo
	*/
	r.POST("user/add", func(c *gin.Context) {
		var user DAO.Userinfo
		err := c.ShouldBindJSON(&user) // 解析 post body 中的json信息，并绑定到结构体中
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "add failed",
				"data": gin.H{},
			})
		} else {
			db.Create(&user)
			c.JSON(http.StatusOK, gin.H{
				"msg":  "add successfully",
				"data": user,
			})
		}
	})

	// 3.3 删除数据
	/**
	@api	/user/delete/:id
	@param	id
	*/
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var users []DAO.Userinfo
		id := c.Param("id")
		//db.First(&user, id)
		db.Where("id = ?", id).Find(&users)
		if len(users) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "not found",
			})
		} else {
			db.Where("id = ?", id).Delete(&users)
			c.JSON(http.StatusOK, gin.H{
				"msg": "delete successfully",
			})
		}
	})

	// 3.4 修改数据
	/**
	@api	/user/update/:id
	@param	id
	*/
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var user DAO.Userinfo
		id := c.Param("id")
		// db.Select("id").Where("id = ?", id).Find(&user)
		if err := db.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// not found
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "user not found",
			})
		} else {
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "update failed",
				})
			} else {
				// db.Where("id = ?", id).Updates(&user)
				db.Save(&user)
				c.JSON(200, gin.H{
					"msg": "update successfully",
				})
			}
		}
	})

	_ = r.Run()
}
