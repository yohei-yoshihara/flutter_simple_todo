package cmd

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "データベースへのデータ登録",
	Run: func(cmd *cobra.Command, args []string) {
		RunSeed()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func RunSeed() {
	os.Remove("./database.db")

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	create table users (
		id integer primary key autoincrement,
		username text not null unique,
		password text not null
	);
	`)
	if err != nil {
		panic(err)
	}

	password := []byte("password")
	hashed, _ := bcrypt.GenerateFromPassword(password, 10)
	_, err = db.Exec("insert into users(username, password) values(?, ?)", "user1", hashed)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("insert into users(username, password) values(?, ?)", "user2", hashed)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	create table folders (
	  id integer primary key autoincrement, 
		name text not null);
	`)
	if err != nil {
		panic(err)
	}

	folderNames := []string{"プライベート", "仕事", "その他"}
	stmt, err := db.Prepare(`
	insert into folders(name) values(?)
	`)
	if err != nil {
		panic(err)
	}
	folderIds := []int64{}
	for _, name := range folderNames {
		result, err := stmt.Exec(name)
		if err != nil {
			panic(err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		folderIds = append(folderIds, id)
	}

	_, err = db.Exec(`
	create table tasks (
	  id integer primary key autoincrement,
		title text not null,
		description text not null,
		completed integer not null,
		due_date text not null,
		folder_id integer not null
	)
	`)
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare("insert into tasks(title, description, completed, due_date, folder_id) values(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}

	tasks := []Task{}
	for i := 0; i < 20; i++ {
		tasks = append(tasks, Task{
			Title:       fmt.Sprintf("Task-%d", i),
			Description: fmt.Sprintf("Task-%d's description", i),
			Completed:   rand.Intn(3) == 0,
			DueDate:     time.Now().AddDate(0, 0, rand.Intn(180)),
			FolderId:    folderIds[rand.Intn(10)%len(folderIds)],
		})
	}
	for _, task := range tasks {
		_, err = stmt.Exec(task.Title,
			task.Description,
			task.Completed,
			task.DueDate.Format(time.RFC3339),
			task.FolderId)
		if err != nil {
			panic(err)
		}
	}

}
