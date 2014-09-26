#!/bin/bash

go test -covermode=count -coverprofile=site.coverage.out .
go test -covermode=count -coverprofile=links.coverage.out ./links
go test -covermode=count -coverprofile=db.coverage.out ./db
go test -covermode=count -coverprofile=urlshortener.coverage.out ./urlshortener

echo "mode: count" > coverage.out
test -f site.coverage.out && tail -n+2 site.coverage.out >> coverage.out
test -f links.coverage.out && tail -n+2 links.coverage.out >> coverage.out
test -f db.coverage.out && tail -n+2 db.coverage.out >> coverage.out
test -f urlshortener.coverage.out && tail -n+2 urlshortener.coverage.out >> coverage.out

if [ "$(cat coverage.out | wc -l)" -gt 1 ]; then
  $(go env GOPATH | awk 'BEGIN{FS=":"} {print $1}')/bin/goveralls -coverprofile=coverage.out -service=codeship -repotoken AA5OvxIquBCzVvXj0NNge0LpAakBiLuiC
fi

rm *.out
