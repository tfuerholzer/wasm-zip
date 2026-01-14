package wasmZip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"maps"
	"strings"
)

type MemoryZip struct {
	ZipWriter *zip.Writer
	Buffer    *bytes.Buffer
}

var memZipMap = make(map[int]*MemoryZip)

func nextEmpty() int {
	for i := 0; i < 50000; i++ {
		if _, exist := memZipMap[i]; !exist {
			return i
		}
	}
	return -1
}

func NewZip() int {
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

func AddFile(index int, name string, content []byte) string {
	zipItem, ok := memZipMap[index]
	if !ok {
		var strs []string
		for key := range maps.Keys(memZipMap) {
			strs = append(strs, fmt.Sprintf("%d", key))
		}
		return fmt.Sprintf("not found, valid keys are %s", strings.Join(strs, ";"))
	}
	created, err := zipItem.ZipWriter.Create(name)
	if err != nil {
		return fmt.Sprintf("file %s could not be created because %s", name, err.Error())
	}
	_, err = created.Write(content)
	if err != nil {
		return fmt.Sprintf("file %s be written not be created because %s", name, err.Error())
	}
	return "success"
}

func GetFile(index int) (*bytes.Buffer, error) {
	zipItem, ok := memZipMap[index]
	if !ok {
		return nil, fmt.Errorf("could not find zipItem %d!", index)
	}
	delete(memZipMap, index)
	err := zipItem.ZipWriter.Close()
	if err != nil {
		return nil, err
	}
	return zipItem.Buffer, nil
}
