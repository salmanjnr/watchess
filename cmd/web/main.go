package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"watchess.org/watchess/pkg/models"
	"watchess.org/watchess/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	logger        *zap.Logger
	config        config
	templateCache map[string]*template.Template
	tournaments   interface {
		Insert(string, string, string, bool, time.Time, time.Time, bool, int) (int, error)
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
	rounds interface {
		Insert(string, string, models.GameReward, models.GameReward, time.Time, int) (int, error)
		Get(int) (*models.Round, error)
		GetByTournament(int) ([]*models.Round, error)
	}
	matches interface {
		Insert(string, string, int) (int, error)
		Get(int) (*models.Match, error)
		GetByRound(int) ([]*models.Match, error)
	}
	games interface {
		Insert(string, string, models.GameResult, string, string, string, int, int) (int, error)
		Get(int) (*models.Game, error)
		GetByMatch(int) ([]*models.Game, error)
		GetByRound(int) ([]*models.Game, error)
	}
	session *sessions.Session
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/watchess?parseTime=true", "MariaDB data source name")
	secret := flag.String("secret", "a3BBA+69e27UvVv&K9P12nasdk@89ue!", "Session encryption key")
	dev := flag.Bool("dev", false, "Run in dev mode")
	flag.Parse()

	if (*addr)[0] != ':' {
		*addr = fmt.Sprintf(":%v", *addr)
	}

	logger, err := getLogger(dev)
	if err != nil {
		log.Fatal("Can't initialize logger")
	}

	// Direct error produced by http.server through zap
	undo, err := zap.RedirectStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		logger.Fatal("Can't redirect std logger output through zap")
	}
	defer undo()

	db, err := openDB(*dsn)
	if err != nil {
		logger.Fatal("Can't open database", zap.Error(err))
	}

	defer db.Close()

	tc, err := newTemplateCache("./ui/html")

	if err != nil {
		logger.Fatal("Couldn't initialize template cache", zap.Error(err))
	}

	session := sessions.New([]byte(*secret))

	app := &application{
		logger:        logger,
		config:        getConfig(),
		templateCache: tc,
		tournaments:   &mysql.TournamentModel{DB: db},
		rounds:        &mysql.RoundModel{DB: db},
		matches:       &mysql.MatchModel{DB: db},
		games:         &mysql.GameModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		session:       session,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Updating ReadTimeout uses the same value for IdleTimeout unless set explicitly
		IdleTimeout: time.Minute,
		// Close connection if it takes more than five sec to read header/body
		// to mitigate the risk from slow-client attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", zap.String("port", *addr))
	err = srv.ListenAndServe()
	logger.Fatal("Failure to initialize server", zap.Error(err))
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

func getLogger(dev *bool) (*zap.Logger, error) {
	if (dev == nil) || (*dev == false) {
		logger, err := zap.NewProduction()
		return logger, err
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := config.Build()
		return logger, err
	}
}
