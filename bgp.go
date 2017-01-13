package main

import (
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/osrg/gobgp/client"
	"github.com/pkg/errors"
)

//Config ...
type Config struct {
	EtcdEndpoints []string
	GobgpAddress  string
}

//BgpService ...
type BgpService struct {
	etcdClient    *clientv3.Client
	gobgpClient   *client.Client
	peersService  *PeersService
	routesService *RoutesService
}

//Shutdown ...
func (bs *BgpService) Shutdown() error {
	if err := bs.etcdClient.Close(); err != nil {
		return nil
	}

	if err := bs.gobgpClient.Close(); err != nil {
		return nil
	}

	return nil
}

//NewBgpService ...
func NewBgpService(config Config) (*BgpService, error) {
	//Initialize etcd client
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   config.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create a new BgpService failed")
	}

	// Initialize gobgp client
	gobgp, err := client.New(config.GobgpAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Intialize peer service

	s := &BgpService{
		etcdClient:  etcd,
		gobgpClient: gobgp,
	}

	return s, nil
}

// Start ...
func (bs BgpService) Start() error {
	go bs.peersService.Start()

	return nil
}
