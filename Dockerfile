FROM golang:1.4.2

RUN mkdir -p /go/src/github.com/samertm/sample-gosu-app
WORKDIR /go/src/github.com/samertm/sample-gosu-app

COPY . /go/src/github.com/samertm/sample-gosu-app

RUN ln -sf conf.prod.toml conf.toml

RUN go get -v github.com/samertm/sample-gosu-app

CMD ["sample-gosu-app"]

EXPOSE 8000
