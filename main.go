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

func main() {
	core.InitLogger()
	c, err := core.InitConfig()
	if err != nil {
		fmt.Println("failed to init log: %v", err)
	}
	mainLogger := core.GetLogger("main")
	mainLogger.Info("starting")

	var httpPort int
	var raftPort int
	flag.IntVar(&httpPort, "-http-port", 18080, "http rest api port")
	flag.IntVar(&raftPort, "-raft-port", 20000, "raft tcp port")

	s := core.CreateServer(c)
	s.Run()
}

