package main

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
	core "github.com/miracl/oidc-samples/integration-tests/crypto"
	"github.com/miracl/oidc-samples/integration-tests/crypto/BN254CX"
)

func TestAuthFail(t *testing.T) {

}

func TestAuth(t *testing.T) {
	// t.Parallel()
	name := uuid.New().String()
	userID := fmt.Sprintf("%v@example.com", name)
	deviceName := "The device of " + name
	pin := 1111

	identity, err := Register(userID, deviceName, pin)
	if err != nil {
		t.Error("Error registering an identity", err)
	}

	// fmt.Println("kor", identity)

	// fmt.Println("userId", userID)
	// fmt.Println("name", name)
	// fmt.Println("deviceName", deviceName)
	// fmt.Println("mpinid", hex.EncodeToString(identity.MPinID))
	// fmt.Println("dtas", identity.DTAs)
	// fmt.Println("token", hex.EncodeToString(identity.Token))

	// identity := Identity{}
	// identity.MPinID = hex2bytes("7b22696174223a313537353837393839362c22757365724944223a2273746f79616e2e6b69726f76406d697261636c2e636f6d222c22634944223a2235633438666333612d386530612d343532662d396230622d366264386566643530353737222c2273616c74223a224c694d496e3458385a6b626547365137435561624a41222c2276223a352c2273636f7065223a5b2261757468225d2c22647461223a5b5d2c227674223a227076227d")
	// identity.Token = hex2bytes("04162c447c631e1c2d1b886de5e4390a8b3a5f9ccad2788033c952f6b431b7f4e41fcbb44251450784702c2e4f4f9055414dfdf0397b9fc3a9ca534e2a901e5189")
	// identity.DTAs = "WyJmZDlmYTY3NzhmYzY3N2RhZDQwZTE0OGY5NTM2YmUwN2Q0MjBjNTlmNWZjNmZkNGYyYWZjZjVhMWQ0MDM1MDdhIiwiODZjMThmNTEzMDVhMGRlMzQzNGM3NzliOWJlNWUxZDgzZjg4NzlmZmRjYjU0NGZjYzQ3MjEyNDJmODczZTYxYSJd"
	// userID := "stoyan.kirov@miracl.com"

	Authenticate(identity, userID, pin)
}

func Register(userID string, deviceName string, pin int) (identity Identity, err error) {
	const G1S = (BN254CX.MFS * 2) + 1 //https://github.com/miracl/core/blob/master/go/TestALL.go#L718

	//authorize
	authorizeResponse, err := authorize(userID, "6157wx51pm48u", "http://localhost:8000/login", "", "openid")
	if err != nil {
		return Identity{}, err
	}

	//activate/initiate
	cvResponse, err := customVerify(userID, deviceName)
	if err != nil {
		return Identity{}, err
	}

	//rps/v2/user
	qrURL, err := url.Parse(authorizeResponse.QRURL)
	if err != nil {
		return Identity{}, err
	}
	regResponse, err := register(userID, deviceName, qrURL.Fragment, cvResponse.ActivateToken)
	if err != nil {
		return Identity{}, err
	}

	//signature
	sigResponse, err := signature(regResponse.MPinID, regResponse.RegOTT)
	if err != nil {
		return Identity{}, err
	}

	//dta/ID
	csResponse, err := clientSecret(sigResponse.CS2URL)

	//Combine bot client secrets
	secret := make([]byte, G1S)
	BN254CX.MPIN_RECOMBINE_G1(hex2bytes(sigResponse.ClientSecretShare), hex2bytes(csResponse.ClientSecret), secret)

	//First extract pin from the combine client secret, in order to get the token
	token := secret
	exitCode := BN254CX.MPIN_EXTRACT_PIN(BN254CX.HASH_TYPE, hex2bytes(regResponse.MPinID), pin, token)
	if exitCode != 0 {
		return Identity{}, newCryptoError(exitCode)
	}

	return Identity{
		MPinID: hex2bytes(regResponse.MPinID),
		Token:  token,
		DTAs:   sigResponse.DTAs,
	}, nil
}

func Authenticate(identity Identity, userID string, pin int) error {
	//we'll need data from this response (wid) for pass2 request
	//otherwise, for some reason, if wid is not passed, the authOTT, that you'll need later is not stored in redis
	authorizeResponse, err := authorize(userID, "6157wx51pm48u", "http://localhost:8000/login", "", "openid")
	if err != nil {
		return err
	}
	qrURL, err := url.Parse(authorizeResponse.QRURL)
	if err != nil {
		return err
	}

	//Get pass1 proof from the token and pin (this is the U param in /pass1)
	rand := core.NewRAND()
	X := make([]byte, 32)
	S := make([]byte, 65)
	U := make([]byte, 65)
	exitCode := BN254CX.MPIN_CLIENT_1(BN254CX.HASH_TYPE, 0, identity.MPinID, rand, X, pin, identity.Token, S, U, nil, nil)
	if exitCode != 0 {
		return newCryptoError(exitCode)
	}

	//rps/v2/pass1
	p1Response, err := pass1(identity, U, "oidc")
	if exitCode != 0 {
		return newCryptoError(exitCode)
	}

	//Get V (used in /pass2) param using Y param from the pass1 response
	exitCode = BN254CX.MPIN_CLIENT_2(X, hex2bytes(p1Response.Y), S)
	if exitCode != 0 {
		return newCryptoError(exitCode)
	}

	//rps/v2/pass2
	p2Response, err := pass2(identity, S, qrURL.Fragment)
	if err != nil {
		return err
	}

	//rps/v2/authenticate
	authResponse, err := authenticate(p2Response.AuthOTT)
	if err != nil {
		return err
	}

	fmt.Println("auth", authResponse)

	return nil
}
