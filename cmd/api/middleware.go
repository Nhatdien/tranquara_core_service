package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each
	// client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu sync.Mutex
		// Update the map so the values are pointers to a client struct.
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {

			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceedResponse(w, r)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (app *application) testPostMiddleWare(previous http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		previous.ServeHTTP(w, r)

		app.logger.PrintInfo("post-middleware called", nil)

	})
}

const userCtxKey = "user"

func (app *application) authMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		pubKey, err := app.loadPublicKey("/publicKey.pem")
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}

		prefix := "Bearer "
		cutBearerToken := strings.TrimPrefix(authHeader, prefix)
		token, err := jwt.Parse(cutBearerToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("error with signing method")
			}

			return pubKey, nil
		})

		if err != nil || !token.Valid {
			app.errorResponse(w, r, http.StatusUnauthorized, "Invalid token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), userCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		app.errorResponse(w, r, http.StatusUnauthorized, "Error when get claims")
	})
}

func (app *application) GetUserFromContext(ctx context.Context) jwt.MapClaims {
	claims, ok := ctx.Value(userCtxKey).(jwt.MapClaims)
	if !ok {
		return nil
	}
	return claims
}

func (app *application) GetUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	claims := ctx.Value(userCtxKey).(jwt.MapClaims)

	return uuid.Parse(claims["sub"].(string))
}
