package main

import (
	"net"
	"log"
	"io"
	"os"
	"sync"
	"crypto/tls"
	"crypto/x509"
)


func linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string string, tlsconfig *tls.Config, wg_connect *sync.WaitGroup, id int) {
	defer wg_connect.Done()
	// Connect to the public listener
	var conn_pub net.Conn
	if tlsconfig != nil {
		var err error
		conn_pub, err = tls.Dial(sock_pub_type, sock_pub_string, tlsconfig)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		var err error
		conn_pub, err = net.Dial(sock_pub_type, sock_pub_string)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	log.Printf("Linked with %s ; id: %d\n", sock_pub_string, id)
	defer conn_pub.Close()
	
	buf := make([]byte, 1500)
	n, err := conn_pub.Read(buf) // Will block here
	if err != nil {
		log.Fatal(err)
	}
	buf = buf[:n]

	// connect to the local listener
	conn_loc, err := net.Dial(sock_loc_type, sock_loc_string)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to %s ; id: %d\n", sock_loc_string, id)
	defer conn_loc.Close()

	_, err = conn_loc.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	wg_connect.Add(1)
	go linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string, tlsconfig, wg_connect, id + 1)

	var wg_copy sync.WaitGroup
	wg_copy.Add(1)
	
	// go conn_copy(&conn_loc, &conn_pub, &wg_copy)
	// go conn_copy(&conn_pub, &conn_loc, &wg_copy)
	go func () {
		defer wg_copy.Done()
		io.Copy(conn_loc, conn_pub)
	}()
	go func () {
		io.Copy(conn_pub, conn_loc)
	}()
	wg_copy.Wait() // Wait until both copies complete, they will copy until they reach EOF
	log.Printf("Closed linkup and connection with id %d\n", id)
}


func connect(conf Conf) {
	sock_loc_type := conf.Sock_loc[:4]
	sock_loc_string := conf.Sock_loc[5:]
	sock_pub_type := conf.Sock_pub[:4]
	sock_pub_string := conf.Sock_pub[5:]

	var tlsconfig *tls.Config = nil

	if conf.Ssl.Enabled == "yes" {
		// Load the server certificate and use it as a root certificate
		certs := x509.NewCertPool()

		serverCert, err := os.ReadFile(conf.Ssl.Server_cert)
		if err != nil {
			log.Fatal(err)
		}
		certs.AppendCertsFromPEM(serverCert)
		// Load client sert/key pair
		certFile := conf.Ssl.Client_cert
		keyFile := conf.Ssl.Client_key

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		tlsconfig = &tls.Config{
			RootCAs: certs,
			Certificates: []tls.Certificate{cert},
		}
	}

	var wg_connect sync.WaitGroup
	id := 0
	wg_connect.Add(1)
	go linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string, tlsconfig, &wg_connect, id)
	wg_connect.Wait()
}