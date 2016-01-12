package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var bot *tgbotapi.BotAPI
var session = make(map[string]bool)

func main() {
	//test()
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
			}
		} else if session[userName] { //talking and received a different command
			delete(session, userName)
		}
		if _, ok := funcMap[endPoint]; ok {
			rawMsg = funcMap[endPoint](update)
		} else {
			rawMsg = "unknown command"
		}
	} else if session[userName] { //user have an existing talking session
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
	var info string
	if len(text) == 1 && text[0] == "/talk" {
		return "now you can talk to me, type any text starts with '/' to exit the talk"
	} else if text[0] == "/talk" {
		info = strings.Join(text[1:], " ")
	} else {
		info = update.Message.Text
	}
	log.Println(info)

	var response string
	for _, r := range info {

		if unicode.Is(unicode.Scripts["Han"], r) {
			log.Println("汉语")
			if len(strings.Split(info, " ")) > 1 {
				return "中文就不要用空格分隔啦"
			}
			response = tlAI(info)
			break
		} else {
			log.Println("英语")
			response = mitAI(info)
			break
		}

	}
	return response
}

func start(update tgbotapi.Update) string {
	return "welcome: " + update.Message.Chat.UserName
}

func tlAI(info string) string {
	key := "a5052a22b8232be1e387ff153e823975"
	tuLingUrl := fmt.Sprintf("http://www.tuling123.com/openapi/api?key=%s&info=%s", key, info)
	resp, err := http.Get(tuLingUrl)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	reply := new(tlReply)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(reply)
	return strings.Replace(reply.Text, "<br>", "\n", -1)
}

type tlReply struct {
	code int    `json:"code"`
	Text string `json:"text"`
}

func mitAI(info string) string {
	mitUrl := "http://fiddle.pandorabots.com/pandora/talk?botid=9fa364f2fe345a10&skin=demochat"
	resp, err := http.PostForm(mitUrl, url.Values{"message": {info}, "botcust2": {"d064e07d6e067535"}})
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	re, _ := regexp.Compile("Mitsuku:</B>(.*?)<br> <br>")
	all := re.FindAll(body, -1)
	if len(all) == 0 {
		return "change another question?"
	}
	found := (string(all[0]))
	log.Println(found)
	ret := strings.Replace(found, `<P ALIGN="CENTER"><img src="http://`, "", -1)
	ret = strings.Replace(ret, `"></img></P>`, "", -1)
	ret = strings.Replace(ret[13:], "<br>", "\n", -1)
	ret = strings.Replace(ret, "Mitsuku", "samaritan", -1)
	return ret
}

func test() {
	str := `<html>
<head>
<title>text page</title>
<script>
<!--
function sf() {document.f.message.focus();}
// -->
</script>
<STYLE type="text/css"><!--
   a:link { color: #FFFF00; text-decoration: none }
   a:visited { color: #FFFF00; text-decoration: none }
   a:hover { color: #FFFFFF; text-decoration: none }
   --></STYLE></HEAD>
</head>

<body onload="sf()" BGCOLOR="#FFFFFF">
<font face="verdana,ariel" size="3" color="000000">

<center>
<FONT FACE="Trebuchet MS,Arial" COLOR="#990000">
<form method="POST" name="f"><input type="hidden" name="botcust2" value="d064e07d6e067535">
<i><b>Type your message to Mitsuku:</b></i><br>
<input autocomplete="off" type="TEXT" name="message" maxlength="500" size="30"><input type="submit"

value="enter">
</form></center>
<P>
<FONT FACE="Trebuchet MS,Arial" COLOR="#000000">

 <B> You:</B>  tell me a joke<br> <B> Mitsuku:</B>   David Hasselhoff walks into a bar and says to the barman, "I want you to call me David Hoff".<br><br> The barman replies "Sure thing Dave... no hassle".<br> <br>
</font>

</CENTER>

</BODY>
</HTML>`
	re, _ := regexp.Compile("Mitsuku:</B>(.*)")
	all := re.FindAll([]byte(str), -1)
	re2, _ := regexp.Compile("Has(.*?)ff")
	ret := re2.ReplaceAllString(string(all[0])[15:], "")

	fmt.Println(ret)
	os.Exit(1)
}
