package main

import(
	"github.com/segmentio/ksuid"
)

func generateToken() string{
	return ksuid.New().String()
}
