package page

import (
	"os"
)

type File struct {
	OSFile    *os.File
	Directory Directory
	Pages     []Page
}

func NewFile(fileName string) *File {

	// create file in disk
	// store file in File.OSFile
	// create directory NewDirectory()
	// create pages

}
