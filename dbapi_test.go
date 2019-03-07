package dbapi

import (
	"log"
	"net/http/httptest"
	"sync"
)

var (
	serverAddr string
	once       sync.Once
)

func startMockServer() {
	server := httptest.NewServer(nil) // Use default mux
	serverAddr = server.Listener.Addr().String()
	log.Println("Testserver listening on ", serverAddr)
}
