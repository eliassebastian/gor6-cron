package ubisoft

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eliassebastian/gor6-cron/internal/pubsub"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SESSIONSURL = "https://public-ubiservices.ubi.com/v3/profiles/sessions"
	APPID       = "39baebad-39e5-4552-8c25-2c9b919064e2"
	USERNAME    = "gor6client@gmail.com"
	PASS        = "GoClientR6!2021"
)

type UbisoftConfig struct {
	client  *http.Client
	ctx     context.Context
	cancel  context.CancelFunc
	session []byte
}

//type UbisoftSession struct {
//	SessionExpiry time.Time `json:"expiration"`
//	SessionKey    string    `json:"sessionKey"`
//	SpaceID       string    `json:"spaceId"`
//	SessionTicket string    `json:"ticket"`
//}

func basicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", USERNAME, PASS))))
}

func CreateConfig() *UbisoftConfig {
	return &UbisoftConfig{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func createSessionURL(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, errors.New("error creating session url")
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{APPID},
		"Authorization": []string{basicToken()},
		"Connection":    []string{"keep-alive"},
	}

	return req, nil
}

func fetchSessionData(ctx context.Context, client *http.Client, r *http.Request) []byte {
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
				bs, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println("error reading response body", err)
					res.Body.Close()
					return nil
				}
				res.Body.Close()
				return bs
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

	req, err := createSessionURL(ctx, SESSIONSURL)
	if err != nil {
		return err
	}

	sd := fetchSessionData(ctx, c.client, req)
	if sd == nil {
		log.Println("Fetch Session Data returned Nil")
		return errors.New("session fetched returned nil")
	}

	c.session = sd
	ke := p.NewMessage(ctx, sd)
	if ke != nil {
		return ke
	}

	return nil
}

func (c *UbisoftConfig) Stop() {
	c.cancel()
}
