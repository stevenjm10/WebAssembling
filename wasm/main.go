//go:build js && wasm

// Package main is compiled to WebAssembly and exposes Go functions to JavaScript.
package main

import (
	"strings"
	"syscall/js"
)

// add returns a + b.
func add(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return js.ValueOf("add expects 2 arguments")
	}
	return js.ValueOf(args[0].Float() + args[1].Float())
}

// fibonacci returns the n-th Fibonacci number (iterative).
func fibonacci(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return js.ValueOf("fibonacci expects 1 argument")
	}
	n := args[0].Int()
	if n < 0 || n > 90 {
		return js.ValueOf("n must be between 0 and 90")
	}
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return js.ValueOf(a)
}

// countPrimes counts the primes below limit (sieve of Eratosthenes).
// Used to compare WebAssembly and JavaScript on CPU-heavy work.
func countPrimes(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return js.ValueOf("countPrimes expects 1 argument")
	}
	limit := args[0].Int()
	if limit < 2 {
		return js.ValueOf(0)
	}
	composite := make([]bool, limit)
	count := 0
	for i := 2; i < limit; i++ {
		if !composite[i] {
			count++
			for j := i * i; j < limit; j += i {
				composite[j] = true
			}
		}
	}
	return js.ValueOf(count)
}

// reverse reverses a string, respecting Unicode characters.
func reverse(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return js.ValueOf("reverse expects 1 argument")
	}
	r := []rune(args[0].String())
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return js.ValueOf(string(r))
}

// greet manipulates the DOM directly from Go.
func greet(this js.Value, args []js.Value) any {
	name := "World"
	if len(args) == 1 && strings.TrimSpace(args[0].String()) != "" {
		name = strings.TrimSpace(args[0].String())
	}
	doc := js.Global().Get("document")
	out := doc.Call("getElementById", "greet-output")
	out.Set("textContent", "Hello, "+name+"! This text was written by Go running as WebAssembly.")
	return nil
}

func main() {
	wasm := js.Global().Get("Object").New()
	wasm.Set("add", js.FuncOf(add))
	wasm.Set("fibonacci", js.FuncOf(fibonacci))
	wasm.Set("countPrimes", js.FuncOf(countPrimes))
	wasm.Set("reverse", js.FuncOf(reverse))
	wasm.Set("greet", js.FuncOf(greet))
	js.Global().Set("goWasm", wasm)

	js.Global().Call("dispatchEvent", js.Global().Get("CustomEvent").New("wasm-ready"))

	// Keep the Go program alive so the exported functions stay callable.
	select {}
}
