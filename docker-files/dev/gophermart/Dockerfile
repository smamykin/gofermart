FROM golang:1.20-alpine3.17 as dev


ENV UID=1000
ENV GID=1000
ENV USER=www-data
ENV GROUP=www-data

RUN apk add --update make

RUN adduser -S $USER -G $GROUP

RUN set -xe \
    && apk --no-cache add shadow \
    && usermod -u $UID $USER && groupmod -g $GID $GROUP \
    && apk del shadow



WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./pkg ./pkg
COPY ./Makefile ./Makefile
COPY ./build ./build

RUN make build-binary

RUN chown $USER:$GROUP . /tmp

USER $USER

EXPOSE 8080

CMD [ "go", "run", "./cmd/gophermart/main.go"]

FROM alpine:3.17 as prod

COPY --from=dev /app/build/gophermart /gophermart

CMD [ "/gophermart" ]
