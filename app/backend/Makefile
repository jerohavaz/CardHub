.PHONY: generate lint fmt fmt-check deps test up down migrate logs smoke clean

MODULES := foundation protocol service-example gateway

generate:
	cd protocol && buf generate

lint:
	cd protocol && buf lint

fmt:
	@set -e; for mod in $(MODULES); do \
		echo "==> $$mod"; \
		(cd $$mod && go fmt ./...); \
	done

fmt-check:
	@set -e; \
	unformatted=$$(find $(MODULES) -type f -name '*.go' -exec gofmt -l {} +); \
	if [ -n "$$unformatted" ]; then \
		printf 'Unformatted Go files:\n%s\n' "$$unformatted"; \
		exit 1; \
	fi

deps:
	@set -e; for mod in $(MODULES); do \
		echo "==> $$mod"; \
		(cd $$mod && go mod tidy); \
	done

test:
	@set -e; for mod in $(MODULES); do \
		echo "==> $$mod"; \
		(cd $$mod && go test ./...); \
	done

up:
	docker compose up --build -d

migrate:
	docker compose run --build --rm migrate

down:
	docker compose down

logs:
	docker compose logs -f gateway service-example migrate

smoke:
	sh scripts/smoke.sh

clean:
	docker compose down -v