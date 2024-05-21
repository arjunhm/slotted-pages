package page

import (
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
		Directory: NewDirectory,
		Pages:     make([]Pages, DIR_ENTRY_LIMIT),
	}, nil

}
