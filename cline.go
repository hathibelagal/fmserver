package main

import (
    "os"
    "os/signal"
    "strconv"
    "net"
    "path/filepath"
)

// processCommandLine determines the directory from which the server was
// started. It also gets the port number on which the server will accept
// connections
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

// handleSIGINT captures the SIGINT interrupt, and shuts down the server.
// After shutting down, it also prints various stats about the server, like
// the total number of requests it served, the number of invalid requests etc.
func handleSIGINT(listener net.Listener) {
    sigint_ch := make(chan os.Signal, 1)
    signal.Notify(sigint_ch, os.Interrupt)

    go func() {
        <-sigint_ch
        logger.Println("Shutting down...")
        listener.Close()
        logger.Println("Aborted requests : " + strconv.Itoa(props.openConnections))
        logger.Println("Total requests served: " + strconv.Itoa(props.totalConnections))
        logger.Println("Erroneous requests : " + strconv.Itoa(props.errors))
        logger.Println("Bye")
        os.Exit(0)
    }()
}

// main is where is all begins. It calls all the necessary functions
// in the right order to get the server started
func main() {
    createLogger()
    processCommandLine()
    listener := createListener(props.port)
    handleSIGINT(listener)
    handleConnections(listener)
}
