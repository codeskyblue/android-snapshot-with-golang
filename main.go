package main

// #include "screen.h"
import "C"
import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

const DEV_FB0 = "/dev/graphics/fb0"

var screenOnce = sync.Once{}

// return: width, height, bytes_per_pixel
func ScreenInfo() (w, h, bytespp, offset int) {
	screenOnce.Do(func() {
		C.init()
	})
	w = int(C.width())
	h = int(C.height())
	bytespp = int(C.bytespp())
	offset = int(C.offset())
	return
}

func Snapshot() (img *image.RGBA, err error) {
	fd, err := os.OpenFile(DEV_FB0, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	fmt.Println(fd.Stat())

	w, h, bpp, offset := ScreenInfo()
	fmt.Println(ScreenInfo())
	size := int(w * h * bpp)
	mapsize := int(offset) + size
	fmt.Println(size, offset, int(fd.Fd()))

	mmap, err := syscall.Mmap(int(fd.Fd()), 0, mapsize, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		fmt.Println(err, len(mmap))
		// invalid argument
		return nil, err
	}
	fmt.Println("OK")
	defer syscall.Munmap(mmap)

	img = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(w), int(h)}})
	piex := (*[1e9]uint8)(unsafe.Pointer(&mmap[int(offset)]))

	for i := 0; i < size; i++ {
		img.Pix[i] = piex[i]
		if i < 3 {
			fmt.Println("PIEX:", piex[i])
		}
	}
	return img, nil
}

func main() {
	img, err := Snapshot()
	if err != nil {
		log.Fatal(err)
	}
	fd, err := os.Create("/data/local/tmp/s.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	if err = png.Encode(fd, img); err != nil {
		log.Fatal(err)
	}

}
