all: test vet fmt lint build

test:
	cp -r fixtures/merge_into fixtures/merge_into_tmp
	go test -race -covermode=atomic -tags librrd ./...
	go test -race -covermode=atomic -tags rrdtool ./...

vet:
	go vet ./...

fmt:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l
	test -z $$(go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 gofmt -l)

lint:
	go list -f '{{.Dir}}' ./... | grep -v /vendor/ | xargs -L1 revive | grep -v '/rrd.go'

build:
	go build -tags rrdtool -o bin/rrdmerge ./cmd/rrdmerge
	go build -tags librrd -o bin/rrdmerge_librrd ./cmd/rrdmerge
	strip -s bin/rrdmerge
	strip -s bin/rrdmerge_librrd
