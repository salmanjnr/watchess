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
	"watchess.org/watchess/pkg/models"
	"watchess.org/watchess/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	tournaments   interface {
		Insert(string, string, string, bool, time.Time, time.Time, bool) (int, error)
		Get(int) (*models.Tournament, error)
		LatestActive(int) ([]*models.Tournament, error)
		LatestFinished(int) ([]*models.Tournament, error)
		Upcoming(int) ([]*models.Tournament, error)
	}
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/watchess?parseTime=true", "MariaDB data source name")
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

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: tc,
		tournaments:   &mysql.TournamentModel{DB: db},
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
