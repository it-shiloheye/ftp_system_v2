package main

import (
	"log"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db "github.com/it-shiloheye/ftp_system_v2/lib/db_access"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	server_config "github.com/it-shiloheye/ftp_system_v2/server-peer/config"
	"github.com/it-shiloheye/ftp_system_v2/server-peer/server"
)

var ServerConfig = server_config.ServerConfig

func init() {
	logging.InitialiseLogging(ServerConfig.LogDirectory)

}

func main() {

	ctx := ftp_context.CreateNewContext()
	log.Println("hello world from server")

	close_db_conn := db.ConnectToDB(ctx)

	defer close_db_conn()
	go logging.Logger.Engine(ctx, ServerConfig.LogDirectory)
	server.ServerLoop(ctx)
}
