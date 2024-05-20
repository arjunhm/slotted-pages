package page

import (
	"encoding/binary"
)

/*
Directory contains
- Record count - 4 bytes
- Slots <PageId, FreeSpaceAvailable> - 8 bytes
*/

const (
	DIR_PAGE_ID_SIZE    = 4
	DIR_FREE_SPACE_SIZE = 4
	DIR_SLOT_SIZE       = DIR_PAGE_ID_SIZE + DIR_FREE_SPACE_SIZE
)

type Directory struct {
	Data []byte
}

func NewDirectory() *Directory {
	d := &Directory{
		Data: make([]byte, PAGE_SIZE),
	}
}

func (d *Directory) GetCount() uint32 {
	return getuint32(d.Data[0:HEADER_SIZE])
}

func (d *Directory) SetCount(count uint32) {
	putuint32(d.Data[0:HEADER_SIZE], count)
}

func (d *Directory) SetSlot(slotOffset, pageID, freeSpace uint32) {
	pageIDEnd := slotOffset + DIR_PAGE_ID_SIZE
	putuint32(d.Data[slotOffset:pageIDEnd], pageID)
	putuint32(d.Data[pageEnd:DIR_SLOT_SIZE], freeSpace)
}

func (d *Directory) GetSlot(slotOffset) (uint32, uint32) {
	pageIDEnd := slotOffset + DIR_PAGE_ID_SIZE
	pageID := getuint32(d.Data[slotOffset:pageIDEnd])
	freeSpace := getuint32(d.Data[pageEnd:DIR_SLOT_SIZE])
	return pageID, freeSpace
}

func (d *Directory) AddData(pageID, freeSpace uint32) {

	count := d.GetCount()
	offset := HEADER_SIZE + (count * DIR_SLOT_SIZE)

	d.SetSlot(offset, pageID, freeSpace)
	d.SetCount(count + 1)
}

func (d *Directory) UpdateData(pageID, freeSpace uint32) {

	count := d.GetCount()
	if pageID > count {
		return
	}

	offset := HEADER_SIZE + ((pageID - 1) * DIR_SLOT_SIZE)
	d.SetSlot(offset, pageID, freeSpace)
}

func (d *Directory) RemoveData(pageID uint32) {}
