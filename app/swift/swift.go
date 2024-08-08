package swift

import (
	"bisgo/app/DB"
	"bisgo/app/P2P/client"
	"bisgo/app/manager"
	"bisgo/config"
	"bisgo/errlog"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type SwiftHandler struct {
	db                          *DB.DBHandler
	p2pClient                   *client.P2PClient
	complianceCheckStateManager *manager.ComplianceCheckStateManager
}

var handler SwiftHandler

func init() {
	handler = SwiftHandler{
		db:                          DB.GetDBHandler(),
		p2pClient:                   client.GetP2PClient(),
		complianceCheckStateManager: manager.GetComplianceCheckStateManager(),
	}
}

func GetSwiftHandler() *SwiftHandler {
	return &handler
}

func (h *SwiftHandler) ListenAndServer() {
	go func() {
		caPath := os.Getenv("CA_CERT_PATH")

		caCert, err := os.ReadFile(caPath)
		if err != nil {
			panic(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		serverCert, err := tls.LoadX509KeyPair(os.Getenv("SRV_CERT_PATH"), os.Getenv("SRV_KEY_PATH"))
		if err != nil {
			panic(err)
		}

		tlsConfig := &tls.Config{
			ClientCAs:    caCertPool,
			ClientAuth:   tls.NoClientCert,
			Certificates: []tls.Certificate{serverCert},
		}

		router := httprouter.New()

		router.HandlerFunc(http.MethodGet, "/test", h.Test)
		router.HandlerFunc(http.MethodPost, "/swift-notify", h.SwiftNotification)

		server := &http.Server{
			Addr:      config.ResolveServerTLSPort(),
			TLSConfig: tlsConfig,
			Handler:   router,
		}

		fmt.Println("Starting TLS server at port", strings.Split(server.Addr, ":")[1])
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	}()
}

func (h *SwiftHandler) Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello TLS from Test")
}

func (h *SwiftHandler) SwiftNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errlog.Println(err)
		return
	}
	fmt.Println(string(body))
}
