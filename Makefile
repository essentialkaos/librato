########################################################################################

# This Makefile generated by GoMakeGen 0.4.0 using next command:
# gomakegen --metalinter .

########################################################################################

.PHONY = fmt deps

########################################################################################

deps:
	git config --global http.https://pkg.re.followRedirects true
	go get -v pkg.re/essentialkaos/ek.v7

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

metalinter:
	test -s $(GOPATH)/bin/gometalinter || (go get -u github.com/alecthomas/gometalinter ; $(GOPATH)/bin/gometalinter --install)
	$(GOPATH)/bin/gometalinter --deadline 30s

########################################################################################
