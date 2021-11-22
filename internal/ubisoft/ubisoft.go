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
	UBIURL = "https://public-ubiservices.ubi.com/v3/profiles/sessions"
)

func establishConn() error {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodPost, UBIURL, nil)
	if err != nil {
		log.Println(req, err)
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json"},
		"Ubi-AppId":     []string{"39baebad-39e5-4552-8c25-2c9b919064e2"},
		"Authorization": []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "test", "test"))))},
		"Connection":    []string{"keep-alive"},
	}

	res, err := client.Do(req)

	if err != nil {
		log.Println(res, err)
	}

	log.Println("SUCCESS:  ", res, err, res.Body)

	defer res.Body.Close()

	var result map[string]interface{}

	err2 := json.NewDecoder(res.Body).Decode(&result)

	if err2 != nil {
		log.Fatalln(err)
	}

	fmt.Println(result)

	return nil
}
