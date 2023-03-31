FROM golang:1.20 as build

WORKDIR /go/src/twomes-api-server

COPY ./go.mod ./go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server ./cmd/server/

# Create /data folder to be copied later.
RUN mkdir /data

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/server /

# Copy /data folder with correct permissions.
COPY --from=build --chown=nonroot /data /data

USER nonroot

VOLUME /data

EXPOSE 8080

CMD ["/server"]