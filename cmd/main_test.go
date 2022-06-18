package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestJwtHTTPServer(t *testing.T) {
	server := NewJwtServer("", "")
	// Start the test server on random port.
	go server.runHTTP("localhost:0")
	// Prepare the HTTP request.
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/jwtkeys", <-server.httpPort), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected to get %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestJwtHTTPSServer(t *testing.T) {
	var (
		serverKey  = "../assets/private.pem"
		serverCert = "../assets/publickey.crt"
	)

	caCert, err := os.ReadFile(serverCert)
	if err != nil {
		t.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// creating https client with client certificate and certificate authority
	httpsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	server := NewJwtServer(serverCert, serverKey)

	// Start the test server on port 8443.
	go server.runHTTPS(":8443")

	httpsReq, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%d/jwtkeys", <-server.httpsPort), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	resp, err := httpsClient.Do(httpsReq)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected to get %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
