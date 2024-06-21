package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Resp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func SendToBark(host, deviceKey, jsonPayload, key, iv string) error {
	base64Ciphertext, err := NewAesCbc(key, iv).Encrypt(jsonPayload)
	if err != nil {
		fmt.Printf("Failed to encrypt: %v", err)
		return err
	}

	data := url.Values{}
	data.Set("iv", iv)
	data.Set("ciphertext", base64Ciphertext)

	apiUrl, _ := url.JoinPath(host, deviceKey)

	resp, err := http.PostForm(apiUrl, data)
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return err

	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	fmt.Println("Response :", string(content))

	if err != nil || !AssetResp(content) {
		return errors.New("failed to read response")
	}

	return nil
}

func SendToBarkDirectly(host, deviceKey, jsonPayload string) error {
	apiUrl, _ := url.JoinPath(host, deviceKey)
	resp, err := http.Post(apiUrl, "application/json; charset=utf-8", strings.NewReader(jsonPayload))
	if err != nil {
		fmt.Printf("Failed to send request: %v", err)
		return err

	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	fmt.Println("Response :", string(content))

	if err != nil || !AssetResp(content) {
		return errors.New("failed to read response")
	}

	return err
}

func AssetResp(result []byte) bool {
	var r Resp
	_ = json.Unmarshal(result, &r)
	if r.Code == 200 {
		return true
	}

	return false
}
