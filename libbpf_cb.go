package libbpfgo

import (
	"C"
	"unsafe"
)
import "fmt"

// revive:disable

// This callback definition needs to be in a different file from where it is declared in C
// Otherwise, multiple definition compilation error will occur

//export perfCallback
func perfCallback(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	pb := eventChannels.get(uint(uintptr(ctx))).(*PerfBuffer)
	pb.eventsChan <- C.GoBytes(data, size)
}

//export perfLostCallback
func perfLostCallback(ctx unsafe.Pointer, cpu C.int, cnt C.ulonglong) {
	pb := eventChannels.get(uint(uintptr(ctx))).(*PerfBuffer)
	if pb == nil {
		fmt.Println("pb is nil")
		return
	}
	if pb.lostChan != nil {
		fmt.Printf("Received lost event: %d\n", cnt)
		pb.lostChan <- uint64(cnt)
	} else {
		fmt.Println("pb.lostChan is nil")
	}
}

//export ringbufferCallback
func ringbufferCallback(ctx unsafe.Pointer, data unsafe.Pointer, size C.int) C.int {
	ch := eventChannels.get(uint(uintptr(ctx))).(chan []byte)
	ch <- C.GoBytes(data, size)

	return C.int(0)
}

// revive:enable
