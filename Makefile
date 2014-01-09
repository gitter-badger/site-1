buid:
	go build

build-production:
	go build -ldflags "-s"

build-debug:
	go build -gcflags "-N -l"

clean:
	rm site
