package service

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
)

type PNGStore interface {
	Save(filename string, imageData bytes.Buffer) (string, error)
}

type DiskPNGStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*PNGInfo
}

type PNGInfo struct {
	Filename string
	Path     string
}

func NewPNGStore(imageFolder string) *DiskPNGStore {
	return &DiskPNGStore{
		imageFolder: imageFolder,
		images:      make(map[string]*PNGInfo),
	}
}

func (store *DiskPNGStore) Save(filename string, imageData bytes.Buffer) (string, error) {

	reader := bytes.NewReader(imageData.Bytes())
	gzipreader, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}

	decompressedPNG, err := io.ReadAll(gzipreader)
	if err != nil {
		return "", err
	}
	imagePath := fmt.Sprintf("%s/%s", store.imageFolder, filename)

	err = os.WriteFile(imagePath, decompressedPNG, os.ModeType)
	if err != nil {
		return "", fmt.Errorf("cannot write the image")
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[filename] = &PNGInfo{
		Filename: filename,
		Path:     imagePath,
	}
	return filename, nil
}
