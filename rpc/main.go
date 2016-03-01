package main

import (
	pb "github.com/evolsnow/samaritan/rpc/protos"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterSamaritanServer(s, &server{})
	lis, err := net.Listen("tcp", ":10010")
	if err != nil {
		log.Fatal("rpc server fataled")
	}
	log.Println("rpc server listening: 127.0.0.1:10010")
	s.Serve(lis)
}
