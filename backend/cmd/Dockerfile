#############################
# STEP 1 build executable binary
#############################
FROM golang:1.19 AS builder

WORKDIR /app

# Copying the whole project
# You might want to change this to copy only what is necessary for the build
COPY backend/ .

RUN go get -v -t -d ./...

# Switching to the directory where your API command resides
WORKDIR /app/cmd/api

# Building the API command
RUN CGO_ENABLED=0 go build -o api .

#############################
# STEP 2 build a small image
##############################
FROM alpine:3.6

# Adding ca-certificates
RUN apk add -U --no-cache ca-certificates

# Copying built binary from builder stage
COPY --from=builder /app/cmd/api/api /usr/bin/api

# Setting up entrypoint
ENTRYPOINT ["/usr/bin/api"]
