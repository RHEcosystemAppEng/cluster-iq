## Build
####################
FROM golang:1.18.3 AS builder

# Code copy
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aws_scanner ./src/aws_scanner/scanner.go

## Run
####################
FROM scratch

# Arguments
ARG VERSION="0.1-alpha"

# Labels
LABEL version=$VERSION
LABEL description="Openshift Inventory AWS Scanner"

# Binary
COPY --from=builder /app/aws_scanner /aws_scanner
# CA Certs to access AWS domains
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Default config
ENV CIQ_DB_HOST="redis"
ENV CIQ_DB_PORT="6379"
ENV CIQ_DB_PASSWORD=""

ENTRYPOINT ["/aws_scanner"]
