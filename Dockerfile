FROM golang:1.4.2

RUN mkdir -p /go/src/github.com/samertm/sample-goji-app
WORKDIR /go/src/github.com/samertm/sample-goji-app

COPY . /go/src/github.com/samertm/sample-goji-app

RUN ln -sf conf.prod.toml conf.toml

RUN go get -v github.com/samertm/sample-goji-app

CMD ["sample-goji-app"]

EXPOSE 8000
