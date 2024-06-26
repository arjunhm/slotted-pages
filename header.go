package page

const HEADER_SIZE = 16 // <PageID,Count,FreeSpaceEnd,SlotID>

type Header struct {
	PageID       uint32
	Count        uint32
	FreeSpaceEnd uint32
	SlotID       uint32
}

func NewHeader(pageID uint32) *Header {
	return &Header{
		PageID:       pageID,
		Count:        0,
		FreeSpaceEnd: PAGE_SIZE - HEADER_SIZE,
		SlotID:       0,
	}
}

func (h *Header) GetPageID() uint32 {
	return h.PageID
}

func (h *Header) GetCount() uint32 {
	return h.Count
}

func (h *Header) SetCount(c uint32) {
	h.Count = c
}

func (h *Header) GetFreeSpaceEnd() uint32 {
	return h.FreeSpaceEnd
}

func (h *Header) SetFreeSpaceEnd(freeSpaceEnd uint32) {
	h.FreeSpaceEnd = freeSpaceEnd
}

func (h *Header) GetSlotID() uint32 {
	return h.SlotID
}

func (h *Header) SetSlotID(id uint32) {
	h.SlotID = id
}

func (h *Header) SetHeader(count, freeSpaceEnd, slotID uint32) {
	h.SetCount(count)
	h.SetFreeSpaceEnd(freeSpaceEnd)
	h.SetSlotID(slotID)
}
