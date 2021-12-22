package dao

import (
	"context"
	"errors"
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"
	"go.uber.org/zap"
	"log"
	"time"
)

type Goods struct {
	Id         int32     `json:"id" db:"id"`   // 物品id,可以用来找对应照片的路径
	UserId     int32     `json:"uid" db:"uid"` //上传该物品的用户id
	Name       string    `json:"name" db:"name"`
	Uname      string    `json:"uname" db:"uname"`
	Price      float64   `json:"price" db:"price"`
	Detail     string    `json:"detail" db:"detail"`
	Cover      string    `json:"cover" db:"cover"` //封面url
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

type PicGoods struct {
	Id      int32  `json:"id" db:"id"`           //唯一表示图片
	GoodsId int32  `json:"goodsId" db:"goodsId"` //对应物品的id
	Path    string `json:"path" db:"path"`       //图片的路径
}

//AddGoods 增加一条物品信息
func AddGoods(ctx context.Context, name, detail string, price float64, uid int32, coverPath string, filePathList []string) (goodsId int32, err error) {
	if name == "" || detail == "" {
		return -1, errors.New("name, price or password empty")
	}

	tx, err := db.SqlDb.Begin()
	if err != nil {
		return -1, err
	}

	var uname string
	row := tx.QueryRowContext(ctx, `select name from z_user where id = ?`, uid)
	if err = row.Scan(&uname); err != nil {
		_ = tx.Rollback()
		return -1, err
	}

	res, err := tx.ExecContext(ctx, `insert into z_goods(id, uid, name, uname, price, detail, cover, create_time) 
values(null, ?, ?, ?, ?, ?, ?, ?)`, uid, name, uname, price, detail, coverPath, time.Now())
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	goodsId64, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	goodsId = int32(goodsId64)

	log.Println("goodsId", goodsId)

	for _, path := range filePathList {
		if _, err = tx.ExecContext(ctx, `insert into z_goods_pic(id, uId, gId, path) values(null, ?, ?, ?)`,
			uid, goodsId, path); err != nil {
			_ = tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	misc.Logger.Info("add good to sql success", zap.Any("goodId", goodsId))
	return
}

func GetGoodsByUserId(ctx context.Context, userId int32) ([]Goods, error) {
	var goodsList []Goods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &goodsList, `select * from z_goods where uid = ?`,
		userId); err != nil {
		return nil, err
	}
	return goodsList, nil
}

func GetGoods(ctx context.Context, page, count int32) ([]Goods, error) {
	var goodsList []Goods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &goodsList, `select * from z_goods where id in 
(select t.id from (select id from z_goods limit ?, ?) as t)`, page, count); err != nil {
		misc.Logger.Warn("GetGoods err", zap.Error(err))
		return nil, err
	}
	return goodsList, nil
}

//GetGoodsPic 根据商品id获取其具体照片
func GetGoodsPic(ctx context.Context, gid int32) ([]PicGoods, error) {
	var picGoods []PicGoods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &picGoods, `select * from z_goods_pic where gid = ?`, gid); err != nil {
		misc.Logger.Warn("get picGoods err", zap.Error(err))
		return nil, err
	}
	return picGoods, nil
}
