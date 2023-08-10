#############################
# STEP 1 build executable binary
############################
FROM golang:1.19 AS builder

WORKDIR /app

# Copying the whole project
COPY backend/ .

RUN go get -v -t -d ./...

# Switching to the directory where your worker command resides d
WORKDIR /app/worker

# Building the worker command
RUN CGO_ENABLED=0 go build -o worker .

# If you wish to run the command for the built binary, uncomment the line below
# RUN ./worker --version

############################
# STEP 2 build a small image
##############################
FROM alpine:3.6

# Adding ca-certificates
RUN apk add -U --no-cache ca-certificates

# Setting the user environment variable
ENV USER=root

# Copying built binary from builder stage
COPY --from=builder /app/worker/worker /usr/bin/worker

# Setting up entrypoint
ENTRYPOINT ["/usr/bin/worker"]
