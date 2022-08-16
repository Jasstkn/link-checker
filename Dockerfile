FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

# hadolint ignore=DL3019
RUN apk add git=2.36.2-r0

ARG LDFLAGS

RUN CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o linkchecker cmd/linkchecker/main.go

# FROM alpine
# hadolint ignore=DL3006
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/linkchecker linkchecker

ENTRYPOINT ["/linkchecker"]
