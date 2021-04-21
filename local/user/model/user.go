package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	Normal   = 0 //正常状态
	Blocked  = 1 //封号状态
	Canceled = 2 //注销状态
)

/*
DB Table Details
用户基本信息表

CREATE TABLE `users` (
  `user_id` varchar(50) NOT NULL COMMENT '用户ID',
  `phone` varchar(30) NOT NULL COMMENT '手机号',
  `password` varchar(30) NOT NULL COMMENT '密码',
  `nickname` varchar(30) NOT NULL COMMENT '昵称',
  `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '注册时间',
  `status` int(11) NOT NULL COMMENT '用户状态',
  `deleted_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '删除时间',
  `avatar` varchar(100) NOT NULL COMMENT '用户头像',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `unique_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

*/

// User struct is a row record of the users table in the  database
type User struct {
	UserID    uuid.UUID  `gorm:"column:user_id;type:varchar(50);primary_key" json:"user_id"` //[ 0] user_id                                        VARCHAR[30]          null: false  primary: true   auto: false
	UserName  string     `gorm:"column:username;type:varchar(50);unique_index:unique_username" json:"userName"`
	Password  string     `gorm:"column:password;type:varchar(100);" json:"password"` //[ 2] password                                       VARCHAR[30]          null: false  primary: false  auto: false
	Status    int        `gorm:"column:status;type:INT;" json:"status"`              //[ 5] status                                         INT                  null: false  primary: false  auto: false
	CreatedAt time.Time  `gorm:"column:created_at;type:DATETIME;" json:"created_at"` //[ 4] created_at                                     DATETIME             null: false  primary: false  auto: false
	UpdatedAt time.Time  `gorm:"column:updated_at;type:DATETIME;" json:"updated_at"` //[16] updated_at                                     DATETIME             strue   primary: false  auto: false
	DeletedAt *time.Time `gorm:"column:deleted_at;type:DATETIME;" json:"deleted_at"` //[ 6] deleted_at                                     DATETIME             null: true   primary: false  auto: false
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "users"
}
