package browserserver

import (
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	"github.com/it-shiloheye/ftp_system_v2/peer/server"
)

var ServerConfig = server_config.ServerConfig
var C_loc = server.C_loc

func CreateBrowserServer() *server.ServerType {

	Srvr := &server.ServerType{
		Port: ServerConfig.BrowserPort,
	}

	Srvr.InitServer(nil)

	return Srvr

}
