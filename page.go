package page

import "errors"

const (
	PAGE_SIZE = 4096 // 4KB
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

// ----- Helper functions -----

func (p *Page) GetFreeSpaceStart() uint32 {
	count := p.Header.GetCount()
	lastSlotOffset := HEADER_SIZE + (count * SLOT_SIZE)
	return lastSlotOffset + SLOT_SIZE
}

func (p *Page) GetAvailableSpace() uint32 {
	return p.Header.GetFreeSpaceEnd() - p.GetFreeSpaceStart()
}

func (p *Page) GetSlot(slotID uint32) *Slot {
	// slotCount := p.Header.GetCount()

	// naive search for now cause why not
	// pls make it binary search later 😭
	for i := 0; i < len(p.Slots); i++ {
		if p.Slots[i].SlotID == slotID {
			return &p.Slots[i]
		}
	}
	return nil
}

// ----- Main functions -----

func (p *Page) Insert(kv KeyValue) error {
	// ? what if key already exists ?

	payloadSize := kv.GetSize()
	keySize := kv.GetKeySize()
	valSize := kv.GetValueSize()

	// check if space exists
	if payloadSize > p.GetAvailableSpace() {
		return errors.New("Insufficient space")
	}

	// get offset from where data can be inserted
	start := p.Header.GetFreeSpaceEnd() - payloadSize

	// add data
	copy(p.Data[start:start+keySize], kv.Key)
	copy(p.Data[start+keySize:start+payloadSize], kv.Value)

	// add slot
	slotID := p.Header.GetSlotID() + uint32(1) // idk if the uint32 part is required
	slot := NewSlot(slotID, start, keySize, valSize)
	p.Slots = append(p.Slots, slot)

	// update Header: count, freeSpaceEnd, slotID
	p.Header.SetHeader(p.Header.GetCount()+1, start, slotID)

	return nil
}

func (p *Page) Update(slotID uint32, kv KeyValue) error {

	slot := p.GetSlot(slotID)
	if slot == nil {
		return errors.New("Invalid Slot ID")
	}

	oldSize := slot.GetSize()
	newSize := kv.GetSize()

	// if same size
	if oldSize == newSize {
		offset := slot.GetOffset()
		keyEnd := offset + slot.GetKeySize()
		valEnd := keyEnd + slot.GetValueSize()

		// insert data
		copy(p.Data[offset:keyEnd], kv.Key)
		copy(p.Data[keyEnd:valEnd], kv.Value)

		// update key and val size in slot
		keySize := kv.GetKeySize()
		valSize := kv.GetValueSize()
		slot.SetSlot(offset, keySize, valSize)
	} else {
		// if freeSpace enough
		if newSize < p.GetAvailableSpace(){
			// delete slot
			err := p.Delete(slotID)
			if err != nil {
				return err
			}
			// insert data
			err = p.Insert(kv)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Insufficient space")
		}
	}

	return nil
}

func (p *Page) Delete(slotID uint32) error {

	// get slot
	slot := p.GetSlot(slotID)
	if slot == nil {
		return errors.New("Invalid Slot ID")
	}
	payloadOffset := slot.GetOffset()

	// get payload size
	payloadSize := slot.GetSize()
	freeSpaceEnd := p.Header.GetFreeSpaceEnd()

	// if not last slot
	if freeSpaceEnd != payloadOffset {
		// move data
		src := freeSpaceEnd
		dest := src + payloadSize
		nBytes := src + payloadOffset
		copy(p.Data[dest:dest+nBytes], p.Data[src:src+nBytes])
	}

	// update freeSpaceEnd
	p.Header.SetFreeSpaceEnd(freeSpaceEnd + payloadSize)

	// update slots with new offset
	var slotIndex int
	for i := 0; i < len(p.Slots); i++ {
		s := p.Slots[i]
		if s.GetOffset() < payloadOffset {
			s.SetOffset(s.GetOffset() + payloadSize)
		}
		if s.GetSlotID() == slotID {
			slotIndex = i
		}
	}

	// delete slot
	p.Slots = append(p.Slots[:slotIndex], p.Slots[slotIndex+1:]...)
	// update count
	p.Header.SetCount(p.Header.GetCount() - 1)
	return nil
}

