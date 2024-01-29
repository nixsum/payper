# payper

Payper is a simple tool used for forwarding a TCP or UNIX socket from one machine to another.
Using this you can set up a backend that is behind a NAT, with the idea that the backend reaches out to the frontend, which is usually not the case.

Note that the connector does not send keepalive messages to the listener, therefore the NAT could decide to close the connection due to inactivity, though on my home router this never happened.
Also note that as of now the connector does not authenticate the listener in any way and vice-versa, Using this on the public internet might not be a good idea in its current state.

# How to use

--- Build the tool and locate it inside \$HOME/go/bin
```
\$ go install github.com/nixsum/payper@latest

go: downloading github.com/nixsum/payper v0.0.0-20240128172000-f4c7f8739437
```
--- Start the listener on your frontend
```
\$ go/bin/payper listen tcp4:0.0.0.0:2222 tcp4:0.0.0.0:3333
```
--- Start the connector on your backend, in this case I will forward connection to my ssh server. You will recieve a message that it is linked with the frontend
```
\$ go/bin/payper connect tcp4:127.0.0.1:22 tcp4:45.76.88.19:3333

2024/01/28 18:09:42 Recieved linkup from 62.73.122.87:61242 ; id: 0
```

--- Now you can connect to port 2222 on the frontend and it will be automatically relayed to the backend.
```
\$ nc -v 45.76.88.19 2222

Ncat: Version 7.93 ( https://nmap.org/ncat )

Ncat: Connected to 45.76.88.19:2222.

SSH-2.0-OpenSSH_8.8
```
