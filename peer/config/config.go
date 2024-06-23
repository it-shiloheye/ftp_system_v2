package server_config

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var ServerConfig = &ServerConfigStruct{}

func init() {
	ReadConfig()
}

type ServerConfigStruct struct {
	ServerId string `json:"server_id"`

	PeerPort    string `json:"peer_port"`
	BrowserPort string `json:"browser_port"`

	CertsDirectory    string `json:"certs_dir"`
	StorageDirectory  string `json:"storage_directory"`
	DatabaseURL       string `json:"db_url"`
	TLS_Cert_Creation time.Time
}

func ReadConfig() (err error) {
	ok := false
	godotenv.Load()
	ServerConfig.ServerId, ok = os.LookupEnv("SERVER_ID")
	if !ok || len(ServerConfig.ServerId) < 1 {
		log.Println(`"SERVER_ID" missing from .env`)
	}
	ServerConfig.PeerPort, ok = os.LookupEnv("PEER_PORT")
	if !ok || len(ServerConfig.PeerPort) < 1 {
		log.Fatalln(`Fatal "PEER_PORT" missing from .env`)
	} else {
		ServerConfig.PeerPort = ":" + ServerConfig.PeerPort
	}
	ServerConfig.BrowserPort, ok = os.LookupEnv("BROWSER_PORT")
	if !ok || len(ServerConfig.BrowserPort) < 1 {
		log.Fatalln(`Fatal "BROWSER_PORT" missing from .env`)
	} else {
		ServerConfig.BrowserPort = ":" + ServerConfig.BrowserPort
	}

	ServerConfig.StorageDirectory, ok = os.LookupEnv("STORAGE_DIRECTORY")
	if !ok || len(ServerConfig.StorageDirectory) < 1 {
		log.Fatalln(`Fatal "STORAGE_DIRECTORY" missing from .env`)
	}

	ServerConfig.CertsDirectory = ServerConfig.StorageDirectory + "/ssl_certs"

	ServerConfig.DatabaseURL, ok = os.LookupEnv("DATABASE_URL")
	if !ok || len(ServerConfig.DatabaseURL) < 1 {
		log.Println(`"DATABASE_URL" missing from .env`)
	}

	log.Println("successfully loaded .env")
	return
}

// Get preferred local outbound ip of this machine
func (scs *ServerConfigStruct) LocalIp() net.IP {

	conn, err1 := net.Dial("udp", "192.168.0.1:80")
	if err1 != nil {
		conn1, err2 := net.Dial("udp", "192.168.1.1:80")
		if err2 != nil {

			log.Fatalln(err1, "\n", err2)
		}

		conn = conn1
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Get preferred web outbound ip of this machine
func WebIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
