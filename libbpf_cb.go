package libbpfgo

import (
	"C"
	"sync/atomic"
	"time"
	"unsafe"
)
import "fmt"

var (
	waitingEvents atomic.Int64
	passedEvents  atomic.Int64
)

func init() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			fmt.Printf("libbpfgo: waitingEvents: %d\n", waitingEvents.Load())
			fmt.Printf("libbpfgo: passedEvents: %d\n", passedEvents.Load())
			passedEvents.Store(0)
		}
	}()
}

// revive:disable

// This callback definition needs to be in a different file from where it is declared in C
// Otherwise, multiple definition compilation error will occur

//export perfCallback
func perfCallback(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	pb := eventChannels.get(uint(uintptr(ctx))).(*PerfBuffer)
	waitingEvents.Add(1)
	pb.eventsChan <- C.GoBytes(data, size)
	passedEvents.Add(1)
	waitingEvents.Add(-1)
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
