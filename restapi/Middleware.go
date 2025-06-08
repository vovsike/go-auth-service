package main

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

//func Auth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token := r.Header.Get("Authorization")
//		jwtValue := strings.Split(token, "Bearer")[1]
//		priv, _ := loadRSAPrivateKeyFromFile("private.pem")
//		_, err := jwt.Parse([]byte(jwtValue), jwt.WithKey(jwa.RS256(), priv.Public()))
//		if err != nil {
//			log.Println("JWT validation failed, error: ", err, "")
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}
