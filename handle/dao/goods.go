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
	Id     int32   `json:"id" db:"id"`   // 物品id,可以用来找对应照片的路径
	UserId int32   `json:"uid" db:"uid"` //上传该物品的用户id
	Name   string  `json:"name" db:"name"`
	Price  float64 `json:"price" db:"price"`
	Detail string  `json:"detail" db:"detail"`
}

type PicGoods struct {
	Id      int32  `json:"id" db:"id"`           //唯一表示图片
	GoodsId int32  `json:"goodsId" db:"goodsId"` //对应物品的id
	Path    string `json:"path" db:"path"`       //图片的路径
}

//AddGoods 增加一条物品信息
func AddGoods(ctx context.Context, name, detail string, price float64, uid int32, filePathList []string) (goodsId int32, err error) {
	if name == "" || detail == "" {
		return -1, errors.New("name, price or password empty")
	}

	tx, err := db.SqlDb.Begin()
	if err != nil {
		return -1, err
	}

	res, err := tx.ExecContext(ctx, `insert into z_goods(id, uid, name, price, detail, create_time) 
values(null, ?, ?, ?, ?, ?)`, uid, name, price, detail, time.Now())
	if err != nil {
		return -1, err
	}
	goodsId64, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	goodsId = int32(goodsId64)

	log.Println("goodsId", goodsId)

	for _, path := range filePathList {
		if _, err = tx.ExecContext(ctx, `insert into z_goods_pic(id, uId, gId, path) values(null, ?, ?, ?)`,
			uid, goodsId, path); err != nil {
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	misc.Logger.Info("add good to sql success", zap.Any("goodId", goodsId))
	return
}
