/*
 * Filename: /mnt/c/Users/PureAdmin/tj/courses/raft/store.go
 * Path: /mnt/c/Users/PureAdmin/tj/courses/raft
 * Created Date: Thursday, January 1st 1970, 8:00:00 am
 * Author: strawhatboy
 *
 * Copyright (c) 2020 Your Company
 */

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	log "github.com/sirupsen/logrus"
)

const DEFAULT_TIMEOUT = 10 * time.Second


type IStore interface {
	raft.FSM
	ApplyPut(key string, value string) interface{}
	ApplyDelete(key string) interface{}
}
// Store
// implements raft.FSM
type Store struct {
	id 	string
	raft   *raft.Raft
	logger *log.Entry

	lock	sync.Mutex
	data	map[string]string
}

type Snapshot struct {
	data	map[string]string
}

var snapShotLogger = GetLogger("snapshot")

func NewStore(singleNode bool, id string, raftAddr string, raftPath string, joinAddr string) (*Store, error) {
	s := &Store{
		logger: GetLogger("store"),
		data:	make(map[string]string),
		id: 	id,
	}
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(id)

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		s.logger.Error("error when resolving tcp addr with: ", raftAddr, err)
		return nil, err
	}

	transport, err := raft.NewTCPTransport(raftAddr, tcpAddr, 10, DEFAULT_TIMEOUT, os.Stdout)
	if err != nil {
		s.logger.Error("error when creating TCP transport with: ", raftAddr, tcpAddr, err)
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(raftPath, 10, os.Stdout)
	if err != nil {
		s.logger.Error("error when creating snapshot store in path: ", raftPath, err)
		return nil, err
	}

	store := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	s.raft, err = raft.NewRaft(conf, s, store, stableStore, snapshotStore, transport)
	if err != nil {
		s.logger.Error("error when creating raft instance ", err)
		return nil, err
	}

	if singleNode {
		// start the cluster
		s.raft.BootstrapCluster(raft.Configuration{
			Servers: []raft.Server{
				{
					ID:	conf.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		})
	} else {
		// should join other cluster
		body := JoinRequestBody{ID: id, Addr: raftAddr}
		if j, err := json.Marshal(body); err == nil {
			if res, err := http.Post(joinAddr + "/join", "text/json", bytes.NewReader(j)); err != nil {
				s.logger.Error("failed to send join request to: ", joinAddr, err)
				return nil, err
			} else {
				s.logger.Info(fmt.Sprintf("sent json %v and got response %v", string(j), res))
			}
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (s *Store) Get(key string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok := s.data[key]
	if !ok {
		return ""
	}
	return v
}

func (s *Store) Put(key string, value string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("current node is not a leader")
	}

	operation := Operation{Name: "put", Key: key, Value: value }
	b, err := operation.marshal()
	if err != nil {
		s.logger.Error("failed to marshal operation put with key&value: ", key, value, err)
		return err
	} else {
		s.logger.Info("marshaled json: ", string(b))
	}

	future := s.raft.Apply(b, DEFAULT_TIMEOUT)
	return future.Error()
}

func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return fmt.Errorf("current node is not a leader")
	}

	operation := Operation{Name: "delete", Key: key}
	b, err := operation.marshal()
	if err != nil {
		s.logger.Error("failed to marshal operation delete with key: ", key, err)
		return err
	}

	future := s.raft.Apply(b, DEFAULT_TIMEOUT)
	return future.Error()
}

func (s *Store) Join(id string, addr string) error {
	s.logger.Info(fmt.Sprintf("going to handle the join request from %v, %v", addr, id))

	future := s.raft.GetConfiguration()
	if err := future.Error(); err != nil {
		s.logger.Error("failed to get raft configuration ", err)
		return err
	}

	for _, server := range future.Configuration().Servers {
		matchID := server.ID == raft.ServerID(id)
		matchAddr := server.Address == raft.ServerAddress(addr)
		if matchID && matchAddr {
			s.logger.Info(fmt.Sprintf("node %v %v already joined", id, addr))
			return nil
		} else if matchID != matchAddr {
			// remove
			f := s.raft.RemoveServer(server.ID, 0, 0)
			if err := f.Error(); err != nil {
				msg := fmt.Sprintf("failed to remove server %v %v", server.ID, server.Address)
				s.logger.Error(msg)
				return err
			}
		}

	}

	f := s.raft.AddVoter(raft.ServerID(id), raft.ServerAddress(addr), 0, 0)
	if err := f.Error(); err != nil {
		s.logger.Error(fmt.Sprintf("failed to add voter %v, %v", id, addr))
		return err
	}

	s.logger.Info(fmt.Sprintf("successfully added voder %v, %v", id, addr))
	return nil
}

// implements IStore
func (s *Store) ApplyPut(key string, value string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = value
	return nil
}

func (s *Store) ApplyDelete(key string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}
// implements raft.FSM
func (s *Store) Apply(l *raft.Log) interface{} {
	o := Operation{}
	err := o.unmarshal(l.Data)
	if err != nil {
		s.logger.Error("failed to unmarshal from raft.Log: ", err)
	} else {
		s.logger.Info("unmarshaled apply log: ", o)
	}

	return o.apply(s)
}

func (s *Store) Restore(rc io.ReadCloser) error {
	data := make(map[string]string)
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		s.logger.Error("failed to restore data ", err)
		return err
	}
	s.data = data
	return nil
}

func (s *Store) Snapshot() (raft.FSMSnapshot, error) {
	newMap := make(map[string]string)
	for k, v := range s.data {
		newMap[k] = v
	}
	return &Snapshot{data: newMap}, nil
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	b, err := json.Marshal(s.data)
	if err != nil {
		snapShotLogger.Error("failed to unmarshal data when persist snapshot ", err)
		return err
	}

	// write to sink
	_, err = sink.Write(b)
	if err != nil {
		snapShotLogger.Error("failed to write data to sink ", err)
		return err
	}

	err = sink.Close()
	if err != nil {
		snapShotLogger.Error("failed to close sink ", err)
		sink.Cancel()
	}

	return err
}

func (s *Snapshot) Release() {}

