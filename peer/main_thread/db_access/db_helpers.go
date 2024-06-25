package db_helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db "github.com/it-shiloheye/ftp_system_v2/lib/db_access"
	db_access "github.com/it-shiloheye/ftp_system_v2/lib/db_access/generated"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_v2/lib/tls_handler/v2"
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	"github.com/jackc/pgx/v5/pgtype"
)

var DB = db.DB
var ServerConfig = server_config.ServerConfig
var Storage = server_config.StorageSingleton
var Logger = logging.Logger

func ticker(loc log_item.Loc, i int, v ...any) {
	log.Printf("%s\n%03d: %s\n", string(loc), i, fmt.Sprint(v...))
}

func ConnectClient(ctx ftp_context.Context) error {
	loc := log_item.Loc(`ConnectClient(ctx ftp_context.Context) error`)
	var err1, err2, err3, err4, err5, err6, err7 error
	var db_peers []*db_access.PeersTable

	if DB == nil {
		log.Fatalln("DB pointer is nil")
	}
	if ServerConfig == nil {
		log.Fatalln("ServerConfig is nil")
	}

	ip_addr := ServerConfig.LocalIp()
	if ip_addr == nil {
		log.Fatalln("ip address missing")
	}

	ticker(loc, 1, "before query")
	db_peers, err1 = DB.ConnectClient(context.TODO(), db.DB_Conn, ip_addr.String())
	if err1 != nil {
		return Logger.LogErr(loc, err1)
	}

	ticker(loc, 2, "before template_cd")
	template_cd := ftp_tlshandler.CertData{
		Organization:  "Shiloh Eye, Ltd",
		Country:       "KE",
		Province:      "Coast",
		Locality:      "Mombasa",
		StreetAddress: "2nd Floor, SBM Bank Centre, Nyerere Avenue, Mombasa",
		PostalCode:    "80100",
		NotAfter: ftp_tlshandler.NotAfterStruct{
			Days: 31,
		},
		IPAddrresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.IPv6loopback,
			ip_addr.To4(),
		},
	}

	if len(db_peers) > 0 {
		ticker(loc, 3, "need to decode PEM")
		db_peer := db_peers[0]
		Storage.PeerId = db_peer.PeerID
		log.Println("db_peers exists")
		tmp_ca_pem := ftp_tlshandler.CAPem{}

		err2 = json.Unmarshal(db_peer.Pem, &tmp_ca_pem)
		if err2 != nil {
			return Logger.LogErr(loc, err2)
		}

		// generate new tls each time
		x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)

		ServerConfig.TLS_Cert, err3 = ftp_tlshandler.GenerateTLSCert(&tmp_ca_pem, x509_tls_cert)
		if err3 != nil {
			return Logger.LogErr(loc, err3)
		}

		return nil
	}
	log.Println("registering new client to db")

	// simple time guard, update cert every 7 days, server restarts every day at least once
	ServerConfig.TLS_Cert_Creation = time.Now()

	tmp_x509 := ftp_tlshandler.ExampleCACert(template_cd)

	tmp_ca_pem, err4 := ftp_tlshandler.GenerateCAPem(tmp_x509)
	if err4 != nil {
		return Logger.LogErr(loc, err4)
	}

	// generate new tls each time
	x509_tls_cert := ftp_tlshandler.ExampleTLSCert(template_cd)

	ticker(loc, 4, "before tls")
	ServerConfig.TLS_Cert, err5 = ftp_tlshandler.GenerateTLSCert(&tmp_ca_pem, x509_tls_cert)
	if err5 != nil {
		return Logger.LogErr(loc, err5)
	}
	ticker(loc, 5, "before encoding to ca_pem")

	var d []byte
	d, err6 = json.Marshal(&tmp_ca_pem)
	if err6 != nil {
		return Logger.LogErr(loc, err6)
	}

	var tmp_peerids []pgtype.UUID

	tmp_peerids, err7 = DB.CreateClient(ctx, db.DB_Conn, &db_access.CreateClientParams{
		IpAddress: ip_addr.String(),
		Pem:       d,
	})
	if err7 != nil {
		return Logger.LogErr(loc, err7)
	}

	Storage.PeerId = tmp_peerids[0]

	return nil
}
