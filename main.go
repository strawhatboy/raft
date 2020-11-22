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
	DEFAULT_HTTP_ADDR = "0.0.0.0:18080"
	DEFAULT_RAFT_ADDR = "0.0.0.0:20000"
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
	var id	string
	flag.StringVar(&httpAddr, "-http-addr", DEFAULT_HTTP_ADDR, "http rest api port")
	flag.StringVar(&raftAddr, "-raft-addr", DEFAULT_RAFT_ADDR, "raft tcp port")
	flag.StringVar(&joinAddr, "-join-addr", "", "the address to join into")
	flag.StringVar(&id, "-id", "", "the id of the node")

	if httpAddr != DEFAULT_HTTP_ADDR {
		c.HttpAddr = httpAddr
	}
	if raftAddr != DEFAULT_RAFT_ADDR {
		c.RaftAddr = raftAddr
	}

	s := core.CreateServer(c)
	s.Run()
}

