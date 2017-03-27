package main

import (
	"log"

	"github.com/sfreiberg/gotwilio"
)

type Cred struct {
	Accountsid string
	Authtoken  string
	From       string
	To         string
}

type config struct {
	Live, Test Cred
}

var parseCred = func(c config) Cred {
	return c.Live
}

func (c Cred) sendSMS(msg string) {
	twilio := gotwilio.NewTwilioClient(c.Accountsid, c.Authtoken)

	// always talking to the same person
	resp, ex, err := twilio.SendSMS(c.From, c.To, msg, "", "")
	if err != nil {
		log.Println("failed to send sms", err)
	}
	if ex != nil {
		log.Printf("twilio error: %#v", ex)
	}
	log.Printf("sms send %#v", resp)
}
