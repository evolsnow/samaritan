package base

import (
	"github.com/evolsnow/samaritan/base/log"
	pb "github.com/evolsnow/samaritan/rpc/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client pb.SamaritanClient

func init() {
	conn, err := grpc.Dial("127.0.0.1:10010", grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to rpc server failed")
	}
	client = pb.NewSamaritanClient(conn)
}

func SendMail(to, subject, body string) (err error) {
	log.Debug("calling rpc.SendMail")
	_, err = client.SendMail(context.Background(), &pb.MailRequest{To: to, Subject: subject, Body: body})
	return
}

func IOSPush(deviceList []string, msg string) {
	log.Debug("calling rpc.SendMail")
	client.ApplePush(context.Background(), &pb.PushRequest{Message: msg, DeviceToken: deviceList})
}
