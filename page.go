package page

import (
	"encoding/binary"
)

const (
	PAGE_SIZE                = 4096 // 4KB
	HEADER_SIZE       uint32 = 12   // 4 page ID, 4 offset, 4 record count
	SLOT_POINTER_SIZE        = 4
	SLOT_LENGTH_SIZE         = 4
	SLOT_SIZE                = SLOT_POINTER_SIZE + SLOT_LENGTH_SIZE
)

type Page struct {
	Data []byte
}

// wrapper, easier to type
func putuint32(buf []byte, val uint32) {
	binary.LittleEndian.PutUint32(buf, val)
}

func AllocPage(pageID uint32) *Page {
	page := Page{
		Data: make([]byte, PAGE_SIZE),
	}

	// Add header data
	putuint32(page.Data[0:4], pageID)
	putuint32(page.Data[4:8], uint32(PAGE_SIZE))
	putuint32(page.Data[8:12], uint32(0))

	return &page
}

//func (p *Page) AddData(key string, val string) {
//}
