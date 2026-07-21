.PHONY: tidy swagger build test frontend-build

GOPROXY ?= https://goproxy.cn,direct
SERVICES := gateway system infra pay member business

tidy:
	GOPROXY=$(GOPROXY) go mod tidy

swagger:
	GOPROXY=$(GOPROXY) go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/system/main.go -o docs --parseDependency

build: swagger
	mkdir -p bin
	$(foreach service,$(SERVICES),GOPROXY=$(GOPROXY) go build -trimpath -o bin/nimbus-$(service) ./cmd/$(service);)

test:
	GOPROXY=$(GOPROXY) go test ./...

frontend-build:
	cd ../frontend && pnpm install --frozen-lockfile && pnpm ts:check && pnpm build:prod
