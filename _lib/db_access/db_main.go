package dbaccess

import (
	"github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgtype"
)

var DB_Conn = &pgx.Conn{}

type void_func = func()

func ConnectToDB(ctx ftp_contxt.Context) void_func {
	DB_Conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return func() {
		conn.Close(ctx)
	}

}
