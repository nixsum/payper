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

func accept(listener_loc, listener_pub net.Listener, wg_listen *sync.WaitGroup, id int) {
	defer wg_listen.Done()
	conn_pub, err := listener_pub.Accept() // Will block here
	if err != nil {
		log.Fatal("Listener error: %s", err)	
	}
	defer conn_pub.Close()
	log.Printf("Recieved linkup from %s ; id: %d\n", conn_pub.RemoteAddr(), id)

	// Accept connection on the local listener
	conn_loc, err := listener_loc.Accept() // Will block here
	if err != nil {
		log.Fatal(err)
	}
	defer conn_loc.Close()
	log.Printf("Recieved connection from %s ; id: %d\n", conn_loc.RemoteAddr(), id)

	wg_listen.Add(1)
	go accept(listener_loc, listener_pub, wg_listen, id + 1) // Call self to accept a new connection

	var wg_copy sync.WaitGroup
	wg_copy.Add(1)
	
	// go conn_copy(&conn_loc, &conn_pub, &wg_copy)
	// go conn_copy(&conn_pub, &conn_loc, &wg_copy)
	go func () {
		io.Copy(conn_loc, conn_pub)
	}()
	go func () {
		defer wg_copy.Done()
		io.Copy(conn_pub, conn_loc)
	}()
	wg_copy.Wait() // Wait until both copies complete, they will copy until they reach EOF
	log.Printf("Closed linkup and connection with id %d\n", id)
}


func listen(conf Conf) {
	sock_loc_type := conf.Sock_loc[:4]
	sock_loc_string := conf.Sock_loc[5:]
	sock_pub_type := conf.Sock_pub[:4]
	sock_pub_string := conf.Sock_pub[5:]

	var listener_pub net.Listener

	if conf.Ssl.Enabled == "yes" {
		// Read the server certificate and key
		certFile := conf.Ssl.Server_cert
		keyFile := conf.Ssl.Server_key

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}
		// Read the approved client certificates
		clientCertPool := x509.NewCertPool()
		for _, filename := range conf.Ssl.Client_certs {
			file, err := os.ReadFile(filename)
			if err != nil {
				log.Printf("Unable to open client certificate %s , continuing...", filename)
			}
			clientCertPool.AppendCertsFromPEM(file)
		}

		tlsconfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs: clientCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}

		var err_listen error
		listener_pub, err_listen = tls.Listen(sock_pub_type, sock_pub_string, tlsconfig)
		if err_listen != nil {
			log.Fatal(err_listen)
		}
	} else {
		var err_listen error
		listener_pub, err_listen = net.Listen(sock_pub_type, sock_pub_string)
		if err_listen != nil {
			log.Fatal(err_listen)
		}
	}
	defer listener_pub.Close()

	listener_loc, err := net.Listen(sock_loc_type, sock_loc_string)
	if err != nil {
		log.Fatal(err)
	}
	defer listener_loc.Close()

	var wg_listen sync.WaitGroup
	id := 0
	wg_listen.Add(1)
	go accept(listener_loc, listener_pub, &wg_listen, id)
	wg_listen.Wait()
}