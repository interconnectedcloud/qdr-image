PROJECT_NAME=qdrouterd
DOCKER_REGISTRY=quay.io
PWD=$(shell pwd)
ROUTER_SOURCE_URL=http://archive.apache.org/dist/qpid/dispatch/1.8.0/qpid-dispatch-1.8.0.tar.gz
PROTON_SOURCE_URL=http://archive.apache.org/dist/qpid/proton/0.28.0/qpid-proton-0.28.0.tar.gz

all: build

build:
	docker build -t qdrouterd-builder:latest builder
	docker run -ti -v $(PWD):/build:z -w /build qdrouterd-builder:latest bash build_tarballs ${ROUTER_SOURCE_URL} ${PROTON_SOURCE_URL}

clean:
	rm -rf proton_build proton_install qpid-dispatch.tar.gz qpid-dispatch-src qpid-proton.tar.gz qpid-proton-src staging build

cleanimage:
	docker image rm -f qdrouterd-builder

push:
# DOCKER_USER and DOCKER_PASSWORD is useful in the CI environment.
# Use the DOCKER_USER and DOCKER_PASSWORD if available
# if not available, assume the user has already logged in
ifneq ($(strip $(DOCKER_USER)$(DOCKER_PASSWORD)),)
	@docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
endif
	docker build -t quay.io/interconnectedcloud/qdrouterd:latest .
	docker push quay.io/interconnectedcloud/qdrouterd:latest

.PHONY: build cleanimage clean push
