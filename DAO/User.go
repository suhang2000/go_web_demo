package DAO

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
