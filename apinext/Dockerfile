# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
# TODO(madadam): Is there an ubuntu-based Go image?
#FROM golang
FROM gcr.io/google_appengine/golang
ENV GO15VENDOREXPERIMENT 1

# Install & build glide
RUN git clone https://github.com/Masterminds/glide.git /go/src/github.com/Masterminds/glide
WORKDIR /go/src/github.com/Masterminds/glide
RUN make bootstrap && make install

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/clarifai/go/apinext
WORKDIR /go/src/github.com/clarifai/go/apinext

# Build the project inside the container.
RUN glide install && go install github.com/clarifai/go/apinext

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/apinext

# Document that the service listens on port 8080.
EXPOSE 8080
