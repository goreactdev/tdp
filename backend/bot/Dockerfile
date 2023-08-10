#############################
# STEP 1 build executable binary
############################
FROM golang:1.19 AS builder

WORKDIR /app

# Copying the whole project
COPY backend/ .

RUN go get -v -t -d ./...

# Switching to the directory where your bot command resides
WORKDIR /app/bot

# Building the bot command
RUN CGO_ENABLED=0 go build -o bot .

# If you wish to run the command for the built binary, uncomment the line below
# RUN ./bot --version

############################
# STEP 2 build a small image
##############################
FROM alpine:3.6

# Adding ca-certificates
RUN apk add -U --no-cache ca-certificates

# Setting the user environment variable
ENV USER=root

# Copying built binary from builder stage
COPY --from=builder /app/bot/bot /usr/bin/bot

# Setting up entrypoint
ENTRYPOINT ["/usr/bin/bot"]
