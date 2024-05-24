package page

import (
	"fmt"
	"testing"
)

const TEST_PAGE_ID = 0

func compareInt(t *testing.T, expected, got uint32, name string) {
	if expected != got {
		t.Errorf("%s Expected: %d, got: %d", name, expected, got)
	}
}

func compareKV(t *testing.T, expKV, gotKV *KeyValue) {
	if expKV.Key != gotKV.Key {
		t.Errorf("Expected %s, got %s", expKV.Key, gotKV.Key)
	}
	if expKV.Value != gotKV.Value {
		t.Errorf("Expected %s, got %s", expKV.Value, gotKV.Value)
	}
}

func compareHeader(t *testing.T, header Header, pageID, count, freeSpaceEnd, slotID uint32) {
	compareInt(t, pageID, header.GetPageID(), "Header.PageID")
	compareInt(t, count, header.GetCount(), "Header.Count")
	compareInt(t, freeSpaceEnd, header.GetFreeSpaceEnd(), "Header.FreeSpaceEnd")
	compareInt(t, slotID, header.GetSlotID(), "Header.LatestSlotID")
}

func compareSlot(t *testing.T, slot *Slot, slotID, offset, keySize, valSize uint32) {
	compareInt(t, slotID, slot.GetSlotID(), "Slot.slotID")
	compareInt(t, offset, slot.GetOffset(), "Slot.offset")
	compareInt(t, keySize, slot.GetKeySize(), "Slot.keySize")
	compareInt(t, valSize, slot.GetValueSize(), "Slot.valSize")
}

func createKV(id int) *KeyValue {
	switch id {
	case 1:
		return NewKeyValue("name", "john")
	case 2:
		return NewKeyValue("pincode", "12345")
	default:
		return NewKeyValue("key", "value")
	}
}

func createPage(t *testing.T, pageID uint32, kvIDLimit int) *Page {
	var kv *KeyValue
	var size uint32
	page := NewPage(pageID)

	for i := 1; i <= kvIDLimit; i++ {
		kv = createKV(i)
		size += kv.GetSize()
		page.Insert(kv)

		count := uint32(i)
		freeSpaceEnd := (PAGE_SIZE - HEADER_SIZE) - size
		slotID := uint32(i)
		compareHeader(t, page.Header, pageID, count, freeSpaceEnd, slotID)

		slot := page.GetSlot(slotID)
		compareSlot(t, slot, slotID, freeSpaceEnd, kv.GetKeySize(), kv.GetValueSize())
	}
	return page
}

func TestNewPage(t *testing.T) {
	fmt.Print("starting tests...\n")
	pageID := uint32(12)
	page := NewPage(pageID)

	gotPageID := page.Header.GetPageID()
	compareInt(t, pageID, gotPageID, "pageID")

	fsEnd := PAGE_SIZE - HEADER_SIZE
	gotFSEnd := page.Header.GetFreeSpaceEnd()
	compareInt(t, fsEnd, gotFSEnd, "fsEnd")
}

func TestInsert(t *testing.T) {
	pageID := uint32(1)
	createPage(t, pageID, 3)
}

func DeleteSlot(t *testing.T, page *Page, pageID, slotID, lastSlotID uint32) {

	slot := page.GetSlot(slotID)
	payloadSize := slot.GetSize() // 12

	// gather expected data
	freeSpaceEnd := page.Header.GetFreeSpaceEnd() + payloadSize
	count := page.Header.GetCount() - uint32(1)
	nextSlot := page.GetSlot(slotID + uint32(1))
	var nextSlotOffset uint32
	if nextSlot != nil {
		nextSlotOffset = nextSlot.GetOffset() + payloadSize
	}

	// delete slot
	page.Delete(slotID)

	if nextSlot != nil {
		// slot i + 1 sould be moved to left
		if *nextSlot != page.Slots[slotID-1] {
			t.Errorf("nextSlot not in expected position")
		}
		// slot i+1 offset should be updated
		compareInt(t, nextSlotOffset, nextSlot.GetOffset(), "offset")
	}

	// slot should not exist
	errSlot := page.GetSlot(slotID)
	if errSlot != nil {
		t.Errorf("Slot %d should have been deleted", slotID)
	}

	// header count, freeSpaceEnd,
	compareHeader(t, page.Header, pageID, count, freeSpaceEnd, lastSlotID)
}

func TestDelete(t *testing.T) {
	// create page with 3 items
	pageID := uint32(1)
	page := createPage(t, pageID, 3)

	slotID := uint32(2)
	DeleteSlot(t, page, pageID, slotID, 3)

	// delete slot 3 (start)
	DeleteSlot(t, page, pageID, uint32(3), 3)

	// delete slot 1 (end)
	DeleteSlot(t, page, pageID, uint32(1), 3)

}

func TestRead(t *testing.T) {
	var expKV, gotKV *KeyValue

	pageID := uint32(1)
	page := createPage(t, pageID, 3)

	expKV = createKV(1)
	gotKV = page.Read(1)
	if expKV.Key != gotKV.Key {
		t.Errorf("Expected %s, got %s", expKV.Key, gotKV.Key)
	}
	if expKV.Value != gotKV.Value {
		t.Errorf("Expected %s, got %s", expKV.Value, gotKV.Value)
	}

	expKV = createKV(2)
	gotKV = page.Read(2)
	compareKV(t, expKV, gotKV)

	expKV = createKV(3)
	gotKV = page.Read(3)
	compareKV(t, expKV, gotKV)

	expKV = createKV(3)
	gotKV = page.Read(2)
	if expKV.Key == gotKV.Key {
		t.Errorf("Expected %s, got %s", expKV.Key, gotKV.Key)
	}
	if expKV.Value == gotKV.Value {
		t.Errorf("Expected %s, got %s", expKV.Value, gotKV.Value)
	}
}

func TestUpdate(t *testing.T) {
	var slotID uint32
	var expKV *KeyValue

	pageID := uint32(1)
	page := createPage(t, pageID, 3)

	// same size
	slotID = uint32(3)
	newKV := NewKeyValue("key", "xalue")
	newKeySize := newKV.GetKeySize()
	newValueSize := newKV.GetValueSize()

	err := page.Update(slotID, newKV)
	if err != nil {
		expKV = page.Read(slotID)
		compareKV(t, expKV, newKV)
	}

	// diff size
	slotID = uint32(2)
	oldSlot := page.GetSlot(slotID)
	newOffset := page.Header.GetFreeSpaceEnd() + oldSlot.GetOffset() - newKV.GetKeySize()
	err = page.Update(slotID, newKV)

	if err != nil {
		newSlotID := uint32(4)
		expKV = page.Read(newSlotID)
		compareKV(t, expKV, newKV)

		slot := page.GetSlot(newSlotID)
		compareSlot(t, slot, newOffset, newKeySize, newValueSize, newSlotID)

	}

	// insufficient space
	slotID = uint32(1)
	value := string(make([]byte, 4096))
	newKV = NewKeyValue("myKey", value)
	err = page.Update(slotID, newKV)
	if err == nil {
		t.Errorf("insufficient space error should be raised")
	}
}
