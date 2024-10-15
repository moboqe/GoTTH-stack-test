.PHONY: echo
echo:
	@echo $(PATH)

#NPM version
.PHONY: tailwind-install
tailwind-install:
	npm install -D tailwindcss

.PHONY: tailwind-init
tailwind-init:
	npx tailwindcss init

.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

#standalone CLI(needs to manually download executable tailwind)
.PHONY: tailwind-init-s
tailwind-init-s:
	./tailwindcss init

.PHONY: tailwind-watch-s
tailwind-watch-s:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build-s
tailwind-build-s:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: dev
dev:
	go build -o ./tmp/main.exe ./cmd/main.go
	air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o ./bin/$(APP_NAME)  ./cmd/main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test:
	  go test -race -v -timeout 30s ./...