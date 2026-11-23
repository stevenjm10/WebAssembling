# WebAssembly with Go – Demo

A minimal example that shows how to write functions in **Go**, compile them to **WebAssembly** (`.wasm`) and call them from **JavaScript** in the browser. The whole project is built and served with **Docker**.

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

Everything is built and run with **Docker**, so you don't need Go installed.

```
 Dockerfile
 ├─ Stage 1 (golang:1.24-alpine)
 │    wasm/main.go ──(GOOS=js GOARCH=wasm go build)──► web/main.wasm
 │    $(go env GOROOT)/lib/wasm/wasm_exec.js ────────► web/wasm_exec.js
 │    serve.go ──(CGO_ENABLED=0 go build)────────────► /server
 └─ Stage 2 (scratch, ~5 MB)
      /app/server + /app/web  ──►  http://localhost:8080

 Browser: index.html ─ loads ─► wasm_exec.js (Go runtime glue)
          app.js ─ WebAssembly.instantiateStreaming(fetch("main.wasm"))
                   go.run(instance)  →  Go registers window.goWasm
                   goWasm.add(7, 35) →  42
```

1. `wasm/main.go` uses the `syscall/js` package to wrap Go functions with `js.FuncOf` and expose them on a global object `window.goWasm`.
2. When Go is ready it fires a `wasm-ready` event, and `select {}` keeps the program alive so the functions stay callable.
3. `wasm_exec.js` (shipped with Go) provides the `Go` class that connects the Wasm module to the browser.
4. `app.js` downloads and instantiates `main.wasm`, then calls `goWasm.*` from button handlers.
5. The **multi-stage Dockerfile** compiles everything in the first stage and copies only the server binary and the web files into an empty `scratch` image, so the final image is tiny and has no shell or extra packages.
6. `serve.go` serves `.wasm` files with the `application/wasm` content type, which browsers require.

## Project structure

```
WebAssembling/
├── Dockerfile       # Multi-stage build: Go -> WebAssembly + static server
├── compose.yaml     # Docker Compose service (port 8080)
├── .dockerignore
├── Makefile         # Shortcuts for the Docker commands
├── wasm/main.go     # Go code compiled to WebAssembly
├── web/
│   ├── index.html   # Demo page
│   ├── app.js       # Loads the .wasm and calls the Go functions
│   └── style.css
├── serve.go         # Tiny static server used inside the container
└── go.mod
```

## Requirements

- [Docker](https://docs.docker.com/get-docker/) (Docker Desktop on Windows/macOS, Docker Engine on Linux)
- A modern browser (Chrome, Edge, Firefox, Safari)

## Run it with Docker

```bash
git clone https://github.com/stevenjm10/WebAssembling.git
cd WebAssembling
```

**Option A – Docker Compose (recommended):**

```bash
docker compose up -d --build   # build the image and start the container
docker compose logs            # see the server output
docker compose down            # stop and remove the container
```

**Option B – plain Docker commands:**

```bash
docker build -t webassembling:latest .
docker run -d --rm --name webassembling -p 8080:8080 webassembling:latest
docker stop webassembling
```

**Option C – Makefile shortcuts:** `make up` / `make down`, or `make docker-run` / `make docker-stop`.

Then open **http://localhost:8080**.

Useful checks:

```bash
docker ps                      # the container should be "Up" with 0.0.0.0:8080->8080/tcp
docker images webassembling    # final image is about 5 MB of content
```

To use another host port, change the mapping, e.g. `-p 3000:8080`, and open http://localhost:3000.

## About the benchmark

Example run in headless Chrome against the Docker container (limit = 5,000,000): Go/Wasm 74.5 ms vs JavaScript 58.2 ms, both finding 348,513 primes. WebAssembly is **not automatically faster** than JavaScript: modern JS engines are very well optimized for simple loops, and Go's Wasm output includes its runtime and garbage collector (`main.wasm` is about 1.7 MB). Wasm shines when you reuse existing native code, need predictable performance, or use languages with a smaller runtime (Rust, C, TinyGo). Try different limits and browsers and compare.

## Going further

- Build a much smaller module with [TinyGo](https://tinygo.org/) by using a TinyGo base image in the first stage of the Dockerfile (`tinygo build -o web/main.wasm -target wasm ./wasm`, with TinyGo's own `wasm_exec.js`).
- Run Wasm outside the browser with WASI: `GOOS=wasip1 GOARCH=wasm go build`.
- Docs: [Go WebAssembly wiki](https://go.dev/wiki/WebAssembly), [MDN – WebAssembly](https://developer.mozilla.org/en-US/docs/WebAssembly).

## Author

Steven Morales
