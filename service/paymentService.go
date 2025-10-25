package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"zarrinpal/models"
	"zarrinpal/repository"
)

const zarinpalRequestURL = "https://sandbox.zarinpal.com/pg/v4/payment/request.json"
const zarinpalVerifyURL = "https://sandbox.zarinpal.com/pg/v4/payment/verify.json"
const zarinpalGatewayURL = "https://sandbox.zarinpal.com/pg/StartPay/"

func RequestPayment(ctx context.Context, payload *models.PaymentRequestPayload, userid int64) (string, error) {
	merchantID := os.Getenv("ZARINPAL_MERCHANT_ID")
	callbackURL := "http://localhost:8080/payment/callback"

	zpReq := &models.ZarinpalRequest{
		MerchantID:  merchantID,
		Amount:      payload.Amount,
		Description: payload.Description,
		CallbackURL: callbackURL,
		Metadata:    payload.Metadata,
	}

	reqBytes, _ := json.Marshal(zpReq)
	req, err := http.NewRequestWithContext(ctx, "POST", zarinpalRequestURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var zarinpalResp models.ZarinpalResponse
	json.NewDecoder(resp.Body).Decode(&zarinpalResp)

	if zarinpalResp.Data.Code != 100 {
		return "", fmt.Errorf("zarinpal request failed with code: %d", zarinpalResp.Data.Code)
	}

	newPayment := &models.Payment{
		Amount:      payload.Amount,
		Description: payload.Description,
		Authority:   zarinpalResp.Data.Authority,
	}
	if err := repository.CreatePayment(ctx, newPayment, userid); err != nil {
		return "", err
	}

	return zarinpalGatewayURL + zarinpalResp.Data.Authority, nil
}

func VerifyPayment(ctx context.Context, authority, status string) (string, error) {
	if status != "OK" {
		repository.UpdatePayment(ctx, authority, "", "FAILED")
		return "", errors.New("payment canceled by user")
	}

	payment, err := repository.GetPaymentByAuthority(ctx, authority)
	if err != nil {
		return "", fmt.Errorf("payment not found for authority: %s", authority)
	}

	merchantID := os.Getenv("ZARINPAL_MERCHANT_ID")
	zpReq := &models.ZarinpalVerifyRequest{
		MerchantID: merchantID,
		Authority:  authority,
		Amount:     payment.Amount,
	}
	reqBytes, _ := json.Marshal(zpReq)

	req, err := http.NewRequestWithContext(ctx, "POST", zarinpalVerifyURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var verifyResp models.ZarinpalVerifyResponse
	json.NewDecoder(resp.Body).Decode(&verifyResp)

	if verifyResp.Data.Code == 100 || verifyResp.Data.Code == 101 {
		refID := strconv.FormatInt(verifyResp.Data.RefID, 10)
		repository.UpdatePayment(ctx, authority, refID, "COMPLETED")
		return refID, nil
	}

	repository.UpdatePayment(ctx, authority, "", "FAILED")
	return "", fmt.Errorf("payment verification failed with code: %d", verifyResp.Data.Code)
}
