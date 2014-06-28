package main

import (
    "net"
    "log"
    "os"
    "strings"
    "io/ioutil"
    "mime"
    "strconv"
    "time"
    "path/filepath"
)

var logger *log.Logger

type ServerProperties struct{
    port string
    totalConnections, errors int64
    directory string
}
var props ServerProperties

// createLogger initializes the logger that will be used to output
// all messages that fmserver generates
func createLogger(){
    logger = log.New(os.Stdout, "FM SERVER: ", log.LstdFlags)
}

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

        // Process this connection using a new goroutine
        // so that other connections don't have to wait
        go func() {
            responseComplete := make (chan bool, 1)
            timedOut := false

            go answer(connection, responseComplete, &timedOut)

            // Close the connection if the response is not generated
            // within 5 seconds
            select {
                case <-responseComplete:
                    return
                case <-time.After(5 * time.Second):
                    logger.Println("Response timed out.")
                    timedOut = true
                    connection.Write(generateError(503))
                    connection.Close()
                    return
            }
        }()
    }
}

// getContentType generates the Content-type header based on
// the extension of the file
func getContentType(filename string) string {    
    parts := strings.Split(filename, ".")
    extension := parts[len(parts)-1]
    return "Content-type: " + mime.TypeByExtension("." + extension)
}

// generateSuccessHeader generates a 200 header for the browser to
// understand the response
func generateSuccessHeader(filename string, bytes int) []byte {
    output := "HTTP/1.1 200 OK" + "\n"
    output = output + "Server: FMServer" + "\n"
    output = output + getContentType(filename) + "\n"
    output = output + "Content-length: " + strconv.Itoa(bytes) + "\n"
    output = output + "\n\n"
    return []byte(output)
}

// generateError generates an appropriate error page for the browser to
// understand the response
func generateError(status int) []byte {
    output := "HTTP/1.1 "
    errorMessage := ""
    switch(status){
        case 404: errorMessage = "404 NOT FOUND"
        case 503: errorMessage = "503 SERVICE UNAVAILABLE"
        default: errorMessage = "500 INTERNAL SERVER ERROR"
    }
    output = output + errorMessage + "\n"
    output = output + "Server: FMServer" + "\n"
    output = output + "Content-type: text/html" + "\n\n"
    output = output + "<html><b>" + errorMessage + "</b></html>"
    return []byte(output)
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
    if (filename == "/") {
        filename = "/index.html"
    }
    
    filename = props.directory + filename

    if (filepath.Dir(filename) != props.directory) {
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
        connection.Close()
        responseComplete <- true
        return
    }
    connection.Write(response)
    connection.Close()
    responseComplete <- true
}

func processCommandLine(){
    var err error
    if (len(os.Args) != 2){
        logger.Fatal ("No port specified")
    }
    props.directory, err = filepath.Abs(filepath.Dir(os.Args[0]))
    if (err != nil) {
        logger.Fatal(err)
    }
    _, err = strconv.Atoi(os.Args[1])
    if (err != nil) {
        logger.Fatal ("Invalid port number specified")
    }
    props.port = ":" + os.Args[1]
    logger.Println("Serving directory : " + props.directory)
}

func main() {
    createLogger()
    processCommandLine()
    listener := createListener(props.port)
    handleConnections(listener)
}
