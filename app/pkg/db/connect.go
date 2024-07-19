package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

var (
	DataBase *pgx.Conn       = nil
	Ctx      context.Context = context.Background()
)

func ConnectDB(log *logrus.Logger) error {
	var err error
	host := os.Getenv("PGHOST")
	port, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		log.Errorf("Invalid port number: %v", err)
		return err
	}
	user := os.Getenv("PGUSER")
	dbname := os.Getenv("PGNAME")
	password := os.Getenv("PGPASSWORD")
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	DataBase, err = pgx.Connect(Ctx, url)
	if err != nil {
		log.Fatal(fmt.Sprintf("no connection database : %v\n", err))
		return err
	}
	log.Info("Successfully connected to the database.")
	return nil
}
