package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid" // go get github.com/google/uuid
)

const (
	baseURL   = "https://api-sandbox.doku.com"
	clientID  = "YOUR-DOKU-CLIENT-ID"
	secretKey = "YOUR-SECRET-KEY-DOKU"
)

// generateDigest = SHA-256(body) lalu base64
func generateDigest(jsonBody string) string {
	h := sha256.Sum256([]byte(jsonBody))
	return base64.StdEncoding.EncodeToString(h[:])
}

// generateSignature = HMAC-SHA256 dari komponen signature
func generateSignature(requestID, timestamp, target, digest string) string {
	component := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		clientID, requestID, timestamp, target, digest,
	)

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(component))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return "HMACSHA256=" + sig
}

func main() {
	target := "/doku-virtual-account/v2/payment-code"

	body := `{
		"order": {
			"invoice_number": "INV-TEST-0001",
			"amount": 150000
		},
		"virtual_account_info": {
			"billing_type": "FIX_BILL",
			"expired_time": 60,
			"reusable_status": false,
			"info1": "Merchant Demo Store"
		},
		"customer": {
			"name": "Test User",
			"email": "test@example.com"
		}
	}`

	requestID := uuid.NewString()
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	digest := generateDigest(body)
	signature := generateSignature(requestID, timestamp, target, digest)

	req, err := http.NewRequest("POST", baseURL+target, bytes.NewBufferString(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Request-Id", requestID)
	req.Header.Set("Request-Timestamp", timestamp)
	req.Header.Set("Signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(respBody))
}