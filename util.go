package main

import (
    "log"
    "os"
    "strings"
    "mime"
    "strconv"
)

var logger *log.Logger

// createLogger initializes the logger that will be used to output
// all messages that fmserver generates
func createLogger(){
    logger = log.New(os.Stdout, "FM SERVER: ", log.LstdFlags)
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

