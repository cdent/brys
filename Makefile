
.PHONY: help gabbi cleangabbi deps gotest

help:
	@echo "gabbi: run the gabbi tests"
	@echo "cleangabbi: remove gabbit virt env"

gabbi: brys
	packr2 clean
	rm -r store || true
	mkdir store
	[ -e .brys.pid ] && kill -TERM $$(cat .brys.pid) || true
	rm .brys.pid || true
	brys & echo "$$!" > .brys.pid
	[ -d .gabbi ] || python3 -mvenv .gabbi
	[ -x .gabbi/bin/gabbi-run ] || .gabbi/bin/pip install gabbi git+https://github.com/cdent/gabbihtml.git#egg=gabbihtml
	sleep 2 && .gabbi/bin/gabbi-run -r gabbihtml.handler:HTMLHandler http://localhost:3333 -- gabbits/*.yaml
	[ -e .brys.pid ] && kill -TERM $$(cat .brys.pid) || true
	rm .brys.pid || true

test: gotest gabbi

gotest: deps
	go test

cleangabbi:
	rm -rf .gabbi

deps:
	go get ./...

brys: deps ${GOBIN}/brys

${GOBIN}/brys: main.go
	go install
