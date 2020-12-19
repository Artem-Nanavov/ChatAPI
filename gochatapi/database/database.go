package database

import (
	"database/sql"
	"io/ioutil"
	"time"
	"api/database/repo"

	_ "github.com/lib/pq" // ...
	"github.com/sirupsen/logrus"
)

// Logger ...
type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
}

// Config ...
type Config struct {
	DSN    string
	Logger Logger
}

// Database ...
type Database struct {
	db     *sql.DB
	dsn    string
	config *Config

	userRepo *repo.UserRepo
}

// NewDatabase ...
func NewDatabase(config *Config) *Database {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}
	return &Database{
		config: config,
	}
}

// Open ...
func (d *Database) Open() error {
	db, err := sql.Open("postgres", d.config.DSN)
	if err != nil {
		d.config.Logger.Error(err.Error())
		return err
	}

	if err := db.Ping(); err != nil {
		d.config.Logger.Error(err.Error())
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	d.db = db
	d.config.Logger.Info("Database connection opened successfully.")
	return nil
}

// Close ...
func (d *Database) Close() error {
	d.config.Logger.Info("Closing database connection ...")
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// User ...
func (d *Database) User() *repo.UserRepo {
	if d.userRepo == nil {
		d.userRepo = repo.NewUserRepo(d.db)
	}
	return d.userRepo
}

// Migrate ...
func (d *Database) Migrate(migrationsPath string) error {
	d.config.Logger.Info("Running migrations ...")
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		d.config.Logger.Error(err.Error())
		return err
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(migrationsPath + "/" + file.Name())

		if err != nil {
			d.config.Logger.Error(err.Error())
			return err
		}

		if _, err := d.db.Exec(string(data)); err != nil {
			d.config.Logger.Error(err.Error())
			return err
		}

		d.config.Logger.Info("\t" + file.Name() + ": done.")
	}
	d.config.Logger.Info("Migrated successfully.")
	return nil
}
