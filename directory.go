package page

import (
	"errors"
)

/*
Directory contains
- Record count - 4 bytes
- Slots <PageId, FreeSpaceAvailable> - 8 bytes
*/

const (
	DIR_COUNT_SIZE   = 4                                                  // # of entries
	DIR_PAGE_ID_SIZE = 4                                                  // bytes for pageID
	DIR_OFFSET_SIZE  = 4                                                  // bytes for offset
	DIR_FREE_SIZE    = 4                                                  // bytes for free space
	DIR_SLOT_SIZE    = DIR_PAGE_ID_SIZE + DIR_OFFSET_SIZE + DIR_FREE_SIZE // 12
	DIR_ENTRY_LIMIT  = 10                                                 // can store 10 entries
	DIR_ENTRY_SIZE   = DIR_SLOT_SIZE * DIR_ENTRY_LIMIT                    // 10 * 12 bytes

)

type Directory struct {
	Count uint32
	Entry []Entry
}

type Entry struct {
	PageID    uint32
	Offset    uint32
	FreeSpace uint32
}

func NewDirectory() *Directory {
	return &Directory{
		Count: 0,
		Entry: make([]Entry, DIR_ENTRY_SIZE),
	}
}

func (d *Directory) UpdateCount() {
	d.Count = uint32(len(d.Entry))
}

// might get rid of this
func (d *Directory) SetCount(count uint32) {
	d.Count = count
}

func (d *Directory) GetCount() uint32 {
	return d.Count
}

func (d *Directory) CreateEntry(pageID, offset, freeSpace uint32) error {

	entryPtr, err := NewEntry(pageID, offset, freeSpace)
	if err != nil {
		return err
	}
	entry := *entryPtr

	d.Entry = append(d.Entry, entry)
	d.UpdateCount()

	return nil
}

func (d *Directory) ReadEntry(pageID uint32) (*Entry, error) {

	if pageID >= d.GetCount() {
		return nil, errors.New("pageID does not exist")
	}

	return &d.Entry[pageID], nil
}

func (d *Directory) UpdateEntry(pageID, offset, freeSpace uint32) error {

	if pageID >= d.GetCount() {
		return errors.New("pageID does not exist")
	}
	d.Entry[pageID].Offset = offset
	d.Entry[pageID].FreeSpace = freeSpace

	return nil
}

// idk about this one for now
func (d *Directory) DeleteEntry() {}

// ------ ENTRY ------

func NewEntry(pageID, offset, freeSpace uint32) (*Entry, error) {

	if pageID >= DIR_ENTRY_LIMIT {
		return nil, errors.New("Entry limit exceeded")
	}

	return &Entry{
		PageID:    pageID,
		Offset:    offset,
		FreeSpace: freeSpace,
	}, nil
}

func (e *Entry) Decode() (uint32, uint32, uint32) {
	return e.PageID, e.Offset, e.FreeSpace
}
