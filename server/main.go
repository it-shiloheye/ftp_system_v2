package main

import (
	"log"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/_lib/context"
	"github.com/it-shiloheye/ftp_system_v2/_lib/logging"
	server_config "github.com/it-shiloheye/ftp_system_v2/server/config"
	"github.com/it-shiloheye/ftp_system_v2/server/server"
)

var ServerConfig = server_config.ServerConfig

func init() {
	logging.InitialiseLogging(ServerConfig.LogDirectory)

}

func main() {

	ctx := ftp_context.CreateNewContext()
	log.Println("hello world from server")
	go logging.Logger.Engine(ctx, ServerConfig.LogDirectory)
	server.ServerLoop(ctx)
}
