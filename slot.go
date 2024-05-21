package page

import (
	"errors"
)

type Slot struct {
	Offset    uint32
	KeySize   uint32
	ValueSize uint32
	Deleted   bool
}

func NewSlot(offset, keySize, valSize uint32) *Slot {
	return &Slot{
		Offset:    offset,
		KeySize:   keySize,
		ValueSize: valSize,
		Deleted:   false,
	}
}

func (s *Slot) GetOffset() uint32 {
	return Slot.Offset
}

func (s *Slot) SetOffset(offset uint32) error {
	if offset < HEADER_SIZE || offset > PAGE_SIZE {
		return errors.New("Invalid offset")
	}
	s.Offset = offset
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

func (s *Slot) GetDeleted() bool {
	return s.Deleted
}

func (s *Slot) SetDeleted(del bool) {
	s.Deleted = del
}

func (s *Slot) SetSlot(offset, keySize, valSize uint32, del bool) {
	s.SetOffset(offset)
	s.SetKeySize(keySize)
	s.SetValueSize(valSize)
	s.SetDeleted(del)
}
