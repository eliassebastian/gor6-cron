package ubisoft

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eliassebastian/gor6-cron/internal/config"
	"github.com/eliassebastian/gor6-cron/internal/models"
	"log"
	"net/http"
	"time"
)

const (
	SESSIONSURL = "https://public-ubiservices.ubi.com/v3/profiles/sessions"
	USERNAME    = "gor6client@gmail.com"
	PASS        = "GoClientR6!2021"
)

var Ubisoft *models.UbisoftSession

func createClient() *http.Client {

	if config.Config != nil {
		return config.Config.Client
	}

	return &http.Client{
		Timeout: time.Second * 10,
	}
}

//TODO Accept Different User Details
func basicToken() string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", USERNAME, PASS))))
}

func Connect() {

}

func EstablishConn() (*http.Client, error) {

	client := createClient()

	req, err := http.NewRequest(http.MethodPost, SESSIONSURL, nil)
	if err != nil {
		log.Println(req, err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{"39baebad-39e5-4552-8c25-2c9b919064e2"},
		"Authorization": []string{basicToken()},
		"Connection":    []string{"keep-alive"},
	}

	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		//TODO Error Handling
		return nil, err
	}

	if res.StatusCode == 400 || res.StatusCode == 401 {
		return nil, nil
	}

	Ubisoft = &models.UbisoftSession{
		Retries:       0,
		MaxRetries:    5,
		RetryTime:     10,
		SessionStart:  time.Now().UTC(),
		SessionPeriod: 175,
	}

	err2 := json.NewDecoder(res.Body).Decode(Ubisoft)
	if err2 != nil {
		log.Fatalln(err)
	}

	fmt.Println(Ubisoft)
	return client, nil
}

func refreshConn() error {
	return nil
}
