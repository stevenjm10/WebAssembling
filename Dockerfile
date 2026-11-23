# ---------- Stage 1: build ----------
# Compiles the Go code to WebAssembly and builds the static file server.
FROM golang:1.24-alpine AS build
WORKDIR /src

COPY go.mod ./
COPY serve.go ./
COPY wasm ./wasm
COPY web ./web

# 1) Copy Go's JavaScript glue file next to the page
# 2) Compile wasm/main.go -> web/main.wasm
# 3) Compile the server as a static Linux binary
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" web/ \
 && GOOS=js GOARCH=wasm go build -o web/main.wasm ./wasm \
 && CGO_ENABLED=0 go build -o /server serve.go

# ---------- Stage 2: run ----------
# Minimal image: only the server binary and the web files.
FROM scratch
WORKDIR /app
COPY --from=build /server /app/server
COPY --from=build /src/web /app/web
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/server"]
