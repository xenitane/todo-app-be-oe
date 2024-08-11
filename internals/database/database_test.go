package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartPostgresContainer() (func(context.Context) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)
	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	database = dbName
	username = dbUser
	password = dbPwd

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}
	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}
	host = dbHost
	port = dbPort.Port()
	return dbContainer.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}
	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatalf("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected pg status to be `up`, got %s", stats["status"])
	}
	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}
	if stats["message"] != "It's Healthy" {
		t.Fatalf("expected pg message to be `It's Healthy`, got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New()
	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
