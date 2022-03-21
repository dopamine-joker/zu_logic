package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"
	"go.uber.org/zap"
	"time"
)

type Goods struct {
	Id         int32     `json:"id" db:"id"`   // 物品id,可以用来找对应照片的路径
	UserId     int32     `json:"uid" db:"uid"` //上传该物品的用户id
	Name       string    `json:"name" db:"name"`
	Uname      string    `json:"uname" db:"uname"`
	Price      float64   `json:"price" db:"price"`
	GType      int32     `json:"gType" db:"type"`
	School     string    `json:"school" db:"school"`
	Detail     string    `json:"detail" db:"detail"`
	Cover      string    `json:"cover" db:"cover"` //封面url
	CreateTime time.Time `json:"create_time" db:"create_time"`
	PicList    []PicGoods
}

type PicGoods struct {
	Id      int32  `json:"id" db:"id"`       //图片id
	UserId  int32  `json:"userId" db:"UId"`  //用户的id
	GoodsId int32  `json:"goodsId" db:"gId"` //对应物品的id
	Path    string `json:"path" db:"path"`   //图片的路径
}

//AddGoods 增加一条物品信息
func AddGoods(ctx context.Context, name, detail string, price float64, gType int32, uid int32, school, coverPath string, filePathList []string) (goodsId int32, err error) {
	if name == "" || detail == "" {
		return -1, errors.New("name, price or password empty")
	}

	tx, err := db.SqlDb.Begin()
	if err != nil {
		return -1, err
	}

	res, err := tx.ExecContext(ctx, `insert into z_goods(id, uid, name, price, type, school, detail, cover, create_time) 
values(null, ?, ?, ?, ?, ?, ?, ?, ?)`, uid, name, price, gType, school, detail, coverPath, time.Now())
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

	for _, path := range filePathList {
		if _, err = tx.ExecContext(ctx, `insert into z_goods_pic(id, UId, gId, path) values(null, ?, ?, ?)`,
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

//GetGoodsByUserId 查询用户的商品
func GetGoodsByUserId(ctx context.Context, userId int32) ([]Goods, error) {
	var goodsList []Goods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &goodsList, `select zu.id as uid, zu.name as uname, zg.id, zg.name, zg.price, zg.type, zg.school, zg.detail, zg.cover, zg.create_time from z_user zu right join
    (select id, uid, name, price, type, school, detail, cover, create_time from z_goods where uid = ?) zg on zu.id = zg.uid`,
		userId); err != nil {
		return nil, err
	}
	return goodsList, nil
}

//GetGoods 翻页查找
func GetGoods(ctx context.Context, page, count int32) ([]Goods, error) {
	var goodsList []Goods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &goodsList, `select zu.id as uid, zu.name as uname, zg.id, zg.name, zg.price, zg.type, zg.school, zg.detail, zg.cover, zg.create_time from z_user zu right join
    (select id, uid, name, price, type, school, detail, cover, create_time from z_goods) zg on zu.id = zg.uid where zg.id in 
(select t.id from (select id from z_goods limit ?, ?) as t)`, page, count); err != nil {
		misc.Logger.Warn("GetGoods err", zap.Error(err))
		return nil, err
	}
	return goodsList, nil
}

//DelGoods 删除货物
func DelGoods(ctx context.Context, gid int32, uid int32) error {
	var err error

	tx, err := db.SqlDb.Begin()
	if err != nil {
		misc.Logger.Error("start tx err", zap.Error(err))
		return err
	}

	//删除货物
	if _, err = tx.ExecContext(ctx, `delete from z_goods where id = ? and uid = ?`, gid, uid); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("delete goods err", zap.Error(err))
		return err
	}

	//删除照片
	if _, err = tx.ExecContext(ctx, `delete from z_goods_pic where gId = ?`, gid); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("delete goods pic err", zap.Error(err))
		return err
	}

	//删除有关评论
	if _, err = tx.ExecContext(ctx, `delete from z_comment where gid = ?`, gid); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("delete goods comment err", zap.Error(err))
		return err
	}

	//删除收藏
	if _, err = tx.ExecContext(ctx, `delete from z_favorites where gid = ?`, gid); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("delete goods favorites err", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("commit tx err", zap.Error(err))
		return err
	}

	return nil
}

//GetGoodsDetail 根据商品id获取具体信息
func GetGoodsDetail(ctx context.Context, gid int32) (Goods, error) {
	var goods Goods
	var picList []PicGoods
	var err error

	tx, err := db.SqlDb.Begin()
	if err != nil {
		misc.Logger.Warn("picList start transaction err", zap.Error(err))
		return Goods{}, err
	}

	row := tx.QueryRowContext(ctx, `select zu.id as uid, zu.name as uname, zg.id, zg.name, zg.price, zg.type, zg.school, zg.detail, zg.cover, zg.create_time from z_user zu right join
    (select id, uid, name, price, type, school, detail, cover, create_time from z_goods where id = ?) zg on zu.id = zg.uid`, gid)

	if err = row.Scan(&goods.UserId, &goods.Uname, &goods.Id, &goods.Name, &goods.Price, &goods.School, &goods.Detail, &goods.Cover, &goods.CreateTime); err != nil {
		_ = tx.Rollback()
		return Goods{}, err
	}

	rows, err := tx.QueryContext(ctx, `select * from z_goods_pic where gid = ?`, gid)
	if err != nil {
		_ = tx.Rollback()
		return Goods{}, err
	}

	for rows.Next() {
		var pic PicGoods
		if err = rows.Scan(&pic.Id, &pic.UserId, &pic.GoodsId, &pic.Path); err != nil {
			_ = tx.Rollback()
			return Goods{}, err
		}
		picList = append(picList, pic)
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return Goods{}, err
	}

	goods.PicList = picList

	return goods, nil
}

//GetGoodsByName 根据名字查找货物
func GetGoodsByName(ctx context.Context, name string) ([]Goods, error) {
	var goodList []Goods
	var err error
	if err = db.SqlDb.SelectContext(ctx, &goodList, `select zu.id as uid, zu.name as uname, zg.id, zg.name, zg.price, zg.type, zg.school, zg.detail, zg.cover, zg.create_time from z_user zu right join
    (select id, uid, name, price, type, school, detail, cover, create_time from z_goods where name like ?) zg on zu.id = zg.uid`,
		fmt.Sprintf("%%%s%%", name)); err != nil {
		misc.Logger.Warn("GetGoods err", zap.Error(err))
		return nil, err
	}
	return goodList, nil
}
