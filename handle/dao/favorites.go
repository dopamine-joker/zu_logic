package dao

import (
	"context"

	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"

	"go.uber.org/zap"
)

type UserFavorites struct {
	Id    int32   `json:"id" db:"id"`
	UId   int32   `json:"uid" db:"uid"`
	GId   int32   `json:"gid" db:"gid"`
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
	Cover string  `json:"cover" db:"cover"`
}

func AddFavorites(ctx context.Context, uid, gid int32) (int32, error) {
	res, err := db.SqlDb.ExecContext(ctx, `insert into z_favorites(id, uid, gid) values(null, ?, ?)`,
		uid, gid)
	if err != nil {
		misc.Logger.Error("insert favorites err", zap.Error(err))
		return -1, err
	}
	fid64, err := res.LastInsertId()
	if err != nil {
		misc.Logger.Error("get insert id err", zap.Error(err))
		return -1, err
	}
	fid32 := int32(fid64)
	return fid32, nil
}

func DeleteFavorites(ctx context.Context, fid int32, uid int32) error {
	if _, err := db.SqlDb.ExecContext(ctx, `delete from z_favorites where id = ? and uid = ?`, fid, uid); err != nil {
		misc.Logger.Error("delete favorites err", zap.Error(err))
		return err
	}
	return nil
}

func GetUserFavorites(ctx context.Context, uid int32) ([]UserFavorites, error) {
	var userFavoritesList []UserFavorites
	var err error
	if err = db.SqlDb.SelectContext(ctx, &userFavoritesList, `select zf.id, zf.uid, zf.gid, zg.name, zg.price, zg.type, zg.cover from z_goods zg 
		right join (select * from z_favorites where uid = ?) zf on zg.id = zf.gid`, uid); err != nil {
		misc.Logger.Error("get favorites list err", zap.Error(err))
		return nil, err
	}
	return userFavoritesList, nil
}
