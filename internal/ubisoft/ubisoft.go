package ubisoft

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	SESSIONSURL = "https://public-ubiservices.ubi.com/v3/profiles/sessions"
	USERNAME    = "gor6client@gmail.com"
	PASS        = "GoClientR6!2021"
)

type UbisoftConfig struct {
	client           *http.Client
	appId            string
	appAuthorisation string
	UbisoftSession
}

type UbisoftSession struct {
	Retries       uint8
	MaxRetries    uint8
	RetryTime     uint8
	SessionStart  time.Time
	SessionPeriod uint16
	SessionExpiry time.Time `json:"expiration"`
	SessionKey    string    `json:"sessionKey"`
	SpaceID       string    `json:"spaceId"`
	SessionTicket string    `json:"ticket"`
}

//TODO Accept Different User Details
func basicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", USERNAME, PASS))))
}

func CreateConfig() *UbisoftConfig {
	return &UbisoftConfig{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		appId:            "39baebad-39e5-4552-8c25-2c9b919064e2",
		appAuthorisation: basicToken(),
	}
}

func (c *UbisoftConfig) Connect() error {
	req, err := http.NewRequest(http.MethodPost, SESSIONSURL, nil)
	if err != nil {
		log.Println(req, err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{c.appId},
		"Authorization": []string{c.appAuthorisation},
		"Connection":    []string{"keep-alive"},
	}

	res, errhttp := c.client.Do(req)
	defer res.Body.Close()

	if errhttp != nil {
		return errhttp
	}

	if res.StatusCode == 400 || res.StatusCode == 401 {
		//TODO Retry Connection
		return nil
	}

	c.UbisoftSession = UbisoftSession{
		Retries:       0,
		MaxRetries:    5,
		RetryTime:     10,
		SessionStart:  time.Now().UTC(),
		SessionPeriod: 175,
	}

	errdec := json.NewDecoder(res.Body).Decode(c.UbisoftSession)
	if errdec != nil {
		log.Fatalln(err)
	}

	fmt.Println(c)
	return nil
}
