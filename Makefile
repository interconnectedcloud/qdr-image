DOCKER := docker
PROJECT_NAME=qdrouterd
DOCKER_REGISTRY=quay.io
DOCKER_ORG=interconnectedcloud
PWD=$(shell pwd)

 # This is the latest version of the Qpid Dispatch Router
DISPATCH_VERSION=1.16.x-DISPATCH-2099
PROTON_VERSION=0.34.0
ROUTER_SOURCE_URL=http://github.com/kgiusti/dispatch/archive/${DISPATCH_VERSION}.tar.gz
PROTON_SOURCE_URL=http://github.com/apache/qpid-proton/archive/${PROTON_VERSION}.tar.gz

# If a DOCKER_TAG is specified, go ahead and use it.
# if DOCKER_TAG is not specified use the DISPATCH_VERSION as the DOCKER_TAG
ifneq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG_VAL=$(DOCKER_TAG)
else
	DOCKER_TAG_VAL=$(DISPATCH_VERSION)
endif

# Ignores pushing latest tag when DISPATCH_VERSION contains freeze or x
PUSH_SKIP=false
ifeq ($(DOCKER_TAG_VAL), latest)
ifneq (,$(findstring freeze,$(DISPATCH_VERSION)))
	PUSH_SKIP=true
endif
ifneq (,$(findstring x,$(DISPATCH_VERSION)))
	PUSH_SKIP=true
endif
endif

all: build

build:
	${DOCKER} build -t qdrouterd-builder:${DOCKER_TAG_VAL} builder
	${DOCKER} run -ti -v $(PWD):/build:z -w /build qdrouterd-builder:${DOCKER_TAG_VAL} bash build_tarballs ${ROUTER_SOURCE_URL} ${PROTON_SOURCE_URL}

clean:
	rm -rf proton_build proton_install qpid-dispatch.tar.gz qpid-dispatch-src qpid-proton.tar.gz qpid-proton-src staging build

cleanimage:
	${DOCKER} image rm -f qdrouterd-builder

buildimage:
	${DOCKER} build -t ${PROJECT_NAME}:latest .

testimages:
	@echo Building atomic test images
	(cd test/container-images && make)

test:
	@echo Running qdr-image tests
	cd test/k8s ; go test -v -count 1 -p 1 -timeout 60m -tags integration ./integration/...

pushlocal:
	${DOCKER} tag ${PROJECT_NAME}:latest 127.0.0.1:5000/${PROJECT_NAME}:latest
	${DOCKER} push 127.0.0.1:5000/${PROJECT_NAME}:latest

push:
# DOCKER_USER and DOCKER_PASSWORD is useful in the CI environment.
# Use the DOCKER_USER and DOCKER_PASSWORD if available
# if not available, assume the user has already logged in
ifneq ($(strip $(DOCKER_USER)$(DOCKER_PASSWORD)),)
	@${DOCKER} login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
endif

# If DISPATCH_VERSION contains freeze and tag to be pushed is latest, it should be skipped
ifneq ($(PUSH_SKIP),true)
	${DOCKER} tag ${PROJECT_NAME}:latest ${DOCKER_REGISTRY}/${DOCKER_ORG}/${PROJECT_NAME}:${DOCKER_TAG_VAL}
	${DOCKER} push ${DOCKER_REGISTRY}/${DOCKER_ORG}/${PROJECT_NAME}:${DOCKER_TAG_VAL}
else
	@echo Push skipped for $(DOCKER_TAG_VAL)
endif

.PHONY: build buildimage cleanimage clean test push pushlocal
