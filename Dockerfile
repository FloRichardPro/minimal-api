#### BUILD STAGE
FROM golang:1.18-alpine AS builder
RUN apk add --update --no-cache build-base
ENV GOPRIVATE=git.wegetfunded.tech/*
WORKDIR $GOPATH/src/wegetfunded/

# Copy project files into the builder
COPY . .

RUN go build -a -mod=vendor -ldflags '-linkmode external -w -s -extldflags "-static"' -o ./minimal-api

#### IMAGE DEFINITION
FROM alpine:3
# Port used by the application, do not override this port
ENV PORT=9090
# Upgrade system dependencies for security patches or bug fixing
RUN apk update && apk upgrade
WORKDIR /app
# Copy the programm with for uid 1000 and gid 1000
COPY --chown=1000:1000 --from=builder /go/src/wegetfunded/minimal-api minimal-api

# Expose port, it is just documentation for ops
EXPOSE ${PORT}
# Define a healthcheck to know if container is healthy with curl by getting a response from /status endpoint and parsing result with jq
# HEALTHCHECK --interval=10s --timeout=10s --start-period=20s --retries=3 \ 
# CMD (curl -s http://localhost:${PORT}/status || echo '{"healthy":null}') | jq -e .healthy
# Run with a non-root user
USER 1000:1000

CMD ["./minimal-api"]