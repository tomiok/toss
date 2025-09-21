ARG  BUILDER_IMAGE=golang:1.24.4-alpine
############################
# STEP 1 build executable binary
############################
FROM ${BUILDER_IMAGE} as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create appuser
ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/toss/
# use modules
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/toss ./cmd/api/main
############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/toss /go/bin/toss

COPY static /static/
COPY migrations /migrations/
COPY platform/templates /platform/templates/

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 9000
# Run the binary.

ENTRYPOINT ["/go/bin/toss"]