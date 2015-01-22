VERSION = $(shell grep -E 'VERSION =' server.go | sed -e 's/^.*VERSION =[^0-9]*/v/g')
PROG = skrassiev/webbench

all:
	@echo Building $(VERSION)
	@docker build --tag $(PROG):$(VERSION) .
	@docker tag $(PROG):$(VERSION) $(DOCKER_REGISTRY)/$(PROG):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(PROG):$(VERSION)
	
