FROM golang

COPY . /go/src/github.com/benjaminmisell/pupur_server
ENV CGO_ENABLED=0
RUN cd /go/src/github.com/benjaminmisell/pupur_server && go get -v && go install

FROM scratch
ENV GOPATH=/go
COPY vendor /go/vendor
COPY --from=0 /go/bin/pupur_server /go/
WORKDIR /go
CMD ["./pupur_server"]
