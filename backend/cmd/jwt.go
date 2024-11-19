package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var cookieDuration time.Duration = time.Hour * 1

type Token struct {
	Token string `json:"token"`
}

var secretKey = []byte("secretpassword")

type Payload struct {
	UserID    int64
	ExpiresAt time.Time
}

func encrypt(payload Payload) (string, error) {
	claims := jwt.MapClaims{}
	claims["userId"] = payload.UserID
	claims["expiresAt"] = payload.ExpiresAt.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func decrypt(session string) (*Payload, error) {
	token, err := jwt.Parse(session, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	id, ok := claims["userId"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid id")
	}

	expiresAt, ok := claims["expiresAt"].(float64)
	if !ok || int64(expiresAt) < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}
	return &Payload{UserID: int64(id), ExpiresAt: time.Unix(int64(expiresAt), 0)}, nil
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac
func GenerateToken(userID int64, w http.ResponseWriter) (string, error) {
	expiresAt := time.Now().Add(cookieDuration)
	session, err := encrypt(Payload{UserID: userID, ExpiresAt: expiresAt})
	if err != nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    session,
		HttpOnly: true,
		Secure:   false, // this must be true if running on a production
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	return session, nil
}

func UpdateToken(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return err
	}
	session := cookie.Value
	payload, err := decrypt(session)
	if err != nil {
		return err
	}

	payload.ExpiresAt = time.Now().Add(cookieDuration)
	session, err = encrypt(*payload)
	if err != nil {
		return err
	}

	cookie = &http.Cookie{
		Name:     "session",
		Value:    session,
		HttpOnly: true,
		Secure:   false, // this must be true if running on a production
		Expires:  payload.ExpiresAt,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac
func VerifyToken(session string) (*Payload, error) {
	payload, err := decrypt(session)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func DeleteSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // this must be true if running on a production
		MaxAge:   0,     // this makes a cookie will be deleted
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}
