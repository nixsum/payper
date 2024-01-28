package main

import (
	"net"
	"log"
	"io"
	"sync"
)


func linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string string, wg_connect *sync.WaitGroup, id int) {
	defer wg_connect.Done()
	// Connect to the public listener
	conn_pub, err := net.Dial(sock_pub_type, sock_pub_string)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Linked with %s ; id: %d\n", sock_pub_string, id)
	defer conn_pub.Close()
	// This is where authentiaction with the remote server should occur.
	// Initiate connection to the local listener once data is recieved on the public connection
	buf := make([]byte, 1500)
	n, err := conn_pub.Read(buf) // Will block here
	if err != nil {
		log.Fatal(err)
	}
	buf = buf[:n]

	// connecto to the local listener
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
	go linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string, wg_connect, id + 1)

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


func accept(listener_loc, listener_pub net.Listener, wg_listen *sync.WaitGroup, id int) {
	defer wg_listen.Done()
	conn_pub, err := listener_pub.Accept() // Will block here
	if err != nil {
		log.Fatal(err)
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


func Listen(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string string) {
	// Start the public and local listeners
	listener_pub, err := net.Listen(sock_pub_type, sock_pub_string)
	if err != nil {
		log.Fatal(err)
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


func Connect(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string string) {
	var wg_connect sync.WaitGroup
	id := 0
	wg_connect.Add(1)
	go linkup(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string, &wg_connect, id)
	wg_connect.Wait()
}