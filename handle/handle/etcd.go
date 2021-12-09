package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/dopamine-joker/zu_logic/misc"
)

type Register struct {
	EtcdAddrs  []string
	DiaTimeout int

	basePath    string
	serverPath  string
	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	client      *clientv3.Client
	server      *RpcLogicServer
	srvTTL      int64
}

func NewRegister(etcdAddrs []string, BasePath, ServerPath string, diaTimeout int) (*Register, error) {
	var err error
	register := &Register{
		EtcdAddrs:   etcdAddrs,
		DiaTimeout:  diaTimeout,
		basePath:    BasePath,
		serverPath:  ServerPath,
		closeCh:     make(chan struct{}),
		keepAliveCh: make(<-chan *clientv3.LeaseKeepAliveResponse),
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: time.Duration(misc.Conf.EtcdCfg.TimeOut) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	register.client = client
	return register, nil
}

func (r *Register) Register(ctx context.Context, server *RpcLogicServer, ttl int64) error {
	var err error
	if r.client == nil {
		return errors.New("register.client == nil")
	}
	r.server = server
	r.srvTTL = ttl

	if err = r.register(ctx); err != nil {
		return err
	}

	go r.keepAlive()

	return nil
}

func (r *Register) register(ctx context.Context) error {
	leaseCtx, cancel := context.WithTimeout(ctx, time.Duration(r.DiaTimeout)*time.Second)
	defer cancel()

	// create new lease
	leaseRps, err := r.client.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseRps.ID

	// lease keepAlive
	if r.keepAliveCh, err = r.client.KeepAlive(ctx, leaseRps.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.server)
	if err != nil {
		return err
	}
	key := r.buildKeyPath()
	_, err = r.client.Put(ctx, key, string(data), clientv3.WithLease(r.leasesID))

	return nil
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			misc.Logger.Info("close, unregister etcd")
			// 删除Key
			if err := r.unregister(context.Background()); err != nil {
				misc.Logger.Error("unregister server failed", zap.String("path", r.buildKeyPath()))
			}
			// 撤销租约
			if _, err := r.client.Revoke(context.Background(), r.leasesID); err != nil {
				misc.Logger.Error("revoke leases fail", zap.Int64("leaseID", int64(r.leasesID)))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(context.Background()); err != nil {
					misc.Logger.Error("register failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(context.Background()); err != nil {
					misc.Logger.Error("register failed", zap.Error(err))
				}
			}

		}
	}
}

func (r *Register) unregister(ctx context.Context) error {
	_, err := r.client.Delete(ctx, r.buildKeyPath())
	return err
}

func (r *Register) buildKeyPath() string {
	return fmt.Sprintf("/%s/%s/%s", r.basePath, r.serverPath, r.server.Addr)
}

func (r *Register) StopServe() {
	r.closeCh <- struct{}{}
}
