package main

import(
	"github.com/segmentio/ksuid"
	"github.com/badoux/checkmail"
	"fmt"
	"log"
)

func generateToken() string{
	return ksuid.New().String()
}

func verifyEmail(email string) bool{
	log.Print(email)
	err := checkmail.ValidateFormat(email)
	if err != nil {
		fmt.Println(err)
		return false
	}
	err = checkmail.ValidateHost(email)
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		fmt.Printf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
		return false
	}else if err != nil{
		log.Print(err)
		return false
	}
	return true
}