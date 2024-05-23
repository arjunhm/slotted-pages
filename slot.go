package page

import (
	"errors"
)

const SLOT_SIZE = 16 // <SlotID,Offset,KeySize,ValueSize>

type Slot struct {
	SlotID    uint32
	Offset    uint32
	KeySize   uint32
	ValueSize uint32
}

func NewSlot(slotID, offset, keySize, valSize uint32) Slot {
	return Slot{
		SlotID:    slotID,
		Offset:    offset,
		KeySize:   keySize,
		ValueSize: valSize,
	}
}

func (s *Slot) GetSlotID() uint32 {
	return s.SlotID
}

func (s *Slot) GetOffset() uint32 {
	return s.Offset
}

func (s *Slot) SetOffset(offset uint32) error {
	if offset < HEADER_SIZE || offset > PAGE_SIZE {
		return errors.New("Invalid offset")
	}
	s.Offset = offset
	return nil
}

func (s *Slot) GetKeySize() uint32 {
	return s.KeySize
}

func (s *Slot) SetKeySize(keySize uint32) {
	s.KeySize = keySize
}

func (s *Slot) GetValueSize() uint32 {
	return s.ValueSize
}

func (s *Slot) SetValueSize(valSize uint32) {
	s.ValueSize = valSize
}

func (s *Slot) GetSize() uint32 {
	return s.GetKeySize() + s.GetValueSize()
}

func (s *Slot) SetSlot(offset, keySize, valSize uint32) {
	s.SetOffset(offset)
	s.SetKeySize(keySize)
	s.SetValueSize(valSize)
}
