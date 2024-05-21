package page

import (
	"os"
)

type DBFile struct {
	File *os.File
}

func CreateFile(fileName string) (*DBFile, error) {

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &DBFile{File: file}, nil
}


func (db *DBFile) ReadPage(pageID uint32) (*Page, error) {
	offset := int(pageID - 1) * PAGE_SIZE
	data := make([]byte, PAGE_SIZE)
	_, err := db.File.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}
	// look into this
	return &Page{Data: data}, nil
}

func (db *DBFile) WritePage(page *Page, pageID uint32) error {
	offset := int(pageID - 1) * PAGE_SIZE
	_, err := db.File.WriteAt(page.Data, offset)
	return err
}

func (db *DBFile) Close() error {
	return db.File.Close()
}



