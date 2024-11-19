package cmd

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func tasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		folderId := r.URL.Query().Get("folderId")

		var rows *sql.Rows
		var err error

		if folderId == "" {
			rows, err = db.Query("select id, title, description, completed, due_date, folder_id from tasks")
		} else {
			rows, err = db.Query("select id, title, description, completed, due_date, folder_id from tasks where folder_id = ?", folderId)
		}
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to query", http.StatusBadRequest)
			return
		}
		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var id int64
			var title string
			var description string
			var completed bool
			var due_date string
			var folderId int64
			err := rows.Scan(&id, &title, &description, &completed, &due_date, &folderId)
			if err != nil {
				log.Println(err)
				http.Error(w, "failed to get task data", http.StatusBadRequest)
				return
			}

			dueDate, err := time.Parse(time.RFC3339, due_date)
			if err != nil {
				log.Println(err)
				http.Error(w, "failed to parse due_date", http.StatusBadRequest)
			}

			tasks = append(tasks, Task{
				ID:          id,
				Title:       title,
				Description: description,
				Completed:   completed,
				DueDate:     dueDate,
				FolderId:    folderId,
			})
		}

		data, err := json.Marshal(tasks)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

// /folders/{id}
func taskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		idString := r.PathValue("id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, "illegal id", http.StatusBadRequest)
			return
		}
		row := db.QueryRow(
			`select title, description, completed, due_date, folder_id from tasks where id = ?`, id)
		var title string
		var description string
		var completed bool
		var due_date string
		var folderId int64
		err = row.Scan(&title, &description, &completed, &due_date, &folderId)
		if err != nil {
			http.Error(w, "database error", http.StatusBadRequest)
			return
		}

		dueDate, err := time.Parse(time.RFC3339, due_date)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to parse due_date", http.StatusBadRequest)
		}

		folder := Task{
			ID:          id,
			Title:       title,
			Description: description,
			Completed:   completed,
			DueDate:     dueDate,
			FolderId:    folderId,
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

func createTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		// b, err := io.ReadAll(r.Body)
		// fmt.Println(string(b))

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Print(err)
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			`insert into tasks(title, description, completed, due_date, folder_id) 
			values(?, ?, ?, ?, ?)`,
			task.Title,
			task.Description,
			task.Completed,
			task.DueDate.Format(time.RFC3339),
			task.FolderId)
		if err != nil {
			http.Error(w, "failed to insert a folder", http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "failed to get id", http.StatusBadRequest)
			return
		}
		task.ID = id
		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func updateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			`update tasks set title=?, description=?, completed=?, due_date=?, folder_id=? where id=?`,
			task.Title,
			task.Description,
			task.Completed,
			task.DueDate.Format(time.RFC3339),
			task.FolderId,
			task.ID)
		if err != nil {
			http.Error(w, "failed to update a folder", http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}

func deleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(WaitValue * time.Second)

		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("delete from tasks where id=?", task.ID)
		if err != nil {
			http.Error(w, "failed to insert", http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "failed to marshal JSON", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(data)
	}
}
