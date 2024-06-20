package server_config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var ServerConfig = &ServerConfigStruct{}

func init() {
	ReadConfig()
}

type ServerConfigStruct struct {
	ServerId          string `json:"server_id"`
	LocalIp           string `json:"local_ip"`
	WebIp             string `json:"web_ip"`
	ServerPort        string `json:"server_port"`
	CertsDirectory    string `json:"certs_dir"`
	StorageDirectory  string `json:"storage_directory"`
	LogDirectory      string `json:"log_directory"`
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
	ServerConfig.ServerPort, ok = os.LookupEnv("SERVER_PORT")
	if !ok || len(ServerConfig.ServerPort) < 1 {
		log.Fatalln(`Fatal "SERVER_PORT" missing from .env`)
	} else {
		ServerConfig.ServerPort = ":" + ServerConfig.ServerPort
	}
	ServerConfig.LocalIp, ok = os.LookupEnv("LOCAL_SERVER_IP")
	if !ok || len(ServerConfig.LocalIp) < 1 {
		log.Fatalln(`Fatal "LOCAL_SERVER_IP" missing from .env`)
	}
	ServerConfig.WebIp, ok = os.LookupEnv("REMOTE_SERVER_IP")
	if !ok || len(ServerConfig.WebIp) < 1 {
		log.Println(`"REMOTE_SERVER_IP" missing from .env`)
	}
	ServerConfig.CertsDirectory, ok = os.LookupEnv("CERTS_DIRECTORY")
	if !ok || len(ServerConfig.CertsDirectory) < 1 {
		log.Fatalln(`Fatal "CERTS_DIRECTORY" missing from .env`)
	}
	ServerConfig.LogDirectory, ok = os.LookupEnv("LOG_DIRECTORY")
	if !ok || len(ServerConfig.LogDirectory) < 1 {
		log.Fatalln(`Fatal "LOG_DIRECTORY" missing from .env`)
	}
	ServerConfig.DatabaseURL, ok = os.LookupEnv("DATABASE_URL")
	if !ok || len(ServerConfig.DatabaseURL) < 1 {
		log.Println(`"DATABASE_URL" missing from .env`)
	}
	return
}
