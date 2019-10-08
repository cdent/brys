
.PHONY: help gabbi cleangabbi

help:
	@echo "gabbi: run the gabbi tests"
	@echo "cleangabbi: remove gabbit virt env"

gabbi: brys
	rm .brys.pid || true
	brys & echo "$$!" > .brys.pid
	[ -d .gabbi ] || python3 -mvenv .gabbi
	[ -x .gabbi/bin/gabbi-run ] || .gabbi/bin/pip install gabbi git+https://github.com/cdent/gabbihtml.git#egg=gabbihtml
	.gabbi/bin/gabbi-run -r gabbihtml.handler:HTMLHandler http://localhost:3333 -- gabbits/*.yaml
	[ -e .brys.pid ] && kill -TERM $$(cat .brys.pid)
	rm .brys.pid || true


cleangabbi:
	rm -rf .gabbi

brys: ${GOBIN}/brys

${GOBIN}/brys: main.go
	go install
