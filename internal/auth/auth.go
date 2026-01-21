package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type CookieData struct {
	Name  string
	Value CookieValue
}

type CookieValue struct {
	UserID    string
	LoginDate int64
}

func CreateSession(w http.ResponseWriter, userID string) error {
	cookieSession := CookieData{
		Name: "_session",
		Value: CookieValue{
			UserID:    userID,
			LoginDate: time.Now().Unix(),
		},
	}
	cookie, err := genCookie(cookieSession)

	if err != nil {
		return err
	}

	http.SetCookie(w, &cookie)

	return nil
}

func GetSession(r *http.Request) (CookieValue, error) {
	cookie, err := r.Cookie("_session")

	if err != nil {
		return CookieValue{}, errors.New("cookie name not valid")
	}

	value, err := decodeAndValidateCookie[CookieValue](cookie)

	if err != nil {
		return CookieValue{}, err
	}

	return value, nil
}

func decodeAndValidateCookie[T any](cookie *http.Cookie) (T, error) {
	var result T

	encodedValue, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return result, err
	}

	if len(encodedValue) < sha256.Size {
		return result, errors.New("cookie value too short")
	}

	signature := encodedValue[:sha256.Size]
	value := encodedValue[sha256.Size:]

	expectedSignature := calcSignature(cookie.Name, value)

	if !hmac.Equal(signature, expectedSignature) {
		return result, errors.New("invalid signature")
	}

	err = json.Unmarshal(value, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func genCookie(cookieData CookieData) (http.Cookie, error) {
	cookieValueMarshal, err := json.Marshal(cookieData.Value)

	if err != nil {
		return http.Cookie{}, err
	}

	signature := calcSignature(cookieData.Name, cookieValueMarshal)

	valueEncoded := base64.URLEncoding.EncodeToString(append(signature, cookieValueMarshal...))

	cookie := http.Cookie{
		Name:     cookieData.Name,
		Value:    valueEncoded,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 400,
	}

	return cookie, nil
}

func calcSignature(cookieName string, cookieValueMarshal []byte) []byte {
	hmac := hmac.New(sha256.New, []byte("asdf123456"))
	hmac.Write([]byte(cookieName))
	hmac.Write(cookieValueMarshal)
	signature := hmac.Sum(nil)

	return signature
}
