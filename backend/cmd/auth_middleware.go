package cmd

import (
	"log"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		isProtected := strings.HasPrefix(path, "/api") && path != "/api/login" && path != "/api/register"
		if !isProtected {
			log.Printf("%s is public zone, skip authentication\n", path)
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("Auth Error: no session cookie")
			http.Error(w, "no session cookie", http.StatusUnauthorized)
			return
		}

		log.Printf("session value = %s", cookie.Value)
		payload, err := VerifyToken(cookie.Value)
		if err != nil {
			log.Print(err)
			log.Println("Auth Error: invalid session")
			http.Error(w, "invalid session", http.StatusUnauthorized)
			return
		}
		log.Printf("Request: %s %s %v", r.Method, r.URL.Path, payload)

		next.ServeHTTP(w, r)

		UpdateToken(w, r)
		log.Println("session is verified")
	})
}
