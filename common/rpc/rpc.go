package rpc

import (
	"github.com/evolsnow/samaritan/common/log"
	pb "github.com/evolsnow/samaritan/gpns/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client pb.GPNSClient
var Chats = make(chan string, 100)

func init() {
	conn, err := grpc.Dial("127.0.0.1:10086", grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to rpc server failed")
	}
	client = pb.NewGPNSClient(conn)
	go receiveChat()
}

func SendMail(to, subject, body string) (err error) {
	log.Debug("calling rpc.SendMail")
	_, err = client.SendMail(context.Background(), &pb.MailRequest{To: to, Subject: subject, Body: body})
	return
}

func SocketPush(tokenList []string, msg string) []string {
	log.Debug("calling rpc.SocketPush")
	spr, err := client.SocketPush(context.Background(), &pb.SocketPushRequest{Message: msg, UserToken: tokenList})
	if err != nil {
		log.Error("socket push err:", err)
	}
	return spr.UserToken
}

func IOSPush(tokenList []string, msg string) {
	log.Debug("calling rpc.IOSPush")
	client.ApplePush(context.Background(), &pb.ApplePushRequest{Message: msg, DeviceToken: tokenList})
}

func receiveChat() {
	req := new(pb.ReceiveChatRequest)
	stream, err := client.ReceiveMsg(context.Background(), req)
	if err != nil {
		log.Error(err)
	}
	for {
		rcv, err := stream.Recv()
		if err != nil {
			log.Warn(err)
			stream, err = client.ReceiveMsg(context.Background(), req)
			continue
		}
		Chats <- rcv.Chat
	}
}
