/*
 * Filename: /mnt/c/Users/PureAdmin/tj/courses/raft/main.go
 * Path: /mnt/c/Users/PureAdmin/tj/courses/raft
 * Created Date: Thursday, January 1st 1970, 8:00:00 am
 * Author: strawhatboy
 *
 * A very simple key-value store with hashicorp raft (https://github.com/hashicorp/raft) as its backend
 *
 * Copyright (c) 2020 Your Company
 */

package main

import (
	"flag"
	"fmt"

	"github.com/strawhatboy/raft/core"
)

const (
	DEFAULT_HTTP_ADDR = "127.0.0.1:18080"
	DEFAULT_RAFT_ADDR = "127.0.0.1:20000"
	DEFAULT_RAFT_PATH = "./raft"
)

func main() {
	core.InitLogger()
	c, err := core.InitConfig()
	if err != nil {
		fmt.Println("failed to init log: ", err)
	}
	mainLogger := core.GetLogger("main")
	mainLogger.Info("starting")

	var httpAddr string
	var raftAddr string
	var joinAddr string
	var raftPath string
	var id	string
	var singleNode bool
	flag.StringVar(&httpAddr, "http-addr", DEFAULT_HTTP_ADDR, "http rest api port")
	flag.StringVar(&raftAddr, "raft-addr", DEFAULT_RAFT_ADDR, "raft tcp port")
	flag.StringVar(&joinAddr, "join-addr", "", "the address to join into")
	flag.StringVar(&id, "id", "", "the id of the node")
	flag.StringVar(&raftPath, "raft-path", DEFAULT_RAFT_ADDR, "raft path for snapshot")

	flag.Parse()

	mainLogger.Info("launching with flags: ", []string{httpAddr, raftAddr, joinAddr, raftPath, id})

	if httpAddr != DEFAULT_HTTP_ADDR {
		c.HttpAddr = httpAddr
	}
	if raftAddr != DEFAULT_RAFT_ADDR {
		c.RaftAddr = raftAddr
	}
	if joinAddr != "" {
		c.JoinAddr = joinAddr
		singleNode = false
	} else {
		singleNode = true
	}
	c.SingleNode = &singleNode
	if id != "" {
		c.ID = id
	}
	if raftPath != DEFAULT_RAFT_PATH {
		c.RaftPath = raftPath
	}

	s, err := core.CreateServer(c)
	if err != nil {
		mainLogger.Error("failed to create server ", err)
		return
	}

	s.Run()
}

