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
	client *http.Client
	appId  string
	ctx    context.Context
	cancel context.CancelFunc
	UbisoftSession
}

type UbisoftSession struct {
	SessionExpiry time.Time `json:"expiration"`
	SessionKey    string    `json:"sessionKey"`
	SpaceID       string    `json:"spaceId"`
	SessionTicket string    `json:"ticket"`
}

func basicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", USERNAME, PASS))))
}

func CreateConfig() *UbisoftConfig {
	return &UbisoftConfig{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		appId: "39baebad-39e5-4552-8c25-2c9b919064e2",
	}
}

func createSessionURL(ctx context.Context, url string, appId string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, errors.New("error creating session url")
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{appId},
		"Authorization": []string{basicToken()},
		"Connection":    []string{"keep-alive"},
	}

	return req, nil
}

func fetchSessionData(ctx context.Context, client *http.Client, r *http.Request) *UbisoftSession {
	bs := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		15 * time.Second,
		30 * time.Second,
	}

	for i, b := range bs {
		select {
		case <-ctx.Done():
			log.Println("Session Fetch Loop Context Done")
			return nil
		default:
			log.Println("Running Client Fetch Iteration:", i)
			res, err := client.Do(r)
			if err != nil {
				return nil
			}

			if res.StatusCode == 200 {
				var u *UbisoftSession
				err := json.NewDecoder(res.Body).Decode(&u)
				if err == nil {
					return u
				}
			}
			log.Println("Retrying Session:", i+1)
			res.Body.Close()
			time.Sleep(b)
		}
	}

	return nil
}

func (c *UbisoftConfig) Connect(ctx context.Context, p *pubsub.Producer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.ctx = ctx
	c.cancel = cancel

	req, err := createSessionURL(ctx, SESSIONSURL, c.appId)
	if err != nil {
		return err
	}

	data := fetchSessionData(ctx, c.client, req)
	if data == nil {
		log.Println("Fetch Session Data returned Nil")
		return errors.New("session fetched returned nil")
	}

	return nil
}

func (c *UbisoftConfig) Stop() {
	c.cancel()
}
