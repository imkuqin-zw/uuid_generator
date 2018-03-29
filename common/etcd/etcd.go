package etcd

import "github.com/imkuqin-zw/uuid_generator/config"
import (
	"github.com/coreos/etcd/clientv3"
	"context"
	"strconv"
	"fmt"
	"github.com/imkuqin-zw/uuid_generator/common"
	"github.com/imkuqin-zw/uuid_generator/common/snowflake"
	"github.com/zieckey/etcdsync"
	"time"
)

const (
	PATH       = "/seqs/"
	UUID_KEY   = "/seqs/snowflake-uuid"
)

type Etcd struct {
	etcdSync map[string]*etcdsync.Mutex
	client *clientv3.Client
	root string
}

func NewEtcd(c *config.Etcd, locks map[string]*config.EtcdSync) (etcdv3 *Etcd, err error) {
	etcdv3 = &Etcd{
		root: c.Root,
		etcdSync: make(map[string]*etcdsync.Mutex),
	}
	cfg := clientv3.Config{
		Endpoints: c.Addrs,
		DialTimeout: c.TimeOut,
	}
	etcdv3.client, err = clientv3.New(cfg)
	if err != nil {
		return
	}
	for _, item := range locks {
		var m *etcdsync.Mutex
		m, err = etcdsync.New(item.Root + item.Name, int(item.Ttl / time.Second), item.Addrs)
		if m == nil || err != nil {
			if err == nil {
				err = fmt.Errorf("etcdsync.New failed")
			}
			fmt.Println("etcdsync.New failed")
			return
		}
		etcdv3.etcdSync[item.Name] = m
	}
	return
}

func (e *Etcd) GetNextByName(name string) (next int64, err error) {
	var prevValue int64
	key := PATH + name
	var resp *clientv3.GetResponse
	if err = e.etcdSync[name].Lock(); err != nil {
		return
	}
	defer e.etcdSync[name].Unlock()
	for {
		if resp, err = e.client.Get(context.Background(), key); err != nil {
			return
		}
		for _, val := range resp.Kvs {
			prevValue, err = strconv.ParseInt(string(val.Value), 10, 64)
			if err != nil {
				err = fmt.Errorf("marlformed value")
				return
			}
		}
		_, err = e.client.Put(context.Background(), key, fmt.Sprint(prevValue+1))
		if err != nil {
			common.CasDelay()
			continue
		}
		next = prevValue + 1
		return
	}
	return
}

func (e *Etcd) GetMachineID() (machineId uint64, err error) {
	var prevValue int
	var resp *clientv3.GetResponse
	for {
		if resp, err = e.client.Get(context.Background(), UUID_KEY); err != nil {
			return
		}
		for _, val := range resp.Kvs {
			prevValue, err = strconv.Atoi(string(val.Value))
			if err != nil {
				err = fmt.Errorf("marlformed value")
				return
			}
		}
		_, err = e.client.Put(context.Background(), UUID_KEY, fmt.Sprint(prevValue+1))
		if err != nil {
			common.CasDelay()
			continue
		}
		machineId = (uint64(prevValue+1) & snowflake.MACHINE_ID_MASK) << 12
		return
	}
}