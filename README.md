# WebAssembly with Go – Demo

A minimal example that shows how to write functions in **Go**, compile them to **WebAssembly** (`.wasm`) and call them from **JavaScript** in the browser.

## What is WebAssembly?

WebAssembly (Wasm) is a compact binary instruction format that runs in all modern browsers next to JavaScript, at near-native speed and inside the same security sandbox. Languages such as C/C++, Rust and Go can be compiled to Wasm, which lets you reuse existing code or run CPU-heavy work (image processing, games, simulations, cryptography) on the web.

## What this demo shows

| # | Feature | Go function |
|---|---------|-------------|
| 1 | Pass numbers from JS to Go and get a result back | `add(a, b)` |
| 2 | Run an algorithm in Wasm | `fibonacci(n)` |
| 3 | Pass strings in both directions (Unicode-safe) | `reverse(text)` |
| 4 | Manipulate the DOM directly from Go | `greet(name)` |
| 5 | Compare WebAssembly vs JavaScript on the same algorithm | `countPrimes(limit)` |

## How it works

```
 wasm/main.go ──(GOOS=js GOARCH=wasm go build)──► web/main.wasm
                                                      │
 web/index.html ─ loads ─► wasm_exec.js (Go runtime glue)
                 └─ app.js ─ WebAssembly.instantiateStreaming(fetch("main.wasm"))
                             go.run(instance)  →  Go registers window.goWasm
                             goWasm.add(7, 35) →  42
```

1. `wasm/main.go` uses the `syscall/js` package to wrap Go functions with `js.FuncOf` and expose them on a global object `window.goWasm`.
2. When Go is ready it fires a `wasm-ready` event, and `select {}` keeps the program alive so the functions stay callable.
3. `wasm_exec.js` (shipped with Go) provides the `Go` class that connects the Wasm module to the browser.
4. `app.js` downloads and instantiates `main.wasm`, then calls `goWasm.*` from button handlers.

## Project structure

```
WebAssembling/
├── wasm/main.go     # Go code compiled to WebAssembly
├── web/
│   ├── index.html   # Demo page
│   ├── app.js       # Loads the .wasm and calls the Go functions
│   └── style.css
├── serve.go         # Tiny static server (serves .wasm as application/wasm)
├── Makefile
└── go.mod
```

## Requirements

- Go 1.24 or newer (`go version`)
- A modern browser (Chrome, Edge, Firefox, Safari)

## Run it

With `make`:

```bash
git clone https://github.com/stevenjm10/WebAssembling.git
cd WebAssembling
make serve
```

Without `make` (Linux/macOS):

```bash
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/
GOOS=js GOARCH=wasm go build -o web/main.wasm ./wasm
go run serve.go
```

Windows (PowerShell):

```powershell
Copy-Item "$(go env GOROOT)\lib\wasm\wasm_exec.js" web\
$env:GOOS="js"; $env:GOARCH="wasm"; go build -o web\main.wasm .\wasm
Remove-Item Env:GOOS, Env:GOARCH
go run serve.go
```

Then open **http://localhost:8080**.

> The page must be served over HTTP. Opening `index.html` directly (`file://`) fails because browsers block `fetch()` of local files and require the `application/wasm` content type for `instantiateStreaming`.

## About the benchmark

Example run in headless Chrome (limit = 5,000,000): Go/Wasm 79.8 ms vs JavaScript 54.5 ms, both finding 348,513 primes. WebAssembly is **not automatically faster** than JavaScript: modern JS engines are very well optimized for simple loops, and Go's Wasm output includes its runtime and garbage collector (`main.wasm` is about 1.7 MB). Wasm shines when you reuse existing native code, need predictable performance, or use languages with a smaller runtime (Rust, C, TinyGo). Try different limits and browsers and compare.

## Going further

- Build a much smaller module with [TinyGo](https://tinygo.org/): `tinygo build -o web/main.wasm -target wasm ./wasm` (use TinyGo's own `wasm_exec.js`).
- Run Wasm outside the browser with WASI: `GOOS=wasip1 GOARCH=wasm go build`.
- Docs: [Go WebAssembly wiki](https://go.dev/wiki/WebAssembly), [MDN – WebAssembly](https://developer.mozilla.org/en-US/docs/WebAssembly).

## Author

Steven Morales
