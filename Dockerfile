FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:latest AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /pilpres

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static-debian11 AS release

WORKDIR /

COPY --from=build /pilpres /pilpres

USER nonroot:nonroot

EXPOSE 8080

CMD ["/pilpres"]
