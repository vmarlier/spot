FROM golang:1.22 as BUILD_IMAGE

WORKDIR /opt

COPY . .

# Build
# CGO_ENABLED=0 statically compiles (ie: doesn't link to system libraries)
# GOOS sets the target OS to build for
# GOARCH sets the architecture to build for
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --chown=nonroot:nonroot --from=BUILD_IMAGE /opt/app /app

CMD ["/app"]
