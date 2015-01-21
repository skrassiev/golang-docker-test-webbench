# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.3

# Copy the local package files to the container's workspace.
#ADD src /go/src
ADD . /go/src/github.com/skrassiev/golang-docker-test-webbench/

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/skrassiev/golang-docker-test-webbench

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/golang-docker-test-webbench"]

# Document that the service listens on port 8080.
EXPOSE 8080
EXPOSE 9001
