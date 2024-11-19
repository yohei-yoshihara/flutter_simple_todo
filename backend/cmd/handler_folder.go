package cmd

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func foldersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		rows, err := db.Query("select id, name from folders")
		if err != nil {
			http.Error(w, "failed to query database", http.StatusBadRequest)
			return
		}
		defer rows.Close()

		folders := []Folder{}
		for rows.Next() {
			var id int64
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				http.Error(w, "failed to get a folder from database", http.StatusBadRequest)
				return
			}
			folders = append(folders, Folder{ID: id, Name: name})
		}

		data, err := json.Marshal(folders)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func folderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)
		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "illegal id", http.StatusBadRequest)
			return
		}
		row := db.QueryRow("select name from folders where id = ?", id)
		var name string
		err = row.Scan(&name)
		if err != nil {
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}
		folder := Folder{ID: id, Name: name}
		data, err := json.Marshal(folder)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func createFolderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)
		var folder Folder

		err := json.NewDecoder(r.Body).Decode(&folder)
		if err != nil {
			http.Error(w, "failed to decode JSON", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("insert into folders(name) values(?)", folder.Name)
		if err != nil {
			http.Error(w, "failed to insert", http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "failed to get id", http.StatusBadRequest)
			return
		}
		folder.ID = id
		data, err := json.Marshal(folder)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func updateFolderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)
		var folder Folder

		err := json.NewDecoder(r.Body).Decode(&folder)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("update folders set name=? where id=?", folder.Name, folder.ID)
		if err != nil {
			http.Error(w, "failed to update", http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(folder)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func deleteFolderHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)
		var folder Folder

		err := json.NewDecoder(r.Body).Decode(&folder)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("delete from folders where id=?", folder.ID)
		if err != nil {
			http.Error(w, "failed to insert", http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(folder)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}
