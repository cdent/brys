
.PHONY: help gabbi cleangabbi

help:
	@echo "gabbi: run the gabbi tests"
	@echo "cleangabbi: remove gabbit virt env"

gabbi:
	[ -d .gabbi ] || python3 -mvenv .gabbi
	[ -x .gabbi/bin/gabbi-run ] || .gabbi/bin/pip install gabbi git+https://github.com/cdent/gabbihtml.git#egg=gabbihtml
	.gabbi/bin/gabbi-run -r gabbihtml.handler:HTMLHandler -- gabbits/*.yaml

cleangabbi:
	rm -rf .gabbi
