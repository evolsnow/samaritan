package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var bot *tgbotapi.BotAPI
var session = make(map[string]bool)

func main() {
	//used for 104
	//go http.ListenAndServeTLS("0.0.0.0:8443", "server.crt", "server.key", nil)
	go http.ListenAndServe("0.0.0.0:8000", nil)

	var err error
	bot, err = tgbotapi.NewBotAPI("164760320:AAEE0sKLgCwHGYJ0Iqz7o-GYH4jVTQZAZho")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//used for 104
	//_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://104.236.156.226:8443/"+bot.Token, "server.crt"))
	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://www.samaritan.tech:8443/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.ListenForWebhook("/" + bot.Token)
	for update := range updates {
		go handlerConnection(update)

	}
}

func handlerConnection(update tgbotapi.Update) {
	userName := update.Message.Chat.UserName
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	rawMsg := fmt.Sprintf("Hi %s, you said: %s", userName, text)
	funcMap := map[string]func(update tgbotapi.Update) string{
		"/start": start,
		"/talk":  talk,
	}
	if string(text[0]) == "/" {
		received := strings.Split(text, " ")
		endPoint := received[0]
		if endPoint == "/talk" { //begin to talk
			if len(received) == 1 { //not like "/talk something"
				session[userName] = true
				log.Println("begin to talk", session[userName])
			}
		} else if session[userName] { //talking and received a different command
			log.Println("received a different cmd")
			delete(session, userName)
		}
		if _, ok := funcMap[endPoint]; ok {
			log.Println("reply from func")
			rawMsg = funcMap[endPoint](update)
		} else {
			rawMsg = "unknown command"
		}
	} else if session[userName] { //user have an existing talking session
		log.Println("existing session")
		rawMsg = talk(update)
	}

	msg := tgbotapi.NewMessage(chatId, rawMsg)
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

}

func talk(update tgbotapi.Update) string {
	text := strings.Split(update.Message.Text, " ")
	if len(text) == 1 && text[0] == "/talk" {
		return "now you can talk to me, type any other command to exit the talk"
	}
	return tlAI(text[0])
}

func start(update tgbotapi.Update) string {
	return "welcome: " + update.Message.Chat.UserName
}

func tlAI(info string) string {
	key := "a5052a22b8232be1e387ff153e823975"
	apiUrl := fmt.Sprintf("http://www.tuling123.com/openapi/api?key=%s&info=%s", key, info)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Content-type", "text/html")
	req.Header.Add("charset", "utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
	var reply tlReply

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(reply)
	return string(reply.text)
}

type tlReply struct {
	code int
	text string
}
