package rpc

import (
	"github.com/evolsnow/samaritan/common/log"
	pb "github.com/evolsnow/samaritan/gpns/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var RpcClientD pb.GPNSClient
var RpcClientF pb.GPNSClient

var Chats = make(chan string, 100)

func NewClientD(server string) pb.GPNSClient {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to domestic rpc server failed")
	}
	go receiveChat()
	return pb.NewGPNSClient(conn)
}

func NewClientF(server string) pb.GPNSClient {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to foreign rpc server failed")
	}
	return pb.NewGPNSClient(conn)
}

func SendMail(to, subject, body string) (err error) {
	log.Debug("calling rpc.SendMail")
	_, err = RpcClientF.SendMail(context.Background(), &pb.MailRequest{To: to, Subject: subject, Body: body})
	if err != nil {
		log.Warn(err)
	} else {
		log.Debug("mail sent")
	}
	return
}

func SocketPush(tokenList []string, msg string) []string {
	log.Debug("calling rpc.SocketPush")
	spr, err := RpcClientD.SocketPush(context.Background(), &pb.SocketPushRequest{Message: msg, UserToken: tokenList})
	if err != nil {
		log.Error("socket push err:", err)
	}
	return spr.UserToken
}

func IOSPush(tokenList []string, msg string) {
	log.Debug("calling rpc.IOSPush")
	RpcClientF.ApplePush(context.Background(), &pb.ApplePushRequest{Message: msg, DeviceToken: tokenList})
}

func receiveChat() {
	req := new(pb.ReceiveChatRequest)
	stream, err := RpcClientD.ReceiveMsg(context.Background(), req)
	if err != nil {
		log.Error(err)
	}
	for {
		rcv, err := stream.Recv()
		if err != nil {
			log.Warn(err)
			stream, _ = RpcClientD.ReceiveMsg(context.Background(), req)
			continue
		}
		Chats <- rcv.Chat
	}
}
