/////////////////////////////////////////////////////
// Based on https://github.com/brian1917/vcodeHMAC //
/////////////////////////////////////////////////////

package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultAuthScheme = "VERACODE-HMAC-SHA-256"

func generateHeader(uri, method, apiKeyID, apiKeySecret string) (string, error) {

	host, err := getHost(uri)
	if err != nil {
		return "Unable to get Host", err
	}
	path, err := getPathParams(uri)
	if err != nil {
		return "Unable to get path params", err
	}

	signingData := formatSigningData(apiKeyID, host, path, method)
	timestamp := getCurrentTimestamp()
	nonce := generateNonce()
	authScheme := defaultAuthScheme
	signature, err := createSignature(authScheme, apiKeySecret, signingData, timestamp, nonce)
	if err != nil {
		return "Error creating signature", err
	}
	return formatHeader(authScheme, apiKeyID, timestamp, nonce, signature), nil
}

func createSignature(authScheme string, apiKeySecret string, signingData string, timestamp int64, nonce string) (string, error) {
	if authScheme == defaultAuthScheme {
		signature := hmacSig(apiKeySecret, signingData, timestamp, nonce)
		return signature, nil
	}
	return "", errors.New("unsupported auth scheme")
}

func hmacSig(apiKeySecret string, signingData string, timestamp int64, nonce string) string {

	timeString := strconv.FormatInt(timestamp, 10)
	apiKeySecDecoded, _ := hex.DecodeString(apiKeySecret)
	nonceDecoded, _ := hex.DecodeString(nonce)

	h := hmac.New(sha256.New, apiKeySecDecoded)
	h.Write(nonceDecoded)
	keyNonce := h.Sum(nil)

	h = hmac.New(sha256.New, keyNonce)
	h.Write([]byte(timeString))
	keyDate := h.Sum(nil)

	h = hmac.New(sha256.New, keyDate)
	h.Write([]byte("vcode_request_version_1"))
	signatureKey := h.Sum(nil)

	h = hmac.New(sha256.New, signatureKey)
	h.Write([]byte(signingData))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

func formatSigningData(apiKeyID string, host string, url string, method string) string {
	apiKeyIDLower := strings.ToLower(apiKeyID)
	hostName := strings.ToLower(host)
	method = strings.ToUpper(method)

	return fmt.Sprintf("id=%s&host=%s&url=%s&method=%s", apiKeyIDLower, hostName, url, method)
}

func formatHeader(authScheme string, apiKeyID string, timestamp int64, nonce string, signature string) string {
	return fmt.Sprintf("%s id=%s,ts=%d,nonce=%s,sig=%s", authScheme, apiKeyID, timestamp, nonce, signature)
}

func getCurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func generateNonce() string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
}

func getHost(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

func getPathParams(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	if len(u.RawQuery) > 0 {
		return fmt.Sprintf("%s?%s", u.Path, u.RawQuery), nil
	}
	return string(u.Path), nil
}
