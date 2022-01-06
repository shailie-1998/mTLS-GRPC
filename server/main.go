package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"GRPC/unary/sumpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	enableTLS := true
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	// Create the TLS credentials
	if enableTLS {
		fmt.Println("mTLS enabled")
		creds, err := loadTLSCredentials()
		if err != nil {
			log.Fatalf("could not load TLS keys: %s", err)
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		fmt.Println("mTLS disabled")
	}
	s := grpc.NewServer(opts...)
	sumpb.RegisterSumServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server %v", err)
	}

}

// Add returns sum of two integers
func (*server) Add(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	a, b := req.GetNumbers().GetA(), req.GetNumbers().GetB()
	sum := a + b
	return &sumpb.SumResponse{Result: sum}, nil
}
