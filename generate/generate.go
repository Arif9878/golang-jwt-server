package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Arif9878/golang-auth-service/assets"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
)

const (
	keyUsage = "sig"
)

type Keys struct {
	Keys []jwk.Key `json:"keys"`
}

func main() {
	keyData, err := assets.GetResources().ReadFile("private.pem")
	if err != nil {
		return
	}

	raw, e := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if e != nil {
		fmt.Printf("failed parse rsa: %s\n", err)
		return
	}

	key, err := jwk.New(raw)
	if err != nil {
		fmt.Printf("failed to create RSA key: %s\n", err)
		return
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		fmt.Printf("expected jwk.RSAPrivateKey, got %T\n", key)
		return
	}

	id := uuid.New()

	key.Set(jwk.KeyIDKey, id.String())
	key.Set(jwk.KeyUsageKey, keyUsage)

	keys := Keys{
		Keys: []jwk.Key{key},
	}

	file, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal key into JSON: %s\n", err)
		return
	}

	_ = ioutil.WriteFile("../assets/jwks.json", file, os.ModePerm)
}
