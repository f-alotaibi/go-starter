server:
	air \
	--build.cmd "go build -o tmp/bin/main ./server.go" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

tailwind:
	tailwindcss -i ./assets/tailwind.css -o ./assets/public/main.css --watch --poll

templ:
	templ generate --watch

dev:
	make -j3 tailwind templ server