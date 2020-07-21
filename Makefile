VERSION=

release: check-version check-master
	git tag v${VERSION} && \
	git push origin v${VERSION}

clear-dist:
	rm -rf ./dist

build-local: clear-dist
	goreleaser build --skip-validate

#
# Checking rules
# https://stackoverflow.com/a/4731504
#

check-version:
ifndef VERSION
	$(error VERSION is not set)
endif

check-master:
ifneq ($(shell git rev-parse --abbrev-ref HEAD),master)
	$(error Not on branch master)
endif
