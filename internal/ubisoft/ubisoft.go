package ubisoft

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eliassebastian/gor6-cron/internal/pubsub"
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
	ctx              context.Context
	cancel           context.CancelFunc
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

func (c *UbisoftConfig) Connect(ctx context.Context, pro *pubsub.Producer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.ctx = ctx
	c.cancel = cancel

	var backoffSchedule = []time.Duration{
		5 * time.Second,
		10 * time.Second,
		15 * time.Second,
		30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SESSIONSURL, nil)
	if err != nil {
		log.Println(req, err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{c.appId},
		"Authorization": []string{c.appAuthorisation},
		"Connection":    []string{"keep-alive"},
	}

	for i, backoff := range backoffSchedule {
		log.Println("Running Client Fetch Iteration:", i)
		res, err := c.client.Do(req)
		if err != nil {
			log.Printf("Error with HTTP Response: %s", err)
			cancel()
			break
		}

		if res.StatusCode == 200 {
			err := pro.NewMessage(ctx, &res.Body)
			if err == nil {
				cancel()
				break
			}
		}

		log.Println("Waiting on Client Fetch Iteration:", i+1)
		res.Body.Close()
		time.Sleep(backoff)
	}

}

func connection(ctx context.Context, config *UbisoftConfig) {

}

func (c *UbisoftConfig) Connect2(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SESSIONSURL, nil)
	if err != nil {
		log.Println(req, err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{c.appId},
		"Authorization": []string{c.appAuthorisation},
		"Connection":    []string{"keep-alive"},
	}

	for i := 0; i < 6; i++ {
		res, errhttp := c.client.Do(req)
		if errhttp != nil {
			return errhttp
		}

		if res.StatusCode == 200 {
			pubsub.NewMessage(res.Body)

		}

		res.Body.Close()
		time.Sleep(10 * time.Second)
	}

	res, errhttp := c.client.Do(req)
	defer res.Body.Close()

	if res.StatusCode == 200 {
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

	//TODO Retry Connection

	return errors.New(fmt.Sprintf(""))
}

func (c *UbisoftConfig) Stop() {
	c.cancel()
}
