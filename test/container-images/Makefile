DOCKER=docker
DOCKER_REGISTRY=docker.io
DOCKER_ORG=atomictests

# $(1) - image name
# $(2) - image tag
# $(3) - Dockerfile directory
define build_push_image
	@echo Building atomic test image: $(1)
	$(DOCKER) build -t $(1):$(2) $(3)
	$(DOCKER) tag $(1):$(2) $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(1):$(2)
	$(DOCKER) tag $(1):$(2) $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(1):latest
	$(if $(DOCKER_USER),\
		$(DOCKER) push $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(1):$(2) \
	,)
	$(if $(DOCKER_USER),\
		$(DOCKER) push $(DOCKER_REGISTRY)/$(DOCKER_ORG)/$(1):latest \
	,)
endef

all: build-push-images

build-push-images:
ifneq ($(strip $(DOCKER_USER)),)
	@$(DOCKER) login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
endif
	$(call build_push_image,jms-amqp-tests,1.1,java/jms-amqp-tests)
