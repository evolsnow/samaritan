/*
Package rpc contains some rpc actions
*/
package rpc

import (
	"fmt"
	"github.com/evolsnow/samaritan/common/log"
	pb "github.com/evolsnow/samaritan/gpns/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var RpcClientD pb.GPNSClient
var RpcClientF pb.GPNSClient

var Chats = make(chan string, 100)

// NewClientD returns a initialized domestic rpc server
func NewClientD(server string) pb.GPNSClient {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to domestic rpc server failed")
	}
	go receiveChat()
	return pb.NewGPNSClient(conn)
}

// NewClientF returns a initialized foreign rpc server
func NewClientF(server string) pb.GPNSClient {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to foreign rpc server failed")
	}
	return pb.NewGPNSClient(conn)
}

// SendMail calls rpc server to send a mail with given subject and mail body
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

// SendSMS calls rpc server to send a sms with given text
func SendSMS(to, text string) (err error) {
	log.Debug("calling rpc.SendSMS")
	resp, err := RpcClientD.SendSMS(context.Background(), &pb.SMSRequest{To: to, Text: text})
	if err != nil {
		log.Warn(err)
		return
	}
	if resp.Success {
		return nil
	}
	return fmt.Errorf(resp.Reason)
}

// SocketPush calls rpc server to push message to client with webSocket
func SocketPush(tokenList []string, msg string, extraInfo map[string]string) []string {
	log.Debug("calling rpc.SocketPush")
	spr, err := RpcClientD.SocketPush(context.Background(), &pb.SocketPushRequest{Message: msg, ExtraInfo: extraInfo, UserToken: tokenList})
	if err != nil {
		log.Error("socket push err:", err)
	}
	return spr.UserToken
}

// AppPush calls rpc server to push message to client with apple push notification system
func AppPush(tokenList []string, msg string, extraInfo map[string]string) {
	log.Debug("calling rpc.IOSPush")
	RpcClientF.ApplePush(context.Background(), &pb.ApplePushRequest{Message: msg, ExtraInfo: extraInfo, DeviceToken: tokenList})
}

//receive chat from rpc server
func receiveChat() {
	req := new(pb.ReceiveChatRequest)
	stream, err := RpcClientD.ReceiveMsg(context.Background(), req)
	if err != nil {
		log.Error(err)
	} else {
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
}
