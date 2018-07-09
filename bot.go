package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"github.com/go-bongo/bongo"
	"os"
	"io"
)

var (
	bot           *tgbotapi.BotAPI
	configuration Config
	phrases       map[int]string
	db            *bongo.Connection
	pending		  = make(map[int]int)
	keyboard tgbotapi.InlineKeyboardMarkup
)

func main() {
	initLog()
	initConfig()
	initStrings()
	initDB()
	initKeyboard()

	var err error
	bot, err = tgbotapi.NewBotAPI(configuration.BotToken)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				log.Print("It is command")
				command := update.Message.Command()
				switch command {
				case "start":
					log.Print("It is start")
					start(update.Message)
				case "cancel":
					log.Print("It is cancel")
					cancel(update.Message)
				}
			}else if update.Message.NewChatMembers != nil && update.Message.Chat.ID == configuration.ChatID{
				usersJoined(update.Message.NewChatMembers)
			}else if update.Message.LeftChatMember != nil && update.Message.Chat.ID == configuration.ChatID{
				userLeft(update.Message.LeftChatMember)
			}else {
				if pending[update.Message.From.ID] == 1{
					requestRepeat(update.Message)
				}else if pending[update.Message.From.ID] == 2{
					submitEth(update.Message)
				}else if pending[update.Message.From.ID] == 3{
					submitEmail(update.Message)
				}
			}
		}else if update.CallbackQuery != nil {
			log.Print("It is callback")
			switch update.CallbackQuery.Data {
			case "join":
				go editJoin(update.CallbackQuery)
			case "submitEth":
				go editSubmit(update.CallbackQuery)
			case "check":
				go editCheck(update.CallbackQuery)
			}
			bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{update.CallbackQuery.ID, "", false, "", 0})
		}

	}
}

func initLog() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func initConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	log.Print("First 10 bytes from config.json")
	log.Print(body[:10])
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	log.Print("First 10 bytes after trim")
	reader := bytes.NewReader(body)
	log.Print(body[:10])
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&configuration)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}

}

func initStrings() {
	file, err := os.Open("strings.json")
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	log.Print("First 10 bytes from strings.json")
	log.Print(body[:10])
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	log.Print("First 10 bytes after trim")
	reader := bytes.NewReader(body)
	log.Print(body[:10])
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(&phrases)
	if err != nil {
		log.Print("ERROR: ")
		log.Panic(err)
	}
}

func initDB() {
	var err error
	config := &bongo.Config{
		ConnectionString: configuration.Address,
		Database:         configuration.DBName,
	}
	db, err = bongo.Connect(config)
	if err != nil {
		log.Panic(err)
	}
	log.Print("Database connected")
}

func initKeyboard(){
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(phrases[12], "join"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(phrases[13], "submitEth"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(phrases[14], "check"),
		),
	)
}
