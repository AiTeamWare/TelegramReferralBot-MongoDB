package main

import "github.com/go-bongo/bongo"

type Config struct {
	BotToken, BotUsername, Address, DBName string
	ChatID                                 int64
	StakesPerJoin, StakesPerRef            int
}

type User struct {
	bongo.DocumentBase `bson:",inline"`
	TelegramID                            int
	Username, Token, EthAddress, Email    string
	RefCount                              int
	IsJoined                              bool
	IsVerified                            bool
	InvitedBy                             int
	StakesJoining, StakesRef, StakesTotal int
}
