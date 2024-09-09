package pkg_login

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

func postBase(url string, payload string, headers map[string]string) (resp *http.Response, err error) {
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return
	}
	for index, val := range headers {
		req.Header.Set(index, val)
	}

	return client.Do(req)
}

func getBase(requestUrl string, headers map[string]string) (resp *http.Response, err error) {
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}
	for index, val := range headers {
		req.Header.Set(index, val)
	}
	return client.Do(req)
}

func rand32Str() string {
	harsher := md5.New()
	harsher.Write([]byte(uuid.New().String()))
	hashBytes := harsher.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)
	if len(hashStr) > 32 {
		return hashStr[:32]
	}

	return hashStr
}
