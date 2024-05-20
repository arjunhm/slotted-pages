package page

import (
	"encoding/binary"
	"fmt"
)

const (
	PAGE_SIZE                = 4096 // 4KB
	HEADER_SIZE       uint32 = 12   // 4 page ID, 4 offset, 4 record count
	SLOT_POINTER_SIZE        = 4    // offset to where the data begins from
	SLOT_KEY_SIZE            = 4    // never used, but idk
	SLOT_VAL_SIZE            = 4
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

func (p *Page) SetFreeOffset(offset uint32) {
	putuint32(p.Data[4:8], offset)
}

func (p *Page) GetRecordCount() uint32 {
	return getuint32(p.Data[8:12])
}

func (p *Page) SetRecordCount(count uint32) {
	putuint32(p.Data[8:12], count)
}

func (p *Page) SetSlot(slotOffset, dataOffset, keyLength, valLength uint32) {

	offsetStart, offsetEnd := slotOffset, slotOffset+SLOT_POINTER_SIZE
	putuint32(p.Data[offsetStart:offsetEnd], uint32(keyLength))

	keyEnd := offsetEnd + SLOT_KEY_SIZE
	putuint32(p.Data[offsetEnd:keyEnd], uint32(keyLength))

	valEnd := keyEnd + SLOT_VAL_SIZE
	putuint32(p.Data[keyEnd:valEnd], uint32(valLength))
}

func (p *Page) GetSlot(slotOffset uint32) (uint32, uint32, uint32) {

	offsetStart, offsetEnd := slotOffset, slotOffset+SLOT_POINTER_SIZE
	dataOffset := getuint32(p.Data[offsetStart:offsetEnd])

	keyEnd := offsetEnd + SLOT_KEY_SIZE
	keyLength := getuint32(p.Data[offsetEnd:keyEnd])

	valEnd := keyEnd + SLOT_VAL_SIZE
	valLength := getuint32(p.Data[keyEnd:valEnd])

	return dataOffset, keyLength, valLength
}

func (p *Page) AddData(key string, val string) {

	offset := p.GetFreeOffset()
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
	p.SetFreeOffset(start)
	// add slot array
	p.SetSlot(slotOffset, start, keyLength, valLength)
	// increment count
	p.SetRecordCount(count + 1)
}

func (p *Page) DeleteData(slotIndex uint32) bool {
	// get count
	count := p.GetRecordCount()

	// check if slotIndex exists
	if slotIndex > count {
		return false
	}

	// get slot offset
	slotOffset := HEADER_SIZE + (slotIndex * SLOT_SIZE)
	offset, keyLength, valLength := GetSlot(slotIndex)
	length := keyLength + valLength

	// set data to 0
	for i := uint32(0); i < length; i++ {
		p.Data[offset+i] = 0
	}

	// update slot
	p.SetSlot(slotOffset, 0, 0)

	return true
}

func (p *Page) Vacuum() {}
