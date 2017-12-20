package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomID gives back a random ID that can be used
// when no id was defined for a module
func RandomID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetJSON is used to fetch data for a given url and getting parsed into a given interface
func GetJSON(urlStr string, target interface{}) error {
	req, _ := http.NewRequest("GET", URLEncoded(urlStr), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// PostJSON is used to post data as JSON to a server
func PostJSON(urlStr string, data interface{}, expectedStatus int) (*http.Response, error) {
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if expectedStatus != resp.StatusCode {
		return resp, fmt.Errorf("Unexpected StatusCode, expected %v got %v", expectedStatus, resp.StatusCode)
	}

	return resp, nil
}

// URLEncoded encodes a string like Javascript's encodeURIComponent()
func URLEncoded(str string) string {
	str = strings.Replace(str, "'", "%27", -1)
	str = strings.Replace(str, " ", "%20", -1)
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}

	return u.String()
}
