package service

import (
	"bytes"
	"fmt"
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
	imagePath := fmt.Sprintf("%s/%s", store.imageFolder, filename)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("Cannot create the image\n")
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Cannot write the image\n")
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[filename] = &PNGInfo{
		Filename: filename,
		Path:     imagePath,
	}
	return filename, nil
}
