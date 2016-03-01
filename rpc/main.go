package main

import (
	"flag"
	"fmt"
	pb "github.com/evolsnow/samaritan/rpc/protos"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var port int
	var addr string
	flag.StringVar(&addr, "l", "", "listen address")
	flag.IntVar(&port, "p", 10010, "port")
	flag.Parse()
	s := grpc.NewServer()
	pb.RegisterSamaritanServer(s, &server{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal("rpc server fataled")
	}
	log.Println("rpc server listening", fmt.Sprintf("%s:%d", addr, port))
	s.Serve(lis)
}
