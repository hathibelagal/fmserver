// Author: Hathibelagal
// FM Server

// This is a simple HTTP server that can handle GET requests, and serve
// static files from a directory
package main

import (
    "net"
    "strings"
    "io/ioutil"
    "time"
    "path/filepath"
)

type ServerProperties struct{
    port string
    totalConnections, errors int
    directory string
    openConnections int
}
var props ServerProperties

// createListener is used to generate a tcp based listener on a
// given port
func createListener(port string) net.Listener {
    listener, err := net.Listen("tcp", port)
    if err != nil {
        logger.Fatal(err)
    }
    return listener
}

// handleConnections tries to accept connections to the listener
// over an endless loop
func handleConnections(listener net.Listener) {
    logger.Println("Listening on port : " + props.port[1:])
    for {
        connection, err := listener.Accept()
        if err != nil {
            logger.Print(err)
        }
        props.openConnections += 1

        // Process this connection using a new goroutine
        // so that other connections don't have to wait
        go func() {
            responseComplete := make (chan bool, 1)
            timedOut := false

            defer func() {
                connection.Close()
                props.openConnections -= 1
            }()

            go answer(connection, responseComplete, &timedOut)

            // Timeout if the response is not generated within
            // 5 seconds
            select {
                case <-responseComplete:
                    return
                case <-time.After(5 * time.Second):
                    logger.Println("Response timed out.")
                    timedOut = true
                    connection.Write(generateError(503))
                    return
            }
        }()
    }
}

// serve tries to understand the request, and generate an appropriate response
// It returns nil if it doesn't understand the request
// It returns valid HTTP responses otherwise
func serve(request string) []byte {
    parts := strings.Fields(request)

    // We need atleast two words to understand the request, and the first
    // word should be GET, as we are handling only the GET verb
    if (len(parts) < 2 || parts[0] != "GET") {
        return nil
    }

    // The second word is the filename
    filename := parts[1]

    // If it is a root request, try to server index.html
    if (strings.HasSuffix(filename, "/")) {
        filename = filename+ "/index.html"
    }
    
    filename = props.directory + filename

    if (strings.HasPrefix(filepath.Dir(filename), props.directory) == false) {
        return generateError(404)
    }

    body, err := ioutil.ReadFile(filename)
    if (err != nil) {
        return generateError(404)
    }
    header := generateSuccessHeader(filename, len(body))
    return append(header, body...)
}

// answer reads the request made by the client and passes it
// on to serve(), to generate a response
func answer(connection net.Conn, responseComplete chan<- bool, timedOut *bool) {
    buffer := make([]byte, 4096) 
    bytesRead, err := connection.Read(buffer)
    if err != nil {
        logger.Print(err)
    }
    request := string(buffer[:bytesRead])
    response := serve(request)
    props.totalConnections += 1

    if (*timedOut) {
        logger.Println("Too late")
        return
    }

    if (response == nil){
        logger.Print("Invalid request. Ignoring")
        props.errors += 1
        responseComplete <- true
        return
    }
    connection.Write(response)
    responseComplete <- true
}
