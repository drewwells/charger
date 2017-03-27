package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sfreiberg/gotwilio"
)

type Cred struct {
	Accountsid string
	Authtoken  string
	From       string
	To         string
}

var (
	accountSid, authToken string
)

type config struct {
	Live, Test Cred
}

var parseCred = func(c config) Cred {
	return c.Live
}

func main() {

	bs, err := ioutil.ReadFile("keys.toml")
	if err != nil {
		log.Fatal("missing keys.toml config")
	}

	var c config
	err = toml.Unmarshal(bs, &c)
	if err != nil {
		log.Fatal("failed to unmarshal config", err)
	}

	stations, err := fetchStatus()
	if err != nil {
		log.Fatal(err)
	}

	if len(stations) > 0 {
		m := fmt.Sprintf("Available! %s", strings.Join(stations, ", "))
		parseCred(c).sendSMS(m)
	}
}

type safeStation Station

type Station struct {
	Status    string
	Available int
}

func (s *Station) UnmarshalJSON(bs []byte) error {
	var ss safeStation
	err := json.Unmarshal(bs, &ss)
	if err == nil {
		s.Status = ss.Status
		s.Available = ss.Available
	}
	return nil
}

type safeCharger struct {
	Name     string
	Stations map[string]Station
}

type Charger struct {
	ID        int `json:"id"`
	Name      string
	Available int
}

func (c *Charger) UnmarshalJSON(bs []byte) error {
	var s safeCharger

	err := json.Unmarshal(bs, &s)
	if err == nil {
		c.Name = s.Name
		for _, v := range s.Stations {
			c.Available = c.Available + v.Available
		}
	}

	return nil
}

type SemaResp struct {
	Data map[string]Charger `json:"aaData"`
}

var stations map[string]int

func fetchStatus() ([]string, error) {

	vals, err := url.ParseQuery("action=locationSearch&address=78759&pseudoParam= 1490292687569")
	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm("https://network.semaconnect.com/get_data.php", vals)
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var r SemaResp
	err = json.Unmarshal(bs, &r)
	if err != nil {
		return nil, err
	}

	// Filter stonebridge accounts
	var sel []string
	for _, v := range r.Data {
		if !strings.Contains(v.Name, "Stonebridge") {
			continue
		}
		if v.Available == 0 {
			continue
		}
		sel = append(sel, v.Name)
	}

	return sel, nil
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
