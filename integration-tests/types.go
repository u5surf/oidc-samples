package main

import "fmt"

type Identity struct {
	MPinID []byte
	Token  []byte
	DTAs   string
}

type CryptoError struct {
	exitCode int
}

func newCryptoError(exitCode int) CryptoError {
	return CryptoError{
		exitCode: exitCode,
	}
}

func (err CryptoError) Error() string {
	return fmt.Sprintf("Crypto error exited with code %v", err.exitCode)
}

type CustomVerificationRequest struct {
	UserID     string `json:"userId"`
	DeviceName string `json:"deviceName"`
}

type CustomVerificationResponse struct {
	MPinID        string `json:"mpinId"`
	HashMPinID    string `json:"hashMPinId"`
	ActivateToken string `json:"actToken"`
	ExpireTime    int64  `json:"expireTime"`
}

type RegisterResponse struct {
	Active     bool   `json:"active"`
	AppID      string `json:"appId"`
	CustomerID string `json:"customerId"`
	ExpireTime int    `json:"expireTime"`
	MPinID     string `json:"mpinId"`
	NowTime    int    `json:"nowTime"`
	RegOTT     string `json:"regOTT"`
}

type RegisterRequest struct {
	DeviceName   string `json:"deviceName"`
	UserID       string `json:"userId"`
	WID          string `json:"wid"`
	ActivateCode string `json:"activateCode"`
}

type SignatureRequest struct {
	MPinID string `json:"-"`
	RegOTT string `json:"-"`
}

type SignatureResponse struct {
	ClientSecretShare string `json:"clientSecretShare"`
	CS2URL            string `json:"cs2url"`
	Curve             string `json:"curve"`
	DTAs              string `json:"dtas"`
}

type Header struct {
	Key, Value string
}

type AuthorizeResponse struct {
	AccessURL string `json:"accessURL"`
	QRURL     string `json:"qrURL"`
	WebOTT    string `json:"webOTT"`
}

type ClientSecretRequest struct {
	CS2URL string `json:"cs2url"`
}

type ClientSecretResponse struct {
	ClientSecret    string `json:"clientSecret"`
	DVSClientSecret string `json:"dvsClientSecret"`
	CreatedAt       int    `json:"createdAt"`
	Message         string `json:"message"`
	Version         string `json:"version"`
}

//Authentication

type Pass1Request struct {
	U      string   `json:"U"`
	MPinID string   `json:"mpin_id"`
	DTAs   string   `json:"dtas"`
	Scope  []string `json:"scope"`
}

type Pass1Response struct {
	Y string `json:"y"`
}

type Pass2Request struct {
	V      string `json:"V"`
	WID    string `json:"WID"`
	MPinID string `json:"mpin_id"`
}

type Pass2Response struct {
	AuthOTT string `json:"authOTT"`
}

type AuthResponse struct {
	WebOTT string `json:"webOTT"`
}

type AuthenticateRequest struct {
	AuthOTT string `json:"authOTT"`
}

type AuthenticateResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
