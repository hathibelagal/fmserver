# FM Server
## A custom server built with Go
### Developed by Hathi

This is a simple HTTP server that can handle only GET requests, and serve
static files from a directory.

To serve the contents of a directory, go to that directory, and just type

`fmserver 8080`

Of course, you can change 8080 to any port you want. **Make sure fmserver is in your PATH before you do this**

The complete installation procedure on Linux is as follows:

- Install **go** and **git** on your machine
- `mkdir /tmp/testproject #This could be any other directory`
- `export GOPATH=/tmp/testproject`
- `cd $GOPATH`
- `go get github.com/hathibelagal/fmserver`
- `cd $GOPATH/src/github.com/hathibelagal/fmserver`
- `go install`

- *Installation is now complete.*
- `export PATH=$PATH:$GOPATH/bin`
- *To serve files from your /tmp directory, the procedure is*
- `cd /tmp`
- `fmserver 8080`


