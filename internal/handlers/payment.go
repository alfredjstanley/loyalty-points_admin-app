package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"wac-offline-payment/internal/models"
	"wac-offline-payment/internal/repository"
)

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type PaymentRequest struct {
	UserMobileNumber     string  `json:"user_mobile_number"`
	MerchantMobileNumber string  `json:"merchant_mobile_number"`
	Amount               float64 `json:"amount"`
	InvoiceID            string  `json:"invoice_id"`
	PaymentMode          string  `json:"payment_mode"`
}

type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []any  `json:"data"`
}

func HandlePoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming JSON payload
	var paymentReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check for duplicate invoice_id with status "SUCCESS"
	exists, err := repository.InvoiceExists(paymentReq.InvoiceID)
	if err != nil {
		http.Error(w, "Failed to check invoice existence: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Duplicate entry: a successful transaction with this invoice_id already exists", http.StatusConflict)
		return
	}

	// Authenticate and get token
	authToken, err := authenticate()
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Make payment request
	paymentRes, err := makePaymentRequest(paymentReq, authToken)
	logStatus := "SUCCESS"
	if err != nil {
		logStatus = "FAILURE"
	}

	// Save log to MongoDB
	logEntry := models.Log{
		UserPhone:     paymentReq.UserMobileNumber,
		MerchantPhone: paymentReq.MerchantMobileNumber,
		Amount:        paymentReq.Amount,
		InvoiceID:     paymentReq.InvoiceID,
		PaymentMode:   paymentReq.PaymentMode,
		Status:        logStatus,
		Response:      paymentRes,
	}
	saveErr := repository.SaveLog(logEntry)
	if saveErr != nil {
		http.Error(w, "Failed to save log: "+saveErr.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to the user
	if err != nil {
		http.Error(w, "Payment processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentRes)
}

func authenticate() (string, error) {
	authURL := os.Getenv("AUTH_URL")
	authUsername := os.Getenv("AUTH_USERNAME")
	authPassword := os.Getenv("AUTH_PASSWORD")

	if authURL == "" || authUsername == "" || authPassword == "" {
		return "", errors.New("AUTH_URL, AUTH_USERNAME, or AUTH_PASSWORD is not set in the environment")
	}

	authPayload := map[string]string{
		"username": authUsername,
		"password": authPassword,
	}

	payloadBytes, _ := json.Marshal(authPayload)
	req, _ := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error during authentication: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Capture the response body for debugging
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to authenticate")
	}

	var authRes AuthResponse
	if err := json.Unmarshal(body, &authRes); err != nil {
		log.Printf("Error unmarshalling auth response: %v", err)
		return "", err
	}

	if !authRes.Success {
		log.Printf("Auth API response indicates failure: %s", authRes.Message)
		return "", errors.New(authRes.Message)
	}

	return authRes.Data.Token, nil
}

func makePaymentRequest(paymentReq PaymentRequest, token string) (*PaymentResponse, error) {
	paymentURL := os.Getenv("PAYMENT_URL")
	if paymentURL == "" {
		return nil, errors.New("PAYMENT_URL is not set in the environment")
	}

	payloadBytes, _ := json.Marshal(paymentReq)

	req, _ := http.NewRequest(http.MethodPost, paymentURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making payment request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("payment request failed with status %d", resp.StatusCode))
	}

	var paymentRes PaymentResponse
	if err := json.Unmarshal(body, &paymentRes); err != nil {
		log.Printf("Error unmarshalling payment response: %v", err)
		return nil, err
	}

	if !paymentRes.Success {
		log.Printf("Payment API response indicates failure: %s", paymentRes.Message)
		return nil, errors.New(paymentRes.Message)
	}

	return &paymentRes, nil
}

// RenderForm renders the HTML form
func RenderForm(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templatesDir + "/form.html")
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}
