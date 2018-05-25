package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "encoding/json"
    "time"
)

type ChatRequest struct {
    Action string
    Username string
    Message string
}

type UDPChatClient struct {
    addr *net.UDPAddr
    serverAddr *net.UDPAddr
    peerAddr *net.UDPAddr
    username string
    peername string
    buffer []byte
    connection *net.UDPConn
    retries int
    waitForRetry time.Duration
}

func NewUDPChatClient(port string, server string, username string, peername string) *UDPChatClient {
    log.Print("port: ", port, ", server: ", server, ", username: ", username, ", peer: ", peername)

    client := new(UDPChatClient)

    var err error
    client.serverAddr, err = net.ResolveUDPAddr("udp4", server)
    if err != nil {
        log.Print("Resolve server address failed.")
        log.Fatal(err)
    }

    client.addr, err = net.ResolveUDPAddr("udp4", port)
    if err != nil {
        log.Print("Resolve local addres failed.")
        log.Fatal(err)
    }

    client.username = username
    client.peername = peername

    client.buffer = make([]byte, 2048)

    client.retries = 60
    client.waitForRetry = 1

    client.connection, err = net.ListenUDP("udp", client.addr)
    if err != nil {
        log.Print("Listen UDP failed.")
        log.Fatal(err)
    }

    client.connectToServer()
    client.peerAddr, err = client.getPeerAddress(peername)
    if err != nil {
        log.Fatal(err)
    }

    return client
}

func (self *UDPChatClient) connectToServer() {
    initChatRequest := ChatRequest{
        "New",
        self.username,
        "",
    }

    jsonRequest, err := json.Marshal(initChatRequest)
    if err != nil {
        log.Print("Marshal Register information failed.")
        log.Fatal(err)
    }
    _, err = self.connection.WriteToUDP(jsonRequest, self.serverAddr)
    if err != nil {
        log.Fatal(err)
    }

    log.Print("Waiting for server response...")
    _, _, err = self.connection.ReadFromUDP(self.buffer)
    if err != nil {
        log.Print("Register to server failed.")
        log.Fatal(err)
    }
}

func (self *UDPChatClient) getPeerAddress(peername string) (*net.UDPAddr, error) {
    connectChatRequest := ChatRequest{
        "Get",
        self.username,
        peername,
    }
    jsonRequest, err := json.Marshal(connectChatRequest)
    if err != nil {
        log.Print("Marshal connection information failed.")
        log.Fatal(err)
    }

    var serverResponse ChatRequest
    for i := 0; i < self.retries; i++ {
        self.connection.WriteToUDP(jsonRequest, self.serverAddr)
        n, _, err := self.connection.ReadFromUDP(self.buffer)
        if err != nil {
            log.Print("Get peer address from server failed.")
            log.Fatal(err)
        }
        err = json.Unmarshal(self.buffer[:n], &serverResponse)
        if err != nil {
            log.Print("Unmarshal server response failed.")
            log.Fatal(err)
        }
        if serverResponse.Message != "" {
            break
        }
        time.Sleep(self.waitForRetry * time.Second)
    }

    if serverResponse.Message == "" {
        log.Fatal("Cannot get peer's address")
    }
    log.Print("Peer ", peername, " address: ", serverResponse.Message)
    peerAddr, err := net.ResolveUDPAddr("udp4", serverResponse.Message)
    if err != nil {
        log.Print("Resolve peer addres failed.")
        log.Fatal(err)
    }

    return peerAddr, err
}

func (self *UDPChatClient) start() {
    go self.listen(self.connection)

    self.handleInput()
}

func (self *UDPChatClient) sendMessage(message string, username string) {
    messageRequest := ChatRequest{
        "Chat",
        username,
        message,
    }
    jsonRequest, err := json.Marshal(messageRequest)
    if err != nil {
        log.Print("Error: ", err)
        return
    }
    self.connection.WriteToUDP(jsonRequest, self.peerAddr)
}

func (self *UDPChatClient) handleInput() {
    for {
        fmt.Print("Input message: ")
        message := make([]byte, 2048)
        fmt.Scanln(&message)
        self.sendMessage(string(message), self.username)
    }
}

func (self *UDPChatClient) recvMessage() (ChatRequest, error) {
    buffer := make([]byte, 2048)
    n, _, err := self.connection.ReadFromUDP(buffer)
    if err != nil {
        log.Print(err)
        return ChatRequest{}, err
    }

    var message ChatRequest
    err = json.Unmarshal(buffer[:n], &message)
    if err != nil {
        log.Print(err)
        return ChatRequest{}, err
    }

    return message, nil
}

func (self *UDPChatClient) listen(conn *net.UDPConn) {
    for {
        message, err := self.recvMessage()
        if err != nil {
            log.Print(err)
            continue
        }

        fmt.Println("\n", message.Username, ":", message.Message)
        fmt.Print("Input message: ")
    }
}

func main() {
    if len(os.Args) < 5 {
        log.Fatal("Usage: ", os.Args[0], " port serverAddr username peername")
    }

    port := os.Args[1]
    serverAddr := os.Args[2]
    username := os.Args[3]
    peer := os.Args[4]

    log.Print("port: ", port, ", serverAddr: ", serverAddr, ", username: ", username, ", peer: ", peer)

    client := NewUDPChatClient(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
    client.start()
}

