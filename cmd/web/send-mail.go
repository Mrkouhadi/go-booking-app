package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mrkouhadi/go-booking-app/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	// an anonymous function that runs in the background
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(msg models.MailData) {

	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	// we haven't specified the template do this
	if msg.Template == "" {
		email.SetBody(mail.TextHTML, msg.Content)
	} else {
		// otherwise get the specified template from disk
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", msg.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		// since data is of type array of bytes ([]bytes) we need to convert it
		mailTmpl := string(data)
		//args:  string to be replace in, txt to be replaced, content, how many times
		msgToSend := strings.Replace(mailTmpl, "[%body%]", msg.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("EMAIL HAS BEEN SENT SUCCESFULLY")
	}
}
