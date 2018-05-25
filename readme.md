# Go UDP Hole Punching Chat

[UDP Hole Punching](https://en.wikipedia.org/wiki/UDP_hole_punching) is technique employed in [NAT](https://en.wikipedia.org/wiki/Network_address_translation) for maintaining [User Datagram Protocol](https://en.wikipedia.org/wiki/User_Datagram_Protocol) (UDP) packet streams that traverse the NAT. It is commonly used for establishing bidirectional UDP connections between Internet hosts in private networks using network address translators.

## Start chat server

```go
go run server/main.go
```

## Start chat client

```go
go run client/main.go <port> <server ip:port> <username> <peername>
```

