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

package main

import (
	"fmt"
	//	"io"
	"log"
	"net"
	"os"
	//	"strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"spos_auth/datas"
	pb "spos_proto/uni"
	"time"
)

const (
	port     = "61051"
	ssl_port = ":61052"
	spos_ssl = false //true
)

//// server is used to implement helloworld.GreeterServer.
//type server struct{}

//// SayHello implements helloworld.GreeterServer
//func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
//	fmt.Printf("%s:Spos SayHello111111 Text=%s\n", ulog.CreateDateString(), in.Name)
//	return &pb.HelloReply{Message: "Hello111 " + in.Name}, nil
//}

func main() {

	var addr_port string

	//	go func() {
	//		for _ = range time.Tick(1e10) {
	//			//		ulog.Ul.Debugf("%s:cur con num=%d;dev con in lis=%d;pc con in lis=%d;\n", ulog.CreateDateSting(), cur_conn_num, uconnect.CountDevConnList(), uconnect.CountPcConnList())
	//			fmt.Printf("%s:\n", ulog.CreateDateString())
	//		}
	//	}()

	if datas.OnlyOpenUnDb() != nil {
		fmt.Printf("main: db open error!\n")
		panic("db connection failure")
	}

	go func() {
		for _ = range time.Tick(5 * time.Minute) {
			WorkUserHeartTimeOutHandle()
		}
	}()

	port_new := os.Getenv("PORT")
	if len(port_new) == 0 {
		port_new = port
	}
	port_new = ":" + port_new

	if spos_ssl {
		addr_port = port_new
		fmt.Println("tls start...,port=", addr_port)
	} else {
		addr_port = port_new
		fmt.Println("not tls start...,port=", addr_port)
	}

	lis, err := net.Listen("tcp", addr_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var s *grpc.Server
	if spos_ssl {

		fmt.Println("spos-lang with ssl")
		creds, err := credentials.NewServerTLSFromFile("./gubstech.crt", "./gubstech.key")

		if err != nil {
			fmt.Println("Failed to generate credentials: ", err)
		}

		s = grpc.NewServer(grpc.Creds(creds))
	} else {
		s = grpc.NewServer()
	}

	pb.RegisterLangServer(s, &lang_server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

//func main() {
//	//	lis, err := net.Listen("tcp", port)
//	//	if err != nil {
//	//		log.Fatalf("failed to listen: %v", err)
//	//	}
//	//	s := grpc.NewServer()
//	//	pb.RegisterGreeterServer(s, &server{})
//	//	// Register reflection service on gRPC server.
//	//	reflection.Register(s)
//	//	if err := s.Serve(lis); err != nil {
//	//		log.Fatalf("failed to serve: %v", err)
//	//	}

//	port := os.Getenv("PORT")
//	fmt.Println("Lang-server start...,port=", port)

//	if len(port) == 0 {
//		port = "50051"
//	}
//	port = ":" + port

//	lis, err := net.Listen("tcp", port)
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}

//	s := grpc.NewServer()
//	pb.RegisterLangServer(s, &lang_server{})
//	// Register reflection service on gRPC server.
//	reflection.Register(s)
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}

//}
