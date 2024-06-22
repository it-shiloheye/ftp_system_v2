package db

import (
	"fmt"
	"log"
	"os"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db_access "github.com/it-shiloheye/ftp_system_v2/lib/db_access/generated"
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgtype"
)

var DB_Conn = &pgx.Conn{}
var DB = &db_access.Queries{}

var db_url string

func init() {
	db_url = os.Getenv("DATABASE_URL")
	if len(db_url) < 1 {
		log.Fatalln(`Fatal: "DATABASE_URL" is missing in .env`)
	}

}

type void_func = func()

func ConnectToDB(ctx ftp_context.Context) void_func {
	DB_Conn, err := pgx.Connect(ctx, db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	DB = NewDB(ctx)
	log.Println("connected to: ", db_url)
	return func() {
		DB_Conn.Close(ctx)
	}
}

func NewDB(ctx ftp_context.Context) *db_access.Queries {

	return db_access.New(DB_Conn)
}
