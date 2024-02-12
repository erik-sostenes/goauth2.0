package business

import (
	"crypto/rsa"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type (
	// TokenGenerator contains the method that generates the token using the JWT standard
	TokenGenerator interface {
		// GenerateToken method that generates the token by id, email
		GenerateToken(id, email string) (string, error)
	}

	// RegisteredClaims is a structure representing JWT own claim types
	RegisteredClaims = gojwt.RegisteredClaims

	// Authorization is a structure representing the payload of a token
	Authorization struct {
		Id, Email string
		RegisteredClaims
	}

	tokenGenerator struct {
		privateKey *rsa.PrivateKey
	}
)

func NewTokenGenerator(privateRSA string) *tokenGenerator {
	if privateRSA == "" {
		panic(`missing privateRSA`)
	}

	privateKey, err := gojwt.ParseRSAPrivateKeyFromPEM([]byte(privateRSA))
	if err != nil {
		panic(err)
	}

	return &tokenGenerator{
		privateKey: privateKey,
	}
}

func (jwt tokenGenerator) GenerateToken(id, email string) (string, error) {
	authorization := &Authorization{
		Id:    id,
		Email: email,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    "goauth",
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, authorization)

	return token.SignedString(jwt.privateKey)
}
