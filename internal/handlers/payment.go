package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type PaymentRequest struct {
	UserMobileNumber     string `json:"user_mobile_number"`
	MerchantMobileNumber string `json:"merchant_mobile_number"`
	Amount               int    `json:"amount"`
	InvoiceID            string `json:"invoice_id"`
	PaymentMode          string `json:"payment_mode"`
}

type PaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []any  `json:"data"`
}

// HandlePoints processes the `/api/points` request
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

	// Step 1: Authenticate and get token
	authToken, err := authenticate()
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Step 2: Make payment request
	paymentRes, err := makePaymentRequest(paymentReq, authToken)
	if err != nil {
		http.Error(w, "Payment processing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Respond to the user with the payment response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentRes)
}

func authenticate() (string, error) {
	// Authentication API details
	authURL := "https://olopo-dev.webc.in/api/auth/login"
	authPayload := map[string]string{
		"username": "admin@olopo.com",
		"password": "UJYHdVHeOU19S3i",
	}

	payloadBytes, _ := json.Marshal(authPayload)
	req, _ := http.NewRequest(http.MethodPost, authURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to authenticate")
	}

	var authRes AuthResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &authRes); err != nil {
		return "", err
	}

	if !authRes.Success {
		return "", errors.New(authRes.Message)
	}

	return authRes.Data.Token, nil
}

func makePaymentRequest(paymentReq PaymentRequest, token string) (*PaymentResponse, error) {
	paymentURL := "https://olopo-dev.webc.in/api/payments/offline/complete"

	payloadBytes, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest(http.MethodPost, paymentURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("payment request failed")
	}

	var paymentRes PaymentResponse
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &paymentRes); err != nil {
		return nil, err
	}

	if !paymentRes.Success {
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

// HandlePaymentSubmission processes the form submission
func HandlePaymentSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract form data
	userPhone := r.FormValue("userPhone")
	merchantPhone := r.FormValue("merchantPhone")
	amount := r.FormValue("amount")
	invoiceID := r.FormValue("invoiceId")
	paymentMode := r.FormValue("paymentMode")

	// Log the submitted data
	log.Printf("Received Payment Details:\nUser Phone: %s\nMerchant Phone: %s\nAmount: %s\nInvoice ID: %s\nPayment Mode: %s",
		userPhone, merchantPhone, amount, invoiceID, paymentMode)

	// Simulate a success response
	message := fmt.Sprintf("Payment processed successfully!<br>User Phone: %s<br>Merchant Phone: %s<br>Amount: %s<br>Invoice ID: %s<br>Payment Mode: %s",
		userPhone, merchantPhone, amount, invoiceID, paymentMode)

	// Send response back to the user
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
