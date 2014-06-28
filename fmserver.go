package main

import (
    "net"
    "log"
    "os"
    "strings"
    "io/ioutil"
)

var logger *log.Logger

type Stats struct{
    totalConnections, errors int64
}
var stats Stats

func createLogger(){
    logger = log.New(os.Stdout, "FM SERVER: ", log.LstdFlags)
    logger.Print("Ready")
}

func createListener(port string) net.Listener {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        logger.Fatal(err)
    }
    return listener
}

func handleConnections(listener net.Listener) {
    for {
        connection, err := listener.Accept()
        if err != nil {
            logger.Print(err)
        }
        go answer(connection)
    }
}

func getContentType(filename string) string {
    knownContentTypes := make(map[string]string)
    knownContentTypes["png"] = "image/png"
    knownContentTypes["jpg"] = "image/jpeg"
    knownContentTypes["txt"] = "text/plain"
    knownContentTypes["html"] = "text/html"

    parts := strings.Split(filename, ".")
    extension := parts[len(parts)-1]
    return knownContentTypes[extension]
}

func serve(request string) []byte {
    parts := strings.Fields(request)
    if (len(parts) < 2 || parts[0] != "GET") {
        return nil
    }
    filename := parts[1]
    body, err := ioutil.ReadFile("/tmp/" + filename)
    if (err != nil) {
        return []byte("ERROR")
    }
    header := []byte("HTTP/1.1 200 OK\nContent-type: "+getContentType(filename)+"\n\n")
    return append(header, body...)
}

func answer(connection net.Conn) {
    buffer := make([]byte, 4096)
    bytesRead, err := connection.Read(buffer)
    if err != nil {
        logger.Print(err)
    }
    request := string(buffer[:bytesRead])
    response := serve(request)
    stats.totalConnections += 1
    if (response == nil){
        logger.Print("Invalid request. Ignoring")
        stats.errors += 1
        connection.Close()
        return
    }
    connection.Write(response)
    connection.Close()
}

func main() {
    createLogger()
    listener := createListener(":8080")
    handleConnections(listener)
}
