VERSION = $(shell grep -E 'VERSION =' server.go | sed -e 's/^.*VERSION =[^0-9]*/v/g')
PROG = skrassiev/webbench

all:
	@echo Building $(VERSION)
	@echo docker build --tag $(PROG):$(VERSION) .
	@echo docker tag $(PROG):$(VERSION) $(DOCKER_REGISTRY)/$(PROG):$(VERSION)
	@echo docker push $(DOCKER_REGISTRY)/$(PROG):$(VERSION)
	
