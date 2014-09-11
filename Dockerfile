FROM octohost/go-1.2

RUN mkdir -p $GOPATH/src/github.com/txgruppi/site
ADD . $GOPATH/src/github.com/txgruppi/site
RUN cd $GOPATH/src/github.com/txgruppi/site; go get ./...

EXPOSE 5000

CMD cd $GOPATH/src/github.com/txgruppi/site; go run main.go
