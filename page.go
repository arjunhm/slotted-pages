package page

import "errors"

const (
	PAGE_SIZE   = 4096 // 4KB
	HEADER_SIZE = 12   // PageID+Count+FreeSpaceEnd
	SLOT_SIZE   = 12   // Offset+KeySize+ValueSize
)

type Page struct {
	Header Header
	Slots  []Slot
	Data   []byte
}

func NewPage(pageID uint32) *Page {
	h := NewHeader(pageID)
	return &Page{
		Header: *h,
		Slots:  make([]Slot, 0),
		Data:   make([]byte, PAGE_SIZE-HEADER_SIZE),
	}
}

func GetSlotOffset(count uint32) uint32 {
	return HEADER_SIZE + (count * SLOT_SIZE)
}

func (p *Page) AddRow(kv KeyValue) error {

	offset := p.Header.GetFreeSpaceEnd()

	dataSize := kv.GetSize()
	keySize := kv.GetKeySize()
	valSize := kv.GetValueSize()

	start := offset - dataSize
	keyEnd := start + keySize
	copy(p.Data[start:keyEnd], kv.Key)
	copy(p.Data[keyEnd:keyEnd+valSize], kv.Value)

	// create slot and add to slot array
	slot := NewSlot(offset, keySize, valSize)
	p.Slots = append(p.Slots, slot)

	// update free space ptr
	p.Header.SetFreeSpaceEnd(start)

	// increment count
	p.Header.SetCount(p.Header.GetCount() + 1)

	return nil
}

func (p *Page) UpdateRow(kv KeyValue) error {
	/*
		have to look into this
		basically, you create a new tuple and
		add it to the page, mark the old one as deleted
	*/
	return nil
}

func (p *Page) DeleteRow(slotIndex uint32) error {

	count := p.Header.GetCount()

	if slotIndex > count {
		return errors.New("Invalid slot index")
	}

	slotOffset := GetSlotOffset(count)
	slot := p.Slots[slotOffset]

	offset := slot.GetOffset()
	dataSize := slot.GetSize()

	for i := uint32(0); i < dataSize; i++ {
		p.Data[offset+i] = 0
	}

	// use this as the main thing later.
	// batch delete later
	p.Slots[slotIndex].SetSlot(0, 0, 0, true)

	return nil
}
