package main

import (
	"os"
	"testing"
	"user-permissions-api/pkg/db"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	err := db.ConnectDB(logrus.New())
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}
