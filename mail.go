package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {

	http.HandleFunc("/_ah/mail/", incomingMail)
	http.handleFunc("/mailnotify", mail)
}

func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		log.Errorf(ctx, "Error reading body: %v", err)
		return
	}
	log.Infof(ctx, "Received mail: %v", b)
}

func mail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	stations, err := fetchStatus()
	if err != nil {
		log.Error("failed to retrieve stations")
		return
	}

	subject := "No chargers available"
	body := ""
	if len(stations) > 0 {
		subject = "Chargers Available!"
		body = fmt.Sprintf("Available! %s", strings.Join(stations, ", "))
	}

	msg := &mail.Message{
		Sender:  "Charger Notify <notify@charger-notify.appspot.com>",
		To:      []string{"drew.wells00@gmail.com"},
		Subject: subject,
		Body:    body,
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "Couldn't send email: %v", err)
	}
}
