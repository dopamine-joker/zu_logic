package handle

import (
	"context"
	"encoding/json"
	"github.com/dopamine-joker/zu_logic/db"
	"github.com/dopamine-joker/zu_logic/handle/dao"
	"github.com/dopamine-joker/zu_logic/misc"
	"go.uber.org/zap"
	"time"
)

func TaskAddOrder(ctx context.Context) {
	var err error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var result []string
				result, err = db.RedisClient.BRPop(ctx, time.Second*5, db.RedisOrderAdd).Result()
				if err != nil {
					misc.Logger.Info("task queue block timeout,no msg err", zap.String("err", err.Error()))
				}
				if len(result) >= 2 {
					if err = consumeAdd(ctx, result[1]); err != nil {
						_ = db.RedisClient.RPush(ctx, db.RedisOrderAdd, result[1]).Err()
					}
				}
			}
		}
	}()
}

func TaskUpdateOrder(ctx context.Context) {
	var err error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var result []string
				result, err = db.RedisClient.BRPop(ctx, time.Second*5, db.RedisOrderUpdate).Result()
				if err != nil {
					misc.Logger.Info("task queue block timeout,no msg err", zap.String("err", err.Error()))
				}
				if len(result) >= 2 {
					if err = consumeUpdate(ctx, result[1]); err != nil {
						_ = db.RedisClient.RPush(ctx, db.RedisOrderUpdate, result[1]).Err()
					}
				}
			}
		}
	}()
}

func consumeAdd(ctx context.Context, msg string) error {
	m := &dao.RedisOrderAdd{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		misc.Logger.Warn("json.Unmarshal err", zap.String("err", err.Error()))
		return err
	}
	misc.Logger.Info("push msg info", zap.Any("RedisMsg", m), zap.Any("redisOrder", m))
	orderId, err := dao.AddOrder(ctx, m.BuyId, m.SellId, m.GId, m.Status)
	misc.Logger.Info("add order to sql", zap.Int32("orderId", orderId))
	if err != nil {
		return err
	}
	return nil
}

func consumeUpdate(ctx context.Context, msg string) error {
	m := &dao.RedisOrderUpdate{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		misc.Logger.Warn("json.Unmarshal err", zap.String("err", err.Error()))
		return err
	}
	misc.Logger.Info("push msg info", zap.Any("RedisMsg", m), zap.Any("redisOrder", m))
	err := dao.UpdateOrder(ctx, m.OrderId, m.Status)
	misc.Logger.Info("update order to sql", zap.Int32("orderId", m.OrderId), zap.Int32("status", int32(m.Status)))
	if err != nil {
		return err
	}
	return nil
}
