# payper

Payper is a simple tool used for forwarding a TCP or UNIX socket from one machine to another.
Using this you can set up a backend that is behind a NAT, with the idea that the backend reaches out to the frontend, which is usually not the case.

# How to use

--- Build the tool and locate it inside \$HOME/go/bin
```
$ go install github.com/nixsum/payper@latest
go: downloading github.com/nixsum/payper v0.0.0-20240128172000-f4c7f8739437
```
--- Generate a listener SSL cert/key pair as well as one for the connector. On the connector the listener certificate has to be present and on the listener the connector certificate has to be present. You can copy them after generating and specify them in the .yml config files. Replace "localhost" with your listener domain name.
```
openssl req -new -nodes -x509 -out ssl/server_cert.pem -keyout ssl/server_key.pem -days 1000 -addext "subjectAltName = DNS:localhost"
openssl req -new -nodes -x509 -out ssl/c0.pem -keyout ssl/k0.pem -days 1000
```
--- Start the listener on your frontend
```
$ go/bin/payper listener.yml
```
--- Start the connector on your backend, in this case I will forward connection to my ssh server. You will recieve a message that it is linked with the frontend
```
$ go/bin/payper connector.yml
2024/01/28 18:09:42 Linked with 45.76.88.19:3333 ; id: 0
```

--- Now you can connect to port 2222 on the frontend and it will be forwarded to the backend.
```
$ nc -v localhost 2222
Ncat: Version 7.93 ( https://nmap.org/ncat )
Ncat: Connected to 127.0.0.1:2222.
SSH-2.0-OpenSSH_8.8
```
