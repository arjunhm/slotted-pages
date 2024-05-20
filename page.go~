package page

import (
	"encoding/binary"
	"fmt"
)

const (
	PAGE_SIZE                = 4096 // 4KB
	HEADER_SIZE       uint32 = 12   // 4 page ID, 4 offset, 4 record count
	SLOT_POINTER_SIZE        = 4    // offset to where the data begins from
	SLOT_LENGTH_SIZE         = 4    // never used, but idk
	SLOT_SIZE                = SLOT_POINTER_SIZE + SLOT_LENGTH_SIZE
)

// wrapper, easier to type
func putuint32(buf []byte, val uint32) {
	binary.LittleEndian.PutUint32(buf, val)
}

func getuint32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

type Page struct {
	Data []byte
}

func (p *Page) DisplayHeaderDetails() {
	fmt.Printf("PageID: %d\n", p.GetPageID())
	fmt.Printf("Offset: %d\n", p.GetFreeOffset())
	fmt.Printf("Count: %d\n", p.GetRecordCount())
	fmt.Print("-----\n")
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

func (p *Page) GetPageID() uint32 {
	return getuint32(p.Data[0:4])
}

func (p *Page) GetFreeOffset() uint32 {
	return getuint32(p.Data[4:8])
}

func (p *Page) GetRecordCount() uint32 {
	return getuint32(p.Data[8:12])
}

func (p *Page) SetRecordCount(count uint32) {
	putuint32(p.Data[8:12], count)
}

func (p *Page) UpdateSlot(slotOffset, dataOffset, dataLength uint32) {
	putuint32(p.Data[slotOffset:slotOffset+SLOT_POINTER_SIZE], uint32(dataOffset))
	putuint32(p.Data[slotOffset+SLOT_POINTER_SIZE:slotOffset+SLOT_SIZE], uint32(dataLength))
}

func (p *Page) AddData(key string, val string) {
	// get free ptr offset
	offset := p.GetFreeOffset()

	// get record count
	count := p.GetRecordCount()

	// get slot
	slotOffset := HEADER_SIZE + (count * SLOT_SIZE)

	// calculate length of data
	keyLength := uint32(len(key))
	valLength := uint32(len(val))
	kvLength := keyLength + valLength
	start := offset - kvLength

	// add data to end of page
	copy(p.Data[start:start+keyLength], key)
	copy(p.Data[start+keyLength:offset], val)

	// update free space pointer
	putuint32(p.Data[4:8], uint32(start))

	// add slot array
	p.UpdateSlot(slotOffset, start, kvLength)

	// increment count
	p.SetRecordCount(count + 1)
}

func (p *Page) DeleteData(slotIndex uint32) bool {
	// get count
	count := getuint32(p.Data[8:12])

	// check if slotIndex exists
	if slotIndex > count {
		return false
	}

	// get slot offset
	slotOffset := HEADER_SIZE + (slotIndex * SLOT_SIZE)
	// get pointer to data
	offset := getuint32(p.Data[slotOffset : slotOffset+SLOT_POINTER_SIZE])
	// get length of data
	length := getuint32(p.Data[slotOffset+SLOT_POINTER_SIZE : slotOffset+SLOT_SIZE])

	// set data to 0
	for i := uint32(0); i < length; i++ {
		p.Data[offset+i] = 0
	}

	// update slot
	p.UpdateSlot(slotOffset, 0, 0)

	return true
}

func (p *Page) Vacuum() {}
