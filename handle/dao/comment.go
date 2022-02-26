package dao

import (
	"context"
	"time"

	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/misc"

	"go.uber.org/zap"
)

type Comment struct {
	Id      int32     `json:"id" db:"id"`
	UId     int32     `json:"uid" db:"uid"`
	GId     int32     `json:"gid" db:"gid"`
	Oid     int32     `json:"oid" db:"oid"`
	Content string    `json:"content" db:"content"`
	Level   int32     `json:"level" db:"level"`
	TIme    time.Time `json:"time" db:"time"`
}

type UserComment struct {
	Comment
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
	Cover string  `json:"cover" db:"cover"`
}

type GoodsComment struct {
	Comment
	UserName string `json:"uname" db:"uname"`
}

//AddComment 增加一条评论
func AddComment(ctx context.Context, userId, goodsId, oid, level int32, content string) (commentId int32, err error) {

	tx, err := db.SqlDb.Begin()
	if err != nil {
		misc.Logger.Error("start transaction err", zap.Error(err))
		return -1, err
	}

	res, err := tx.ExecContext(ctx, `insert into z_comment(id, uid, gid, oid, content, level) 
values(null, ?, ?, ?, ?, ?)`, userId, goodsId, oid, content, level)
	if err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("insert a comment error", zap.Error(err))
		return -1, err
	}
	commentId64, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		misc.Logger.Error("get last insert id err", zap.Error(err))
		return -1, err
	}
	commentId = int32(commentId64)

	if _, err = tx.ExecContext(ctx, `update z_order set status = 4 where id = ?`, oid); err != nil {
		misc.Logger.Error("update order err", zap.Error(err))
		return -1, err
	}

	_ = tx.Commit()
	return commentId, nil
}

//GetCommentByUserId 根据用户id，用户查询其评论
func GetCommentByUserId(ctx context.Context, userId int32) ([]UserComment, error) {
	var commentList []UserComment
	var err error
	if err = db.SqlDb.SelectContext(ctx, &commentList, `select zc.id, zc.uid, zc.gid, zc.oid, zc.content, zc.level, zc.time, zg.name, zg.price, zg.cover from z_goods zg 
right join (select * from z_comment zc where uid = ?) zc on zg.id = zc.gid`, userId); err != nil {
		misc.Logger.Error("get comment error", zap.Error(err))
		return nil, err
	}
	return commentList, nil
}

//GetCommentByGoodsId 根据货物id拉取评论
func GetCommentByGoodsId(ctx context.Context, gid int32) ([]GoodsComment, error) {
	var commentList []GoodsComment
	var err error
	if err = db.SqlDb.SelectContext(ctx, &commentList, `select zc.id, zc.uid, zc.gid, zc.oid, zc.content, zc.level, zc.time, zu.name as uname from z_user zu right join 
    (select * from z_comment where gid = ?) zc on zc.uid = zu.id;`, gid); err != nil {
		misc.Logger.Error("get comment list by gid err", zap.Error(err))
		return nil, err
	}
	return commentList, nil
}

//DeleteComment 删除评论
func DeleteComment(ctx context.Context, commentId int32) error {
	if _, err := db.SqlDb.ExecContext(ctx, `delete from z_comment where id = ?`, commentId); err != nil {
		misc.Logger.Error("delete comment err", zap.Error(err))
		return err
	}
	return nil
}
