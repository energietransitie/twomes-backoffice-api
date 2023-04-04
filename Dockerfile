FROM golang:1.20 as build

WORKDIR /go/src/twomes-api-server

# Create /data folder to be copied later.
RUN mkdir /data

# Download dependencies.
COPY ./go.mod ./go.sum .
RUN go mod download

# Build healthcheck binary.
COPY ./cmd/healthcheck/ ./cmd/healthcheck/
RUN CGO_ENABLED=0 go build -o /go/bin/healthcheck ./cmd/healthcheck/

# Build server binary.
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server/

FROM gcr.io/distroless/static-debian11

# Copy /data folder with correct permissions.
COPY --from=build --chown=nonroot /data /data

# Copy healthcheck binary.
COPY --from=build /go/bin/healthcheck /

# Copy server binary.
COPY --from=build /go/bin/server /

USER nonroot

VOLUME /data

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=1s --start-period=10s --retries=3 \
    CMD ["/healthcheck"]

CMD ["/server"]