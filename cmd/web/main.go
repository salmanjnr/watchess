package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"watchess.org/watchess/pkg/models"
	"watchess.org/watchess/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	config        config
	templateCache map[string]*template.Template
	tournaments   interface {
		Insert(string, string, string, bool, time.Time, time.Time, bool) (int, error)
		Get(int) (*models.Tournament, error)
		LatestActive(int) ([]*models.Tournament, error)
		LatestFinished(int) ([]*models.Tournament, error)
		Upcoming(int) ([]*models.Tournament, error)
	}
	users interface {
		Insert(string, string, string, models.UserRole) (int, error)
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	session *sessions.Session
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/watchess?parseTime=true", "MariaDB data source name")
	secret := flag.String("secret", "a3BBA+69e27UvVv&K9P12nasdk@89ue!", "Session encryption key")
	flag.Parse()

	if (*addr)[0] != ':' {
		*addr = fmt.Sprintf(":%v", *addr)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	tc, err := newTemplateCache("./ui/html")

	if err != nil {
		errorLog.Fatalf("Couldn't initialize template cache due to error %v\n", err)
	}

	session := sessions.New([]byte(*secret))

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		config:        getConfig(),
		templateCache: tc,
		tournaments:   &mysql.TournamentModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		session:       session,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Updating ReadTimeout uses the same value for IdleTimeout unless set explicitly
		IdleTimeout: time.Minute,
		// Close connection if it takes more than five sec to read header/body
		// to mitigate the risk from slow-client attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
