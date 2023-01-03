# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.17-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY *.go ./...

# compile application
RUN go build ./...

ADD cmd cmd
ADD pkg pkg
ADD go.mod go.mod
ADD go.sum go.sum

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 3456

# command to be used to execute when the image is used to start a container
CMD [ "go", "run",  "cmd/main.go", "server", "--postgres-host", "host.docker.internal", "--postgres-password", "postgres" ]