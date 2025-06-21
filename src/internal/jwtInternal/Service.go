package jwtInternal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"os"
	"strconv"
	"time"
)

type ServiceImpl struct {
}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) GenerateToken(userId int) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		Expiration(time.Now().Add(5 * time.Minute)).
		Subject(strconv.Itoa(userId)).
		Issuer("https://localhost:8080").
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}
	privKey, err := loadRSAPrivKey(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256(), privKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}
	return signed, nil
}

func (s *ServiceImpl) ValidateToken(t []byte) (bool, error) {
	pubKey, err := loadRSAPubKey(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return false, fmt.Errorf("failed to load public key: %v", err)
	}
	_, err = jwt.Parse(t, jwt.WithKey(jwa.RS256(), pubKey))

	if err != nil {
		return false, fmt.Errorf("failed to validate token: %v", err)
	}
	return true, nil
}

func (s *ServiceImpl) ExtractUserID(token []byte) (int, error) {
	pubKey, err := loadRSAPubKey(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		return 0, fmt.Errorf("loading public key: %w", err)
	}

	parsed, err := jwt.Parse(token, jwt.WithKey(jwa.RS256(), pubKey))
	if err != nil {
		return 0, fmt.Errorf("parsing token: %w", err)
	}

	var userID int

	err = parsed.Get("sub", &userID)
	if err != nil {
		return 0, errors.New("subject claim not found in token")
	}

	return userID, nil
}

func loadRSAPrivKey(privKey string) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode([]byte(privKey))
	if privPem == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}

	return rsaKey, nil
}

func loadRSAPubKey(pubKey string) (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode([]byte(pubKey))
	if pubPem == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	key, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}

	return rsaKey, nil
}
