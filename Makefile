PROJECT_NAME=qdrouterd
PWD=$(shell pwd)
ROUTER_SOURCE_URL=http://archive.apache.org/dist/qpid/dispatch/1.7.0/qpid-dispatch-1.7.0.tar.gz
PROTON_SOURCE_URL=http://archive.apache.org/dist/qpid/proton/0.27.1/qpid-proton-0.27.1.tar.gz

all: build

build:
	docker build -t qdrouterd-builder:latest builder
	docker run -ti -v $(PWD):/build:z -w /build qdrouterd-builder:latest bash build_tarballs ${ROUTER_SOURCE_URL} ${PROTON_SOURCE_URL}

clean:
	rm -rf proton_build proton_install qpid-dispatch.tar.gz qpid-dispatch-src qpid-proton.tar.gz qpid-proton-src staging build

.PHONY: build clean
