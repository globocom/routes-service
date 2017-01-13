package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/deckarep/golang-set"
	"github.com/labstack/gommon/log"
	"github.com/osrg/gobgp/client"
	"github.com/osrg/gobgp/config"
	"github.com/pkg/errors"
)

type Neighbor struct {
	PeerAs          uint32
	NeighborAddress string
}

//PeersService ...
type PeersService struct {
	etcd   *clientv3.Client
	bgp    *client.Client
	prefix string
}

//NewPeersService ...
func NewPeersService(etcd *clientv3.Client, bgp *client.Client) *PeersService {
	s := &PeersService{
		etcd:   etcd,
		bgp:    bgp,
		prefix: "/bgp-service/neighbors",
	}

	return s
}

//Start ...
func (ps *PeersService) Start() error {
	go ps.watchChanges()

	return nil
}

//Start ...
func (ps *PeersService) watchChanges() error {
	rch := ps.etcd.Watcher.Watch(context.Background(), ps.prefix, []clientv3.OpOption{clientv3.WithPrefix()}...)
	for _ = range rch {
		if err := ps.handleChanges(); err != nil {
			return err
		}

	}

	return nil
}

func (ps *PeersService) handleChanges() error {
	desiredSet, err := ps.getDesiredSet()
	if err != nil {
		return err
	}

	currentSet, err := ps.getCurrentSet()
	if err != nil {
		return err
	}

	neighborsToAdd := desiredSet.Difference(currentSet)
	neighborsToDelete := currentSet.Difference(desiredSet)

	// Adding missing neighbors
	for n := range neighborsToAdd.Iter() {
		neighbor := n.(Neighbor)
		if err := ps.addNeighbor(neighbor); err != nil {
			return errors.Wrap(err, "[PeersService] Add neighbor failed")
		}

		log.Debugf("[PeersService] Added neighbor: %#v", neighbor)
	}

	// Deleting neighbors
	for n := range neighborsToDelete.Iter() {
		neighbor := n.(Neighbor)
		if err := ps.deleteNeighbor(neighbor); err != nil {
			return errors.Wrap(err, "[PeersService] Delete neighbor failed")
		}

		log.Debugf("[PeersService] Deleted neighbor: %#v", neighbor)
	}

	return nil
}

func (ps *PeersService) addNeighbor(n Neighbor) error {
	peer := &config.Neighbor{
		Config: config.NeighborConfig{
			PeerAs:          n.PeerAs,
			NeighborAddress: n.NeighborAddress,
		},
	}

	return ps.bgp.AddNeighbor(peer)
}

func (ps *PeersService) deleteNeighbor(n Neighbor) error {
	peer := &config.Neighbor{
		Config: config.NeighborConfig{
			PeerAs:          n.PeerAs,
			NeighborAddress: n.NeighborAddress,
		},
	}

	return ps.bgp.DeleteNeighbor(peer)
}

func (ps *PeersService) getDesiredSet() (mapset.Set, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := ps.etcd.KV.Get(ctx, ps.prefix, []clientv3.OpOption{clientv3.WithPrefix()}...)
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "get neighbors from etcd failed")
	}

	desiredSet := mapset.NewSet()

	for _, ev := range resp.Kvs {
		var n Neighbor
		if err := json.Unmarshal(ev.Value, &n); err != nil {
			return nil, errors.Wrap(err, "Unmarshal neighbor failed")
		}

		desiredSet.Add(n)
	}

	return desiredSet, nil
}

func (ps *PeersService) getCurrentSet() (mapset.Set, error) {
	neighbors, err := ps.bgp.ListNeighbor()
	if err != nil {
		return nil, err
	}

	currentSet := mapset.NewSet()

	for _, n := range neighbors {
		currentSet.Add(Neighbor{
			PeerAs:          n.Config.PeerAs,
			NeighborAddress: n.Config.NeighborAddress,
		})
	}

	return currentSet, nil
}

func (ps *PeersService) storeNeighbor(n Neighbor) error {
	key := fmt.Sprintf("%s/%s", ps.prefix, n.NeighborAddress)
	value, err := json.Marshal(n)
	if err != nil {
		return errors.Wrap(err, "[PeersService] Marshal neighbor failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = ps.etcd.KV.Put(ctx, key, string(value))
	cancel()
	if err != nil {
		return errors.Wrap(err, "store neighbor on etcd failed")
	}
	return nil
}
