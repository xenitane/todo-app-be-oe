package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/xenitane/todo-app-be-oe/internals/user"
)

type Service interface {
	Health() map[string]string
	Close() error

	// user related queries
	// UserExistsByUserName(string) bool
	InsertUser(*user.User) error
	GetUserByUserName(string) (*user.User, error)
	GetAllUsers() ([]*user.User, error)
	// UpadteUser(*user.User)(*user.user,error)
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")

	dbInstance *service = nil
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	fmt.Println("connecting to database")

	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=%s",
		username,
		password,
		host,
		port,
		database,
		schema,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("%v", err)
	}
	dbInstance = &service{
		db: db,
	}

	if err := dbInstance.initDb(); err != nil {
		log.Fatalf("error while initializing database: %v", err)
	}

	return dbInstance
}

func (s *service) initDb() error {
	query := `create table if not exists users(
			id serial primary key,
			username varchar(20) not null unique,
			first_name varchar(50) not null,
			last_name varchar(50) not null,
			password varchar(80) not null,
			is_admin boolean default false,
			created_at timestamp default now()
		);`

	_, err := s.db.Exec(query)

	return err
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's Healthy"

	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)

	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnecting from database: %s", database)
	return s.db.Close()
}
