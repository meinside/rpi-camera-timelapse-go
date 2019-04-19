# Dockerfile for Golang application

FROM balenalib/raspberrypi3-debian-golang:latest

# Working directory outside $GOPATH
WORKDIR /src

# Copy go module files and download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy source files
COPY ./ ./

# Build source files statically
RUN CGO_ENABLED=0 go build \
		-installsuffix 'static' \
		-o /app \
		.

# for rpi binaries
RUN apt-get update -y && \
		apt-get install -y apt-utils libraspberrypi-bin

# Copy config file
COPY ./config.json /

# Open ports (if needed)
#EXPOSE 8080
#EXPOSE 80
#EXPOSE 443

# Entry point for the built application
ENTRYPOINT ["/app"]
