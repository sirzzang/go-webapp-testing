package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {

	var creds Credentials

	// read a json payload
	err := app.readJSON(w, r, &creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// look up the user by email address
	user, err := app.DB.GetUserByEmail(creds.Username)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unahtorized"), http.StatusUnauthorized)
		return
	}

	// generate tokens
	tokenPairs, err := app.generateTokenPairs(user)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	// send token to user
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)

}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := r.Form.Get("refresh_token")
	claims := &Claims{}

	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.JWTSecret), nil
	})

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second {
		log.Println("claims.ExpiresAt:", time.Unix(claims.ExpiresAt.Unix(), 0))
		log.Println(time.Now())
		app.errorJSON(w, errors.New("refresh token does not need renewed yet"), http.StatusTooEarly)
		return
	}

	// get the user id from the claims
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUser(userId)
	if err != nil {
		app.errorJSON(w, errors.New("unknown user"), http.StatusBadRequest)
		return
	}

	// generate new access token
	tokenPairs, err := app.generateTokenPairs(user)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// we're going to use this refresh handler to refresh tokens for API users
	// so perhaps one API is calling our API and there's no actual session or anything like that
	// we're just having an API talk to an API
	// but we'll also use this same handler to refresh web users
	// in other words, say a single page application, something written in vue or react,
	// or even a plain javascript application
	// when you're interacting our API using a web browser, we need somewhere to store that refresh token
	// and the best place to put that, is a cookie
	// store refresh token in a cookie
	http.SetCookie(w, &http.Cookie{
		// tell older browsers this cookie is a secure cookie
		// in other words, it'll apply certain kind of security protocols to that cookie
		Name: "__Host-refresh_token",
		// entire web application, entire web API
		Path:     "/",
		Value:    tokenPairs.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	})

	// send back the refersh token as part of the JSON body
	// and also set a cookie, because depending on who calls this,
	// they're either going to use the cookie for a single page web application,
	// or the refresh token is available to them in the JSON payload,
	// if it's an API to API call, for example
	_ = app.writeJSON(w, http.StatusOK, tokenPairs)

}

func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {

}
