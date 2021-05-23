ARG  BUILDER_IMAGE=golang:1.16
ARG  DISTROLESS_IMAGE=gcr.io/distroless/static
############################
# STEP 1 build executable binary
############################
FROM ${BUILDER_IMAGE} as builder

# Ensure ca-certficates are up to date
RUN update-ca-certificates

WORKDIR $GOPATH/src/mypackage/myapp/

# use modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/goapp .

# WKHTMLTOPDF binaries
FROM surnet/alpine-wkhtmltopdf:3.9-0.12.5-full as wkhtmltopdf

# Base image
FROM ghcr.io/ironpeakservices/iron-alpine/iron-alpine:3.13.2

# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
    libstdc++ \
    libx11 \
    libxrender \
    libxext \
    libssl1.1 \
    ca-certificates \
    fontconfig \
    freetype \
    ttf-dejavu \
    ttf-droid \
    ttf-freefont \
    ttf-liberation \
    ttf-ubuntu-font-family \
    && apk add --no-cache --virtual .build-deps \
    msttcorefonts-installer \
    \
    # Install microsoft fonts
    && update-ms-fonts \
    && fc-cache -f \
    \
    # Clean up when done
    && rm -rf /tmp/* \
    && apk del .build-deps

RUN $APP_DIR/post-install.sh

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /bin/wkhtmltopdf

# Copy our static executable
COPY --from=builder /go/bin/goapp $APP_DIR

WORKDIR $APP_DIR

EXPOSE 8787
USER $APP_USER

ENTRYPOINT ["./goapp"]