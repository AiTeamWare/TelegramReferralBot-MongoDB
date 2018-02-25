package main

import "github.com/go-bongo/bongo"

type Config struct {
	BotToken, BotUsername, Address, DBName string
}

type User struct {
	bongo.DocumentBase `bson:",inline"`
	TelegramID                  int
	Username, Token, EthAddress string
	RefCount                    int
	IsJoined					bool
}
