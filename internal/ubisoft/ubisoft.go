package ubisoft

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eliassebastian/gor6-cron/internal/config"
	"log"
	"net/http"
	"time"
)

const (
	SESSIONSURL = "https://public-ubiservices.ubi.com/v3/profiles/sessions"
	USERNAME    = "test"
	PASS        = "test"
)

type UbisoftInfo struct {
	//client     		*http.Client
	maxRetries    int8
	sessionStart  time.Time
	sessionPeriod time.Duration
	sessionExpiry time.Time
	sessionKey    string
	spaceID       string
	sessionTicket string
}

var Ubisoft *UbisoftInfo

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

	if err != nil {
		log.Println(res, err)
	}

	log.Println("SUCCESS:  ", res.StatusCode, res, err, res.Body)

	defer res.Body.Close()

	var result map[string]interface{}

	err2 := json.NewDecoder(res.Body).Decode(&result)

	if err2 != nil {
		log.Fatalln(err)
	}

	fmt.Println(result)

	return client, nil
}

func refreshConn() error {
	return nil
}
