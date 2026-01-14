package wasmZip

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
)

func readFile() (*[]byte, error) {
	content, err := os.ReadFile("/workspaces/repos/wasm-zip/web/main.js")
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func TestWasmZip(t *testing.T) {
	content, err := readFile()
	index := NewZip()
	result := AddFile(index, "test", *content)
	if result != "success" {
		t.Error("failed!")
	}
	zippedFileContent, err := GetFile(index)
	if err != nil {
		t.Error("failed!")
	}
	if zippedFileContent == nil {
		t.Error("failed!")
	}
	snapshot := zippedFileContent.Bytes()
	bytesReader := bytes.NewReader(snapshot)
	reader, err := zip.NewReader(bytesReader, int64(len(snapshot)))
	if err != nil {
		os.WriteFile("test.zip", snapshot, 0666)
		t.Error("failed!")
	}
	open, err := reader.Open("test")
	if err != nil {
		t.Error("failed!")
	}
	defer open.Close()
	stat, err := open.Stat()
	if err != nil {
		t.Error("failed!")
	}
	decompSize := stat.Size()
	byteRes := make([]byte, decompSize)
	_, err = open.Read(byteRes)
	if err != nil {
		t.Error("failed!")
	}
	if len(byteRes) != len(*content) {
		t.Error("failed!")
	}
}
