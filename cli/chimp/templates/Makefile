templ:
	go tool templ generate --watch

# Run air to detect any go file changes to re-build and re-run the server. Use datastar for hotreload
server:
	env GOENV=dev go tool air \
	--build.cmd "go build -o tmp/bin/main ." \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "true" \
	--misc.clean_on_exit true

tailwind-clean:
	go tool gotailwind -i ./css/input.css -o ./static/style.css --clean

# Run tailwindcss to generate the styles.css bundle in watch mode.
tailwind-watch:
	go tool gotailwind -i ./css/input.css -o ./static/style.css --watch

# Start development server
dev:
	make tailwind-clean
	make -j3 tailwind-watch templ server
