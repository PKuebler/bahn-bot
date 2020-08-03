# Start by building the application.
FROM golang:1.14 as build

WORKDIR /go/src/app
ADD . /go/src/app

ARG CGO_ENABLED=1

RUN go get -d -v ./...

RUN go build -o /go/bin/bahn-bot

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/bahn-bot /
CMD ["/bahn-bot", "bot"]