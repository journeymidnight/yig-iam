.PHONY: build
GOPATH = $(PWD)/build
export GOPATH
GOBIN = $(PWD)/build/bin
export GOBIN
URL = github.com/journeymidnight
REPO = yig-iam
URLPATH = $(PWD)/build/src/$(URL)

build:
	@[ -d $(URLPATH) ] || mkdir -p $(URLPATH)
	@[ -d $(GOBIN) ] || mkdir -p $(GOBIN)
	@ln -nsf $(PWD) $(URLPATH)/$(REPO)
	go build $(URL)/$(REPO)
	go build -o yig-iam-tools $(URL)/$(REPO)/tools
	cp -f yig-iam $(PWD)/build/bin/
	cp -f yig-iam-tools $(PWD)/build/bin/

clean:
	rm -rf build
