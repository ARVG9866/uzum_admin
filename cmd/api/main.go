package main

import (
	"log"

	"github.com/Shemistan/uzum_admin/cmd/api/serv"
	"github.com/Shemistan/uzum_admin/cmd/conf"
	"github.com/Shemistan/uzum_admin/internal/models"
	desc "github.com/Shemistan/uzum_admin/pkg/auth_v1"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort = ":8080"
)

func main() {
	cfg, err := conf.NewConfig()
	if err != nil {
		log.Fatal("failed to get congigs", err.Error())
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("failed to init DB", err.Error())
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Println("failed to close connection to DV:", err.Error())
		}
	}()

	conn, err := grpc.Dial(cfg.App.AuthClient, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewAuthV1Client(conn)

	srv, err := serv.GetServ(cfg, db, c)
	if err != nil {
		log.Fatal("failed to get serv", err.Error())
	}

	log.Println("Admin server is running at port:", httpPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func initDB(cnf models.Config) (*sqlx.DB, error) {
	sqlConnectionString := conf.GetSqlConnectionString(cnf)

	db, err := sqlx.Open("postgres", sqlConnectionString)
	if err != nil {
		return nil, err
	}

	// Проверка доступности БД
	if err = db.Ping(); err != nil {
		log.Println("failed to ping DB")
		return nil, err
	}

	log.Println("Connection to DB success")

	return db, nil
}
