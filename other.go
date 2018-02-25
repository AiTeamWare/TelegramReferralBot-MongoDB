package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"gopkg.in/mgo.v2/bson"
)

func requestRepeat(message *tgbotapi.Message) {
	sendMessage(message.Chat.ID, phrases[4] + message.Text + phrases[5], nil)
	pending[message.From.ID] = 2
}

func submit(message *tgbotapi.Message) {
	var user User
	//db.First(&user, "id = ?", message.From.ID)
	err := db.Collection("users").FindOne(bson.M{"telegramid": message.From.ID}, &user)
	if err != nil {
		log.Panic(err)
	}
	user.EthAddress = message.Text
	//db.Save(&user)
	err = db.Collection("users").Save(&user)
	if err != nil {
		log.Panic(err)
	}
	sendMessage(message.Chat.ID, phrases[7] + message.Text + "\n\n" + phrases[8] +
		"t.me/" + configuration.BotUsername + "?start=" + user.Token, keyboard)
	delete(pending, message.From.ID)
}

func usersJoined(users *[]tgbotapi.User){
	for _, val := range *users{
		var user User
		err := db.Collection("users").FindOne(bson.M{"telegramid": val.ID}, &user)
		if err != nil {
			continue
		}
		user.IsJoined = true
		err = db.Collection("users").Save(&user)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("User %v joined", val.ID)
	}
}

func userLeft(u *tgbotapi.User){
	var user User
	err := db.Collection("users").FindOne(bson.M{"telegramid": u.ID}, &user)
	if err != nil {
		return
	}
	user.IsJoined = false
	err = db.Collection("users").Save(&user)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("User %v left", u.ID)
}