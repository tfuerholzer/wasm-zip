package main

import (
	"archive/zip"
	"bytes"
	"math/rand"
	"syscall/js"
)

var memZipMap = make(map[int]*MemoryZip)

func nextEmpty() int {
	for {
		newIndex := rand.Int()
		if _, ok := memZipMap[newIndex]; !ok {
			return newIndex
		}
	}
}

type MemoryZip struct {
	ZipWriter *zip.Writer
	Buffer    *bytes.Buffer
}

func Initialize(this js.Value, args []js.Value) any {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	mapIndex := nextEmpty()
	memZipData := MemoryZip{
		ZipWriter: writer,
		Buffer:    buf,
	}
	memZipMap[mapIndex] = &memZipData
	return mapIndex
}

func AddFile(this js.Value, args []js.Value) any {
	id := args[0].Int()
	name := args[1].String()
	content := args[2]

	memZipData, found := memZipMap[id]
	if !found {
		return "Zip file not initialized!"
	}
	goContent := make([]byte, content.Length())
	js.CopyBytesToGo(goContent, content)
	writer, err := memZipData.ZipWriter.Create(name)
	if err != nil {
		return err.Error()
	}
	_, err = writer.Write(goContent)
	if err != nil {
		return err.Error()
	}
	return "success"
}

func main() {
    // this shit does the same as Window.initialize = . So it polutes the Global namespace
	js.Global().Set("initialize", js.FuncOf(Initialize))
	js.Global().Set("addFile", js.FuncOf(AddFile))
	/* this select keeps the main process alive (without actually consuming cpu cycles). Without it the exported 
	functions do not stay callable */
	select {}
}
