package api

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"test_system/config"
	"time"
)

type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

func basicAuthHandler(h HandlerWithError) HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return errors.New("unauthorized")
		}
		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
			http.Error(w, "Bad request", http.StatusUnauthorized)
			return errors.New("bad request")
		}

		payload, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
		if err != nil {
			http.Error(w, "bad request", http.StatusUnauthorized)
			return errors.New("bad request")
		}
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		if user, ok := config.TestConfig.Credentials[pair[0]]; ok {
			if user.Password != pair[1] {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return errors.New("unauthorized")
			}
			return h(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return errors.New("unauthorized")
		}
	}
}

// Think about codes
func logMiddleware(h HandlerWithError) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		err := basicAuthHandler(h)(w, r)
		duration := time.Since(start)
		if err != nil {
			log.Printf("Method: %s, URI: %s, Duration %v, Error: %s", method, uri, duration, err)
		} else {
			log.Printf("Method: %s, URI: %s, Duration: %s", method, uri, duration)
		}
	}
	return http.HandlerFunc(logFunc)
}
