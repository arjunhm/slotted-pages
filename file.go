package page

import (
	"bytes"
	"encoding/gob"
	"os"
)

type File struct {
	OSFile    *os.File
	Directory Directory
	Pages     []Page
}

func NewFile(fileName string) (*File, error) {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &File{
		OSFile:    file,
		Directory: *NewDirectory(),
		Pages:     make([]Page, DIR_ENTRY_LIMIT),
	}, nil

}

func (f *File) ReadFile() error {

	var buffer bytes.Buffer

	fileStats, err := f.OSFile.Stat()
	if err != nil {
		return err
	}

	// fileSize := fileStats.Size()
	_, err = f.OSFile.Read(buffer.Bytes())
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&f.Directory)
	if err != nil {
		return err
	}

	err = decoder.Decode(&f.Pages)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) WriteFile() error {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(f.Directory)
	if err != nil {
		return err
	}

	err = encoder.Encode(f.Pages)
	if err != nil {
		return err
	}

	_, err = f.OSFile.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
