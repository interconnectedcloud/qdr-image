PROJECT_NAME=qdrouterd
DOCKER_REGISTRY=quay.io
PWD=$(shell pwd)

 # This is the latest version of the Qpid Dispatch Router
DISPATCH_VERSION=1.10.0
PROTON_VERSION=0.29.0
PROTON_SOURCE_URL=http://archive.apache.org/dist/qpid/proton/${PROTON_VERSION}/qpid-proton-${PROTON_VERSION}.tar.gz
ROUTER_SOURCE_URL=http://archive.apache.org/dist/qpid/dispatch/${DISPATCH_VERSION}/qpid-dispatch-${DISPATCH_VERSION}.tar.gz

# If a DOCKER_TAG is specified, go ahead and use it.
# if DOCKER_TAG is not specified use the DISPATCH_VERSION as the DOCKER_TAG
ifneq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG_VAL=$(DOCKER_TAG)
else
	DOCKER_TAG_VAL=$(DISPATCH_VERSION)
endif

all: build

build:
	docker build -t qdrouterd-builder:${DOCKER_TAG_VAL} builder
	docker run -ti -v $(PWD):/build:z -w /build qdrouterd-builder:${DOCKER_TAG_VAL} bash build_tarballs ${ROUTER_SOURCE_URL} ${PROTON_SOURCE_URL}

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

	docker build -t quay.io/interconnectedcloud/qdrouterd:${DOCKER_TAG_VAL} .
	docker push quay.io/interconnectedcloud/qdrouterd:${DOCKER_TAG_VAL}

.PHONY: build cleanimage clean push
