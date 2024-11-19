package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		row := db.QueryRow("select id, password from users where username=?", user.Username)
		var id int64
		var hashed string
		err = row.Scan(&id, &hashed)
		if err != nil {
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(user.Password))
		if err != nil {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(id, w)
		if err != nil {
			http.Error(w, "failed to create a token", http.StatusInternalServerError)
			return
		}
		log.Printf("token: %s", token)

		data, err := json.Marshal(
			User{ID: id, Username: user.Username})
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}

		log.Printf("User %s is login", user.Username)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func registerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			http.Error(w, "failed to hashed password", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("insert into users(username, password) values(?, ?)", user.Username, hashed)
		if err != nil {
			http.Error(w, "failed to insert a record to db", http.StatusBadRequest)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "failed to get id", http.StatusBadRequest)
			return
		}
		user.ID = id

		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DeleteSession(w)
		fmt.Fprint(w, "session is deleted")
	}
}
