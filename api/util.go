// Package api contains all api handlers for the application.
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

// HandlerWithError is a type for http handlers that can return an error
type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

// errorHandler is a helper function that writes an error message to the response and returns an error
// w is the http.ResponseWriter that we want to write the error message to
// message is the message that we want to write to the response
// code is the status code that we want to return
func errorHandler(w http.ResponseWriter, message string, code int) error {
	http.Error(w, message, code)
	return errors.New(message)
}

// basicAuthHandler is a middleware that checks for basic auth
// h is the handler that we want to call if the user is authorized
// returns a HandlerWithError that contains the logic for checking the basic auth header
func basicAuthHandler(h HandlerWithError) HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return errorHandler(w, "unauthorized", http.StatusUnauthorized)
		}
		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Basic" {
			return errorHandler(w, "bad request", http.StatusBadRequest)
		}

		payload, err := base64.StdEncoding.DecodeString(authHeaderParts[1])
		if err != nil {
			return errorHandler(w, "bad request", http.StatusBadRequest)
		}
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			return errorHandler(w, "unauthorized", http.StatusUnauthorized)
		}
		if user, ok := config.TestConfig.Credentials[pair[0]]; ok {
			if user.Password != pair[1] {
				return errorHandler(w, "unauthorized", http.StatusUnauthorized)
			}
			return h(w, r)
		} else {
			return errorHandler(w, "unauthorized", http.StatusUnauthorized)
		}
	}
}

// logMiddleware is a middleware that logs requests and use basicAuthHandler middleware
// h is the handler that we want to be logged and checked for basic auth
// returns a http.Handler that contains the logic for logging requests
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
