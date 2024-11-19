package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var WaitValue time.Duration = 3

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "リバースプロキシとAPIサーバの起動",
	RunE: func(cmd *cobra.Command, args []string) error {
		connect, err := cmd.Flags().GetString("connect")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		waitValue, err := cmd.Flags().GetInt("wait")
		if err != nil {
			return err
		}
		WaitValue = time.Duration(waitValue)

		run(connect, port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 8000, "待ち受けるポート番号")
	serveCmd.Flags().String("connect", "http://localhost:3000", "接続先")
	serveCmd.Flags().Int("wait", int(WaitValue), "待機時間(秒)")
}

type Folder struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

type Task struct {
	ID          int64     `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	DueDate     time.Time `json:"due_date"`
	FolderId    int64     `json:"folder_id"`
}

type User struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

/*
  /folders
	  * GET - get all folders
	/folders/create
	  * POST - create a new folder
	/folders/update
	  * POST - update a folder
	/folders/delete
	  * POST - delete a folder
	/folders/:id
	  * GET - get a folder

  /tasks
	  * GET - get all tasks
	/tasks/create
	  * POST - create a new task
	/tasks/update
	  * POST - update a task
	/tasks/delete
	  * POST - delete a task
	/tasks/:id
	  * GET - get a task
*/

func run(connect string, port int) {
	remote, err := url.Parse(connect)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ***** folder *****

	mux := http.NewServeMux()

	mux.HandleFunc("/api/folders", foldersHandler(db))
	mux.HandleFunc("/api/folders/{id}", folderHandler(db))
	mux.HandleFunc("/api/folders/create", createFolderHandler(db))
	mux.HandleFunc("/api/folders/update", updateFolderHandler(db))
	mux.HandleFunc("/api/folders/delete", deleteFolderHandler(db))

	mux.HandleFunc("/api/tasks", tasksHandler(db))
	mux.HandleFunc("/api/tasks/{id}", taskHandler(db))
	mux.HandleFunc("/api/tasks/create", createTaskHandler(db))
	mux.HandleFunc("/api/tasks/update", updateTaskHandler(db))
	mux.HandleFunc("/api/tasks/delete", deleteTaskHandler(db))

	mux.HandleFunc("/api/login", loginHandler(db))
	mux.HandleFunc("/api/register", registerHandler(db))
	mux.HandleFunc("/api/logout", logoutHandler())

	// ***** Reverse Proxy *****
	reverseProxyHandler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			// w.Header().Set("xxx", "yyy")
			p.ServeHTTP(w, r)
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	mux.HandleFunc("/", reverseProxyHandler(proxy))

	middleware := AuthenticationMiddleware(mux)

	fmt.Println("listening...")
	err = http.ListenAndServe(":"+strconv.Itoa(port), middleware)
	if err != nil {
		panic(err)
	}
}
