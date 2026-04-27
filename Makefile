GO ?= GO111MODULE=on CGO_ENABLED=1 go

run:
	$(GO) run -v main.go

test:
	gotestsum --format-hide-empty-pkg -- ./... --race

update:
	git pull
	docker compose up -d --build
	docker compose logs -f backend
