package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"strconv"
	"log"
	"gopkg.in/mgo.v2/bson"
)

func start(message *tgbotapi.Message) {
	fields := strings.Fields(message.Text)
	if len(fields) == 1 {
		var user User
		//db.First(&user, "id = ?", message.From.TelegramID)
		err := db.Collection("users").FindOne(bson.M{"telegramid": message.From.ID}, &user)
		if err != nil {
			user := User{
				TelegramID: message.From.ID,
				Username:   message.From.FirstName,
				Token:      generateToken(),
			}
			err = db.Collection("users").Save(&user)
			if err != nil {
				log.Panic(err)
			}

		}
		//else {
		//	  sendMessage(message.Chat.ID, phrases[10], keyboard)
		//}
		sendMessage(message.Chat.ID, phrases[0]+message.From.FirstName+phrases[1], keyboard)

	} else if len(fields) == 2 {
		var user User
		//db.Find(&user, "token = ?", fields[1])
		err := db.Collection("users").FindOne(bson.M{"token": fields[1]}, &user)
		if err != nil {
			user = User{}
			//db.First(&user, "id = ?", message.From.ID)
			err = db.Collection("users").FindOne(bson.M{"telegramid": message.From.ID}, &user)
			if err != nil {
				user := User{
					TelegramID: message.From.ID,
					Username:   message.From.FirstName,
					Token:      generateToken(),
				}
				//db.Create(&user)
				err = db.Collection("users").Save(&user)
				if err != nil {
					log.Panic(err)
				}
			}
			//else {
			//	sendMessage(message.Chat.ID, phrases[10], keyboard)
			//	return
			//}
			sendMessage(message.Chat.ID, phrases[0]+message.From.FirstName+phrases[1], keyboard)

		} else {
			user2 := User{}
			//db.First(&user2, "id = ?", message.From.ID)
			err := db.Collection("users").FindOne(bson.M{"telegramid": message.From.ID}, &user2)
			if err == nil {
				//sendMessage(message.Chat.ID, phrases[10], keyboard)
				sendMessage(message.Chat.ID, phrases[0]+message.From.FirstName+phrases[1], keyboard)

			} else {
				//user.RefCount++
				////db.Save(&user)
				//err = db.Collection("users").Save(&user)
				//if err != nil {
				//	log.Panic(err)
				//}

				user2 = User{
					TelegramID: message.From.ID,
					Username:   message.From.FirstName,
					Token:      generateToken(),
					InvitedBy:  user.TelegramID,
				}

				//db.Create(&user2)
				err = db.Collection("users").Save(&user2)
				if err != nil {
					log.Panic(err)
				}
				sendMessage(message.Chat.ID, phrases[0]+message.From.FirstName+phrases[1], keyboard)

			}
		}

	}

}

func cancel(message *tgbotapi.Message) {
	delete(pending, message.From.ID)
	sendMessage(message.Chat.ID, phrases[6], keyboard)
}

func editJoin(query *tgbotapi.CallbackQuery) {
	log.Printf("[%s] %s", query.From.FirstName, "clicked Join")
	editMessage(query.Message.Chat.ID, query.Message.MessageID, phrases[2])
}

func editSubmit(query *tgbotapi.CallbackQuery) {
	log.Printf("[%s] %s", query.From.FirstName, "clicked Sumbit")
	var user User
	//db.First(&user, "id = ?", query.From.ID)
	err := db.Collection("users").FindOne(bson.M{"telegramid": query.From.ID}, &user)
	if err != nil {
		log.Panic(err)
	}
	if !user.IsJoined {
		editMessage(query.Message.Chat.ID, query.Message.MessageID, phrases[15])
		return
	}

	if user.EthAddress == "" {
		editMessage(query.Message.Chat.ID, query.Message.MessageID, phrases[3])
		pending[query.From.ID] = 1
	} else if user.Email == "" {
		editMessage(query.Message.Chat.ID, query.Message.MessageID, phrases[16])
		pending[query.From.ID] = 3
	} else {
		editMessage(query.Message.Chat.ID, query.Message.MessageID, phrases[11])
	}

}

func editCheck(query *tgbotapi.CallbackQuery) {
	log.Printf("[%s] %s", query.From.FirstName, "clicked Check")
	var user User
	//db.First(&user, "id = ?", query.From.ID)
	err := db.Collection("users").FindOne(bson.M{"telegramid": query.From.ID}, &user)
	if err != nil {
		log.Panic(err)
	}
	text := phrases[8] + "t.me/" +
		configuration.BotUsername + "?start=" + user.Token + "\n\n" +
		phrases[9] + strconv.Itoa(user.RefCount) + "\n" + phrases[19] + strconv.Itoa(user.StakesTotal)
	//if user.IsJoined && user.EthAddress != ""{
	//	text += "\n"
	//}
	editMessage(query.Message.Chat.ID, query.Message.MessageID, text)
}
