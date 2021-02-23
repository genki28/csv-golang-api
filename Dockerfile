FROM golang:1.12.6-alpine

WORKDIR /app
ADD . /app

RUN apk add --no-cache alpine-sdk git && go get github.com/oxequa/realize

EXPOSE 8000
CMD ["realize", "start", "--run"]