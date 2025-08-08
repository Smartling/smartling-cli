MAINTAINER = Alex Koval <akoval@smartling.com>
DESCRIPTION = CLI for Smartling Platform

all: clean get build
	@

build: darwin windows.exe linux
	@

get:
	go get

clean:
	rm -rf bin pkg
	mkdir bin

_PKG = pkg/build

_CONTROL = echo >> $(_PKG)/DEBIAN/control

deb: _pkg-init
	mkdir -p $(_PKG)/usr/bin
	cp bin/smartling.linux $(_PKG)/usr/bin/smartling
	mkdir -p $(_PKG)/DEBIAN
	$(_CONTROL) "Package: smartling"
	$(_CONTROL) "Version: $(VERSION)"
	$(_CONTROL) "Architecture: all"
	$(_CONTROL) "Section: unknown"
	$(_CONTROL) "Priority: extra"
	$(_CONTROL) "Maintainer: $(MAINTAINER)"
	$(_CONTROL) "Homepage: https://github.com/Smartling/smartling-cli"
	$(_CONTROL) "Description: $(DESCRIPTION)"
	dpkg -b $(_PKG) pkg/smartling-$(VERSION)_all.deb
	rm -rf $(_PKG)

_SPEC = echo >> $(_PKG)/smartling.spec

rpm: _pkg-init
	$(_SPEC) "Name: smartling"
	$(_SPEC) "Version: $(VERSION)"
	$(_SPEC) "Release: 1%{?dist}"
	$(_SPEC) "Summary: $(DESCRIPTION)"
	$(_SPEC) "License: MIT"
	$(_SPEC) "%description"
	$(_SPEC) "%install"
	$(_SPEC) "mkdir -p %{buildroot}/%{_bindir}"
	$(_SPEC) "cp $(PWD)/bin/smartling.linux %{buildroot}/%{_bindir}/smartling"
	$(_SPEC) "%files"
	$(_SPEC) "%{_bindir}/smartling"
	$(_SPEC) "%define _rpmdir $(_PKG)"
	rpmbuild -bb $(_PKG)/smartling.spec
	cp $(_PKG)/*/*.rpm pkg/
	rm -rf $(_PKG)

_pkg-init:
	rm -rf $(_PKG)
	mkdir -p $(_PKG)
	$(eval VERSION ?= \
		$(shell git rev-list --count HEAD).$(shell git rev-parse --short HEAD))

%:
	GOOS=$(basename $@) go build -o bin/smartling.$@

tidy:
	go mod tidy

_linter:
	go install github.com/mgechev/revive@v1.10.0

lint:
	golangci-lint run ./... && \
	revive -config .revive.toml ./...

_mockery-install:
	go install github.com/vektra/mockery/v3@v3.3.4

mockery:
	mockery --config .mockery.yml

test_unit:
	go test ./cmd/...

# add binary and config to tests/cmd/bin/ before run test integration
test_integration:
	go test ./tests/cmd/files/push/...
	go test ./tests/cmd/files/pull/...
	go test ./tests/cmd/files/list/...
	go test ./tests/cmd/files/status/...
	go test ./tests/cmd/files/rename/...
	go test ./tests/cmd/files/delete/...
	go test ./tests/cmd/projects/...
	go test ./tests/cmd/init/...
