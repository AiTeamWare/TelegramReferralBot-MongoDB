package main

import (
	"reflect"
	"strings"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMessage(chatId int64, text string, keyboard interface{}){
	msg := tgbotapi.NewMessage(chatId, text)
	//msg.ParseMode = tgbotapi.ModeMarkdown
	typeOfKeyboard := reflect.TypeOf(keyboard)
	if typeOfKeyboard == nil {
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{true, false}
	}else {
		typeOfKeyboardString := strings.Split(typeOfKeyboard.String(), ".")[1]
		switch typeOfKeyboardString {
		default:
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{true, false}
		case "ReplyKeyboardMarkup":
			msg.ReplyMarkup = keyboard
		case "InlineKeyboardMarkup":
			msg.ReplyMarkup = &keyboard
		}
	}

	_, err := bot.Send(msg)
	if err != nil {
		log.Print(err)
	}
	log.Printf("[Bot] SENT %s TO %v", msg.Text, msg.ChatID)
}

func editMessage(chatId int64, messageId int, text string){
	msg := tgbotapi.NewEditMessageText(chatId, messageId, text)
	//msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = &keyboard
	bot.Send(msg)
	log.Printf("[Bot] EDITED %s TO %v", msg.Text, msg.ChatID)
}