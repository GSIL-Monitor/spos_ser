/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto
//go:generate root@uniser03:/usr/lhh/gopath/src/spos_proto/uni/voice# protoc -I . --go_out=plugins=grpc:. chat.proto

//	pb "spos_proto/uni/voice"
package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"spos_auth/datas"
	proto "spos_proto/uni/voice" //

	"time"
)

func main() {
	log.Println("start chat server...")

	if datas.OnlyOpenUnDb() != nil {
		fmt.Printf("main: db open error!\n")
		panic("db connection failure")
	}

	go func() {
		for _ = range time.Tick(5 * time.Minute) {
			WorkUserHeartTimeOutHandle()
		}
	}()

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "61055"
	}
	port = ":" + port

	fmt.Println("Chat-server start...,port=", port)
	lis, err := net.Listen("tcp", port) //192.168.10.82:61055, 139.224.221.196:61056

	server := grpc.NewServer()
	// register ChatServer
	proto.RegisterChatServer(server, &SposChat{})
	if err != nil {
		panic(err)
	}

	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
