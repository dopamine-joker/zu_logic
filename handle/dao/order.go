package dao

import (
	"context"
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"
	"go.uber.org/zap"
	"time"
)

type Order struct {
	Id           int32     `json:"id" db:"id"`
	BuyId        int32     `json:"buyid" db:"buyid"`
	BuyUserName  string    `json:"bname" db:"bname"`
	SellId       int32     `json:"sellid" db:"sellid"`
	SellUserName string    `json:"sname" db:"sname"`
	GId          int32     `json:"gid" db:"gid"`
	GName        string    `json:"gname" db:"gname"`
	Price        float64   `json:"price" db:"price"`
	Cover        string    `json:"cover" db:"cover"`
	Status       int32     `json:"status" db:"status"`
	Time         time.Time `json:"time" db:"time"`
}

type OrderStatus int32

const (
	COMMIT OrderStatus = iota
	SUCCESS
	FAIL
)

//AddOrder 增加一条订单信息
func AddOrder(ctx context.Context, buyid, sellid, gid int32, status OrderStatus) (int32, error) {

	tx, err := db.SqlDb.Begin()
	if err != nil {
		misc.Logger.Error("start trasaction err", zap.Error(err))
		return -1, err
	}

	res, err := tx.ExecContext(ctx, `insert into z_order(id, buyid, sellid, gid, status, time) values(null, ?, ?, ?, ?, ?);`,
		buyid, sellid, gid, int32(status), time.Now())
	if err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("add order err", zap.Error(err))
		return -1, err
	}

	oid64, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("get id err err", zap.Error(err))
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("commit tx err", zap.Error(err))
		return -1, err
	}

	oid := int32(oid64)
	return oid, nil
}

//GetBuyOrder 得到自己购买的订单
func GetBuyOrder(ctx context.Context, buyid int32) ([]Order, error) {
	var list []Order
	var err error
	if err = db.SqlDb.SelectContext(ctx, &list, `select t2.*, name as gname, price, cover from z_goods as g join 
    (select t.*, name as sname from z_user as u2 join 
        (select o.*, u.name as bname from z_order as o join 
            (select id, name from z_user where id = ?) as u where o.buyid = u.id) as t 
    where t.sellid = u2.id) as t2 
where t2.gid = g.id;`, buyid); err != nil {
		misc.Logger.Error("get buy order err", zap.Error(err))
		return list, err
	}
	return list, nil
}

//GetSellOrder 得到自己卖出的订单
func GetSellOrder(ctx context.Context, sellid int32) ([]Order, error) {
	var list []Order
	var err error
	if err = db.SqlDb.SelectContext(ctx, &list, `select t2.*, name as gname, price, cover from z_goods as g join 
    (select t.*, name as bname from z_user as u2 join 
        (select o.*, u.name as sname from z_order as o join 
            (select id, name from z_user where id = ?) as u 
        where o.sellid = u.id) as t where t.buyid = u2.id) as t2 
where t2.gid = g.id;`, sellid); err != nil {
		misc.Logger.Error("get buy order err", zap.Error(err))
		return list, err
	}
	return list, nil
}

//UpdateOrder 更新订单
func UpdateOrder(ctx context.Context, id int32, status OrderStatus) error {
	var err error
	if _, err = db.SqlDb.ExecContext(ctx, `update z_order set status = ? where id = ?`, int32(status), id); err != nil {
		misc.Logger.Error("update order err", zap.Error(err))
		return err
	}
	return nil
}
