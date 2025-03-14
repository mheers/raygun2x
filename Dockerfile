ARG go="golang:1.24.1-alpine3.21"
ARG base="alpine:3.21"

FROM --platform=$BUILDPLATFORM ${go} AS builder

RUN apk add --no-cache bash git

ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Copy the code from the host and compile it
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

ADD . ./

RUN [ "$(uname)" = Darwin ] && system=darwin || system=linux; \
    ./ci/go-build.sh --os ${system} --arch $(echo $TARGETPLATFORM  | cut -d/ -f2)

# final stage
FROM ${base}
WORKDIR /app
COPY --from=builder /go/src/app/goapp /usr/local/bin/raygun2x
ENTRYPOINT ["/usr/local/bin/raygun2x"]
CMD ["help"]
