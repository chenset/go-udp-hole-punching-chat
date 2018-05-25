package main

// https://gist.github.com/reterVision/33a72d70194d4a3c272e

import (
    "net"
    "log"
    "encoding/json"
    "fmt"
)

var userIP map[string]string

type ChatRequest struct {
    Action string
    Username string
    Message string
}

type UDPChatServer struct {
    userIP map[string]string
    udpAddr *net.UDPAddr
    connection *net.UDPConn
}

func NewUDPChatServer(port string) *UDPChatServer {
    server := new(UDPChatServer)

    server.userIP = map[string]string{}
    resolvedPort, err := net.ResolveUDPAddr("udp4", port)
    if err != nil {
        log.Fatal(err)
    }
    server.udpAddr = resolvedPort

    return server
}

func (self *UDPChatServer) listen() {
    log.Print("Start server")

    conn, err := net.ListenUDP("udp", self.udpAddr)
    self.connection = conn
    if err != nil {
        log.Fatal(err)
    }

    for {
        self.handleClient(conn)
    }
}

func (self *UDPChatServer) handleClient(conn *net.UDPConn) {
    log.Print("Handle client")

    var buf [2048]byte

    n, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
        return
    }

    var chatRequest ChatRequest
    err = json.Unmarshal(buf[:n], &chatRequest)
    if err != nil {
        log.Print(err)
        return
    }

    log.Print("Request: ", chatRequest)

    switch chatRequest.Action {
    case "New":
        self.onChatRequestActionIsNew(chatRequest, addr)
    case "Get":
        self.onChatRequestActionIsGet(chatRequest, addr)
    }
    fmt.Println("User table:", self.userIP)
}

func (self *UDPChatServer) onChatRequestActionIsNew(request ChatRequest, addr *net.UDPAddr) {
        remoteAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
        fmt.Println(remoteAddr, "connecting")
        self.userIP[request.Username] = remoteAddr

        messageRequest := ChatRequest{
            "Chat",
            request.Username,
            remoteAddr,
        }
        jsonRequest, err := json.Marshal(&messageRequest)
        if err != nil {
            log.Print(err)
        }
        self.connection.WriteToUDP(jsonRequest, addr)
}

func (self *UDPChatServer) onChatRequestActionIsGet(request ChatRequest, addr *net.UDPAddr) {
        peerAddr := ""
        if _, ok := self.userIP[request.Message]; ok {
            peerAddr = self.userIP[request.Message]
        }

        messageRequest := ChatRequest{
            "Chat",
            request.Username,
            peerAddr,
        }
        jsonRequest, err := json.Marshal(&messageRequest)
        if err != nil {
            log.Print(err)
        }
        _, err = self.connection.WriteToUDP(jsonRequest, addr)
        if err != nil {
            log.Print(err)
        }
}

func main() {
    server := NewUDPChatServer(":9999")
    server.listen()
}

