package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RazorpayClient interface {
	CreateOrder(amount int64, currency, receipt string) (*RazorpayOrder, error)
	VerifySignature(orderID, paymentID, signature string) bool
	Refund(paymentID string) error
}

type razorpayClient struct {
	key    string
	secret string
}

func NewRazorpayClient(key, secret string) RazorpayClient {
	return &razorpayClient{
		key:    key,
		secret: secret,
	}
}

/////   create order  /////

type RazorpayOrder struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

func (r *razorpayClient) CreateOrder(
	amount int64,
	currency string,
	receipt string,
) (*RazorpayOrder, error) {

	body := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
		"receipt":  receipt,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(r.key, r.secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var order RazorpayOrder
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

///////  verify signature  /////////

func (r *razorpayClient) VerifySignature(
	orderID, paymentID, signature string,
) bool {

	payload := fmt.Sprintf("%s|%s", orderID, paymentID)

	h := hmac.New(sha256.New, []byte(r.secret))
	h.Write([]byte(payload))
	expected := hex.EncodeToString(h.Sum(nil))

	return expected == signature
}

func (r *razorpayClient) Refund(paymentID string) error {

    url := fmt.Sprintf(
        "https://api.razorpay.com/v1/payments/%s/refund",
        paymentID,
    )

    req, _ := http.NewRequest("POST", url, nil)
    req.SetBasicAuth(r.key, r.secret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    if resp.StatusCode >= 300 {
        return errors.New("refund failed")
    }

    return nil
}