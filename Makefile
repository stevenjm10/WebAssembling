.PHONY: build serve clean

build:
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./wasm

serve: build
	go run serve.go

clean:
	rm -f web/main.wasm web/wasm_exec.js
