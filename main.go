package main

import (
	"os"
	"fmt"
	"log"
)


func help() {
		fmt.Println("Usage: payper <mode> <local socket> <public socket>")
		fmt.Println("<mode>: There are 2 options here:") 
		fmt.Println("\t\"listen\" is to be used on your frontend with <local socket> being what your frontend expects your backend to be and <public socket> what the payper instance on your backend has as public socket.")
		fmt.Println("\t\"connect\" is to be used on your backend with <local socket> being where your backend expects connections and <public socket> what the payper instance on your backend has as public socket.")
		fmt.Println("The sockets can be in the following format:")
		fmt.Println("\tunix:<path to socket>\n\ttcp4:<IPv4 address>:<TCP port>")
		os.Exit(0);
}

func main() {
	if len(os.Args) == 1 {
		help()
	}

	mode := os.Args[1]

	if mode == "--help" {
		help()
	} 

	if len(os.Args) != 4 {
		log.Fatal("Not enough command line args!")
	}

	sock_loc_type := os.Args[2][:4]
	sock_loc_string := os.Args[2][5:]
	sock_pub_type := os.Args[3][:4]
	sock_pub_string := os.Args[3][5:]

	if mode == "listen" {
		Listen(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string)
	} else if mode == "connect" {
		Connect(sock_pub_type, sock_pub_string, sock_loc_type, sock_loc_string)
	}
}

