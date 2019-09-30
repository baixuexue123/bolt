package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type pgid uint64

type page struct {
	id       pgid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}

var pageSize int

func init() {
	pageSize = os.Getpagesize()
}

func pageInBuffer(b []byte, id pgid) *page {
	return (*page)(unsafe.Pointer(&b[id*pgid(pageSize)]))
}

func readPage() {
	f, err := os.Open("my1.db")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var buf [0x1000 * 4]byte
	if _, err := f.ReadAt(buf[:], 0); err == nil {
		for i := 0; i < 4; i++ {
			page := pageInBuffer(buf[:], pgid(i))
			fmt.Println(*page)
		}
	}
}

// fdatasync flushes written data to a file descriptor.
func fdatasync(file *os.File) error {
	return syscall.Fdatasync(int(file.Fd()))
}

func writePage() {
	// Create two pages on a buffer.
	buf := make([]byte, pageSize*4)
	for i := 0; i < 4; i++ {
		p := pageInBuffer(buf[:], pgid(i))
		p.id = pgid(i)
		p.flags = uint16(i * i)
	}

	f, err := os.OpenFile("my1.db", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Write the buffer to data file.
	if _, err := f.WriteAt(buf, 0); err != nil {
		log.Fatal(err)
	}
	if err := fdatasync(f); err != nil {
		log.Fatal(err)
	}
}

func main() {

	fmt.Printf("pageSize: %d\n", pageSize)

	writePage()
	readPage()

}
