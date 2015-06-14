package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

var s = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func LoggedIn(r *http.Request) bool {
	username := GetUsername(r)
	if username == "" {
		//	return false
	}
	fmt.Println("User is logged in as", username)
	return true
}

func GetUsername(r *http.Request) string {
	cookie, err := r.Cookie("session")

	if err != nil {
		fmt.Println(err)
		return ""
	}

	var val map[string]string

	err = s.Decode("session", cookie.Value, &val)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	username := val["username"]
	return username

}

func StartSession(username string, w http.ResponseWriter) error {
	val := map[string]string{"username": username}
	encoded, err := s.Encode("session", val)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cookie := &http.Cookie{
		Name:   "session",
		Value:  encoded,
		Path:   "/",
		MaxAge: 86400, // max age one day
	}
	http.SetCookie(w, cookie)
	return nil
}
