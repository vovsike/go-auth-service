package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"log"
	"net/http"
	"os"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtValue := r.Header.Get("Authorization")
		priv, _ := loadRSAPrivateKeyFromFile("private.pem")
		_, err := jwt.Parse([]byte(jwtValue), jwt.WithKey(jwa.RS256(), priv.Public()))
		if err != nil {
			log.Println("JWT validation failed, error: ", err, "")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loadRSAPrivateKeyFromFile(fileName string) (*rsa.PrivateKey, error) {
	privKeyBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	privPem, _ := pem.Decode(privKeyBytes)
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
