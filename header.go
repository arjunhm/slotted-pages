package page

type Header struct {
	PageID uint32
	Count uint32
	FreeSpaceEnd uint32
}

func NewHeader(pageID uint32) *Header {
	return &Header{
		PageID: pageID,
		Count: 0,
		FreeSpaceEnd: PAGE_SIZE-HEADER_SIZE
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

func (h *Header) GetFreeSpace() uint32 {
	return h.FreeSpaceEnd
}

func (h *Header) SetFreeSpaceEnd(freeSpaceEnd uint32) {
	h.FreeSpaceEnd = freeSpaceEnd
}

