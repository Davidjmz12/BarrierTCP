package bar

import (
	"fmt"
	"github.com/DistributedClocks/GoVector/govec"
	"testing"
	"time"
)

func simulateClient(bm *BarrierManager, t time.Duration) {
	go waitSendThree(bm, t)
	bm.StartBarrier()
}

func waitSendThree(bm *BarrierManager, t time.Duration) {
	for i := 1; i <= 3; i++ {
		time.Sleep(1000 * t * time.Millisecond)
		bm.Ready(false)
	}
	fmt.Println("I ACCESS")
}

func waitBlock(bm *BarrierManager, t time.Duration) {

	go func() {
		time.Sleep(1000 * t * time.Millisecond)
		bm.Ready(true)
		fmt.Println("I ACCESS")
	}()

	bm.StartBarrier()
}

func TestBlockingBarrier(t *testing.T) {
	logger1 := govec.InitGoVector("Process_1", "log/LogFile_1", govec.GetDefaultConfig())
	logger2 := govec.InitGoVector("Process_2", "log/LogFile_2", govec.GetDefaultConfig())
	bm1 := New(3000, 3, "./barrier.txt", 1, "", logger1)
	bm2 := New(3000, 1, "./barrier.txt", 2, "", logger2)

	go simulateClient(&bm1, 4)

	waitBlock(&bm2, 1)

}

func TestBarrier(t *testing.T) {

	logger1 := govec.InitGoVector("Process_1", "log/LogFile_1", govec.GetDefaultConfig())
	logger2 := govec.InitGoVector("Process_2", "log/LogFile_2", govec.GetDefaultConfig())
	bm1 := New(1000, 3, "./barrier2.txt", 1, "", logger1)
	bm2 := New(1000, 3, "./barrier2.txt", 2, "", logger2)

	go simulateClient(&bm1, 1)

	simulateClient(&bm2, 3)
}
