package networkpeer

import (
	"log"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	db_helpers "github.com/it-shiloheye/ftp_system_v2/peer/main_thread/db_access"
	"github.com/it-shiloheye/ftp_system_v2/peer/server"
)

var ServerConfig = server_config.ServerConfig
var Logger = logging.Logger

func CreatePeerServer(ctx ftp_context.Context) *server.ServerType {
	loc := log_item.Locf(`CreatePeerServer(ctx ftp_context.Context) *server.ServerType`)
	err1 := db_helpers.ConnectClient(ctx)
	if err1 != nil {
		Logger.LogErr(loc, err1)
		<-time.After(time.Second)
		log.Fatalln("fatal termination")
	}
	Srvr := &server.ServerType{
		Port: ServerConfig.PeerPort,
	}

	Srvr.InitServer(ServerConfig.TLS_Cert, "peer")

	return Srvr

}
