PHONY: build
build:
	cd main && go build -o ../bin/oconf

PHONY: install
install:
	cd main && go build -o ${GOPATH}/bin/oconf