.PHONY: build clean

clean: ## remove all build releated artifacts
	-rm -rf build

build: clean ## build all binaries
	mkdir -p build
	go build -o build/git-gh-clone gh-clone/main.go
	go build -o build/git-gh-foreach-repository gh-foreach-repository/main.go
	go build -o build/git-gh-gists gh-gists/main.go

install-linux: build ## install gh tool chain on linux
	cp build/git-gh-* ${HOME}/.local/bin/
	wget -O build/subcall.tgz https://github.com/sascha-andres/subcall/releases/download/v1.0.0/subcall_1.0.0_linux_amd64.tar.gz
	cd build && tar xzf subcall.tgz
	chmod u+x build/subcall
	cp build/subcall ${HOME}/.local/bin/git-gh
	cp build/subcall ${HOME}/.local/bin/git-gh-foreach

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
