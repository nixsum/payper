package main

import (
	"os"
	"fmt"
	"log"
	"gopkg.in/yaml.v3"
	"sync"
	"time"
)

type SslConf struct {
	Enabled string
	Server_cert string
	Server_key string
	Client_cert string
	Client_key string
	Client_certs []string
}

type Conf struct {
	Sock_loc string
	Sock_pub string
	Tcp_keepalive int
	Ssl SslConf
}

func readconf(filename string) (*map[string]Conf, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make(map[string]Conf)
	
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}

	//TODO: create regex to detect invalid config

	return &data, err
}


func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n" +
		"  payper [mode] [local_socket] [public_socket]\n" +
		"  payper [config_file]\n\n" +
		"Options are:\n" +
		"  mode\t\tthis can be either 'connector' or 'listener'\n" +
		"    connector\tthis is meant to be used on the backend\n" +
		"    listener\tthis is meant to be used on the frontend\n" +
		"  local_socket\t\tthe socket to connect to or listen on, depending on mode\n" +
		"  public_socket\t\tthe public socket to connect to or listen on, this is where the 2 payper instances will meet\n" +
		"    sockets can be in the following formats\n" +
		"      unix:<path to socket>\n" +
		"      tcp4:<IPv4 address or FQDN>\n")
}


func main() {
	switch len(os.Args) {
	case 1:
		usage()
		os.Exit(0)
	case 2:
		if os.Args[1] == "--help" {
			usage()
			os.Exit(0)
		} else {
			conf, err := readconf(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}

			var wg_main sync.WaitGroup

			for k, v := range *conf {
				if k == "listener" {
					wg_main.Add(1)
					go listen(v)
				} else if k == "connector" {
					wg_main.Add(1)
					go connect(v)
				}
				time.Sleep(1 * time.Second)
			}

			wg_main.Wait()
		}
	case 3:
		usage()
		os.Exit(1)
	case 4:
		conf := Conf{Sock_loc: os.Args[2], Sock_pub: os.Args[3], Tcp_keepalive: 0}

		if os.Args[1] == "listen" {
			listen(conf)
		} else if os.Args[1] == "connect" {
			connect(conf)
		} else {
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}

