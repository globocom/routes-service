package main

import (
	"context"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/osrg/gobgp/client"
)

//PeersService ...
type PeersService struct {
	etcd *clientv3.Client
	bgp  *client.Client
}

//NewPeersService ...
func NewPeersService(etcd *clientv3.Client, bgp *client.Client) *PeersService {
	s := &PeersService{
		etcd: etcd,
		bgp:  bgp,
	}

	return s
}

//Start ...
func (ps *PeersService) Start() error {
	rch := ps.etcd.Watcher.Watch(context.Background(), "/routes-service/peers", []clientv3.OpOption{clientv3.WithPrefix()}...)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

	return nil
}
