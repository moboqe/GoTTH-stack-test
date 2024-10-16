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
	go build -o D:\VSCode\GoProjects\gotth-test\tmp\main.exe D:\VSCode\GoProjects\gotth-test\cmd\main.go && air

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.Environment=production" -o \bin\cmd\main.go

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: test
test:
	  go test -race -v -timeout 30s ./...

# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:9001" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o D:/VSCode/GoProjects/gotth-test/tmp/bin/main.exe" --build.full_bin "D:/VSCode/GoProjects/gotth-test/tmp/bin/main.exe" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/styles.css --minify --watch

# run esbuild to generate the index.js bundle in watch mode.

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

# start all 5 watch processes in parallel.
live: 
	make -j5 live/templ live/server live/tailwind  live/sync_assets
