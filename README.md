## Slotted Pages

Slotted pages are a memory management technique used in database systems to store variable-sized records efficiently within a fixed-size page. A slotted page structure typically includes:

Page Header: Contains metadata such as the number of records, free space offset, and a pointer to the slot directory.
Slot Directory: An array of pointers, each pointing to the location of a record within the page. Slots can be reused when records are deleted.
Data Area: The actual storage space for records, which can vary in size. Records are packed together to minimize fragmentation.

![slotted pages overview](img/slotted-pages2.png "Slotted Pages Overview")
