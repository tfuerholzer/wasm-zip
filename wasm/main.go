package main

import (
	"fmt"
	"syscall/js"

	"tfurholzer.github.io/wasm-zip/internal/wasmZip"
)

func JsNewZip(this js.Value, args []js.Value) any {
	zipId := wasmZip.NewZip()
	return js.ValueOf(zipId)
}

func JsAddFile(this js.Value, args []js.Value) any {
	index := args[0].Int()
	name := args[1].String()
	length := args[2].Get("length").Int()
	byteArr := make([]byte, length)
	js.CopyBytesToGo(byteArr, args[2])
	return wasmZip.AddFile(index, name, byteArr)
}

func JsGetFile(this js.Value, args []js.Value) any {
	index := args[0].Int()
	content, err := wasmZip.GetFile(index)
	if err != nil {
		return fmt.Sprint("failed to get file %d - %s", index, err.Error())
	}
	goBytes := content.Bytes()
	jsUnit8 := js.Global().Get("Uint8Array").New(len(goBytes))
	js.CopyBytesToJS(jsUnit8, goBytes)
	return jsUnit8
}

func main() {
	moduleDecl := js.Global().Get("Object").New()
	moduleDecl.Set("newZip", js.FuncOf(JsNewZip))
	moduleDecl.Set("addFile", js.FuncOf(JsAddFile))
	moduleDecl.Set("getFile", js.FuncOf(JsGetFile))
	js.Global().Set("wasmZipModule", moduleDecl)
	/* this select keeps the main process alive (without actually consuming cpu cycles). Without it the exported
	functions do not stay callable */
	select {}
}
