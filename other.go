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

func submitEth(message *tgbotapi.Message) {
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
	//sendMessage(message.Chat.ID, phrases[7] + message.Text + "\n\n" + phrases[8] +
	//	"t.me/" + configuration.BotUsername + "?start=" + user.Token, keyboard)
	//delete(pending, message.From.ID)
	sendMessage(message.Chat.ID, phrases[16], nil)
	pending[message.From.ID] = 3
}

func submitEmail(message *tgbotapi.Message) {
	var user User
	//db.First(&user, "id = ?", message.From.ID)
	err := db.Collection("users").FindOne(bson.M{"telegramid": message.From.ID}, &user)
	if err != nil {
		log.Panic(err)
	}

	if ok := verifyEmail(message.Text); ok{
		user.Email = message.Text
		user.IsVerified = true
		user.StakesJoining += configuration.StakesPerJoin
		user.StakesTotal += configuration.StakesPerJoin
		//db.Save(&user)
		err = db.Collection("users").Save(&user)
		if err != nil {
			log.Panic(err)
		}
		if user.InvitedBy != 0{
			var user2 User
			err := db.Collection("users").FindOne(bson.M{"telegramid": user.InvitedBy}, &user2)
			if err != nil {
				log.Panic(err)
			}
			user2.RefCount++
			user2.StakesRef += configuration.StakesPerRef
			user2.StakesTotal += configuration.StakesPerRef
			err = db.Collection("users").Save(&user2)
			if err != nil {
				log.Panic(err)
			}
		}
		sendMessage(message.Chat.ID, phrases[7] + user.EthAddress + "\n" + phrases[17] + user.Email + "\n" +
			phrases[8] + "t.me/" + configuration.BotUsername + "?start=" + user.Token, keyboard)
		delete(pending, message.From.ID)
	}else {
		sendMessage(message.Chat.ID, phrases[18], nil)
	}

}

func usersJoined(users *[]tgbotapi.User){
	for _, val := range *users{
		var user User
		err := db.Collection("users").FindOne(bson.M{"telegramid": val.ID}, &user)
		if err != nil {
			user := User{
				TelegramID: val.ID,
				Username: val.FirstName,
				Token: generateToken(),
				IsJoined: true,
			}
			err = db.Collection("users").Save(&user)
			if err != nil {
				log.Panic(err)
			}
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
	user.StakesRef = 0
	user.StakesJoining = 0
	user.StakesTotal = 0
	err = db.Collection("users").Save(&user)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("User %v left", u.ID)
}