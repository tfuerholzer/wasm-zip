#!/usr/bin/env bash
env GOOS=js GOARCH=wasm go build -o web/main.wasm -ldflags="-s -w" $(pwd) 
cp $(go env GOROOT)/lib/wasm/wasm_exec.js ./web
B64_CONTENT=$(base64 -w 0 web/main.wasm)
JS_CONTENT=$(printf "const go = new Go()\nconst wasmBytes=Buffer.from('%s','base64')\nconst { instance } = await WebAssembly.instantiate(wasmBytes, go.importObject)\ngo.run(instance);\nexport const mainModule = instance" "$B64_CONTENT")

printf "/*GENERATED FROM wasm-zip repo main.go file using build.sh*/\n\n%s\n\n%s" "$(cat ./web/wasm_exec.js)" "$JS_CONTENT" > ./web/main.js
