package dao

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"
)

type User struct {
	Id         int32     `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Name       string    `json:"name" db:"name"`
	Password   string    `json:"password" db:"password"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

//AddUser 往数据库增加一名用户
func AddUser(ctx context.Context, email, name, password string) (userId int32, err error) {
	if email == "" || name == "" || password == "" {
		return -1, errors.New("email, name or password empty")
	}
	oUser := GetUserByEmail(ctx, email)
	if oUser.Id > 0 {
		return oUser.Id, errors.New("the user already exists")
	}
	if _, err = db.SqlDb.QueryContext(ctx, `insert into z_user(id, email, name, password, create_time) 
values(null, ?, ?, ?, ?)`, email, name, password, time.Now()); err != nil {
		return -1, err
	}
	u := GetUserByEmail(ctx, email)
	misc.Logger.Info("add user to sql success", zap.Any("user", u))
	return u.Id, nil
}

//GetUserByEmail 根据email查找用户
func GetUserByEmail(ctx context.Context, email string) User {
	var user User
	if err := db.SqlDb.GetContext(ctx, &user,
		`select * from z_user where email = ?`, email); err != nil {
		misc.Logger.Warn("GetUserByEmail err, no this user", zap.String("email", email))
		return User{}
	}
	return user
}
