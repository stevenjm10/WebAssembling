// 1. Load and start the Go WebAssembly module.
const go = new Go(); // provided by wasm_exec.js

window.addEventListener("wasm-ready", () => {
  const status = document.getElementById("status");
  status.textContent = "WebAssembly module loaded ✔";
  status.className = "status ready";
  document.querySelectorAll("button").forEach((b) => (b.disabled = false));
});

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((result) => go.run(result.instance))
  .catch((err) => {
    const status = document.getElementById("status");
    status.textContent = "Error loading WebAssembly: " + err;
    status.className = "status error";
  });

// 2. Call the functions exported by Go (window.goWasm).
const $ = (id) => document.getElementById(id);

$("btn-add").onclick = () => {
  const r = goWasm.add(Number($("a").value), Number($("b").value));
  $("add-output").textContent = `Result: ${r}`;
};

$("btn-fib").onclick = () => {
  const n = Number($("fib-n").value);
  $("fib-output").textContent = `fibonacci(${n}) = ${goWasm.fibonacci(n)}`;
};

$("btn-reverse").onclick = () => {
  $("reverse-output").textContent = goWasm.reverse($("text").value);
};

$("btn-greet").onclick = () => goWasm.greet($("name").value);

// Same sieve algorithm in JavaScript, for comparison.
function countPrimesJS(limit) {
  if (limit < 2) return 0;
  const composite = new Uint8Array(limit);
  let count = 0;
  for (let i = 2; i < limit; i++) {
    if (!composite[i]) {
      count++;
      for (let j = i * i; j < limit; j += i) composite[j] = 1;
    }
  }
  return count;
}

$("btn-bench").onclick = () => {
  const limit = Number($("limit").value);
  let t0 = performance.now();
  const wasmCount = goWasm.countPrimes(limit);
  const wasmMs = performance.now() - t0;

  t0 = performance.now();
  const jsCount = countPrimesJS(limit);
  const jsMs = performance.now() - t0;

  $("bench-output").textContent =
    `WebAssembly (Go): ${wasmCount} primes in ${wasmMs.toFixed(1)} ms | ` +
    `JavaScript: ${jsCount} primes in ${jsMs.toFixed(1)} ms`;
};
