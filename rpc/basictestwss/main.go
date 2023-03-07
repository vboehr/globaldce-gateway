package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
	//"context"
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/gorilla/rpc/v2"
)

func main() {
	//////////////////////////////////////////

	// Generate a new RSA private key.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Create a self-signed X.509 certificate.
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Fatal(err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	// Load the certificate and private key.
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatal(err)
	}

	//////////////////////////////////////////////
	// Create a new Gorilla Mux router.
	router := mux.NewRouter()

	// Register the WebSocket endpoint.
	router.HandleFunc("/ws", handleWebSocket)

	// Create a new HTTPS server.
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	// Start the server.
	log.Println("Starting server on https://localhost:8080...")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}

// handleWebSocket handles incoming WebSocket connections.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket connection.
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	///////////////////////////////////////////
	jsonrpcserver := rpc.NewServer()
	jsonrpcserver.RegisterService("Arith", arith)
	///////////////////////////////////////////

	// Handle the WebSocket connection.
	for {
		/*
		// Read a message from the WebSocket connection.
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from WebSocket connection:", err)
			break
		}

		// Handle the message here. You can implement your own RPC protocol here.
		log.Printf("Received message: %s", message)

		// Send a response back to the client.
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, client!")); err != nil {
			log.Println("Failed to send message to WebSocket connection:", err)
			break
		}*/
		// Read the next message from the WebSocket connection.
		_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
		
			// Handle the JSON-RPC message and send the response back.
			resp, err := jsonrpcserver.HandleBytes(msg)
			if err != nil {
				log.Println(err)
				break
			}
			if resp != nil {
				err = conn.WriteMessage(websocket.TextMessage, resp)
				if err != nil {
					log.Println(err)
					break
				}
			}
	}

	// Close the WebSocket connection.
	conn.Close()
}
