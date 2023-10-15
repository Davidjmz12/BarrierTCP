package bar

import (
	"BarrierTCP/barrier/barMS"
	"BarrierTCP/barrier/com"
	"fmt"
	"github.com/DistributedClocks/GoVector/govec"
	"github.com/fatih/color"
	"log"
	"os"
	"time"
)

type BMessage struct {
	TypeMsg string
	Pid     int
}

type BarrierManager struct {
	barrier      *barMS.BarrierSys
	rtt          time.Duration
	logger       *govec.GoLog
	waitTo       chan bool
	responseChan chan bool
	size         int
}

func New(rtt time.Duration, size int, connFile string, me int, fileLog string, logger *govec.GoLog) (bm BarrierManager) {
	log.SetFlags(log.Lmicroseconds)
	if len(fileLog) > 0 {
		f, err := os.OpenFile(fileLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		com.CheckError(err)
		log.SetOutput(f)
	}
	waitChan := make(chan bool, size)
	responseChan := make(chan bool)
	barrierTypes := []barMS.Message{BMessage{}, barMS.End{}}
	bs := barMS.New(me, connFile, barrierTypes, logger)

	bm.barrier = &bs
	bm.rtt = rtt
	bm.logger = logger
	bm.size = size
	bm.waitTo = waitChan
	bm.responseChan = responseChan

	return
}

func (bm *BarrierManager) sendMessages(typeMsg string, received map[string][]bool) {
	me := bm.barrier.Me
	for i, value := range received[typeMsg] {
		j := i + 1
		if !value || typeMsg == "ack" {
			msg := BMessage{typeMsg, me}
			bm.barrier.Send(j, msg, false)
		}
	}
}

func (bm *BarrierManager) waitReplies(waitingWhat string, received map[string][]bool) {

	timeToResend := time.Tick(bm.rtt * time.Millisecond)
	end := false
	for !end {
		select {
		case <-timeToResend:
			com.PrintMsgLog(color.CyanString, "Resending Messages")
			bm.sendMessages(waitingWhat, received)
		case msg := <-bm.barrier.Mbox:
			var msgFormatted BMessage
			ok := com.TryDecode(msg, &msgFormatted)
			com.PrintMsgLog(color.GreenString, "Message Received:", msgFormatted)
			if !ok {
				fmt.Println("Not a good type.")
				os.Exit(1)
			}

			received[msgFormatted.TypeMsg][msgFormatted.Pid-1] = true
			if msgFormatted.TypeMsg == "ack" {
				received["ready"][msgFormatted.Pid-1] = true
			}
			end = com.AllTrue(received[waitingWhat])
		}
	}
}

func (bm *BarrierManager) Ready(mode bool) {
	bm.waitTo <- mode
	if mode {
		<-bm.responseChan
	}
}

func (bm *BarrierManager) waitSync() {

}

func (bm *BarrierManager) waitAll() {
	for i := 0; i < bm.size; i++ {
		<-bm.waitTo
	}
	close(bm.responseChan)
}

func (bm *BarrierManager) StartBarrier() {

	bm.logger.LogLocalEvent("Starting Barrier", govec.GetDefaultLogOptions())
	received := make(map[string][]bool)
	received["ready"] = make([]bool, bm.barrier.Size())
	received["ack"] = make([]bool, bm.barrier.Size())

	received["ready"][bm.barrier.Me-1] = true
	received["ack"][bm.barrier.Me-1] = true

	com.PrintMsgLog(color.YellowString, "Waiting for barrier limits to be done...")

	bm.waitAll()

	com.PrintMsgLog(color.HiMagentaString, "Starting Barrier System")

	com.PrintMsgLog(color.CyanString, "Sending messages")
	bm.sendMessages("ready", received)

	com.PrintMsgLog(color.YellowString, "Waiting replies...")
	bm.waitReplies("ready", received)
	com.PrintMsgLog(color.HiGreenString, "All replies received")

	com.PrintMsgLog(color.CyanString, "Sending the acks...")
	bm.sendMessages("ack", received)

	com.PrintMsgLog(color.YellowString, "Waiting the acks...")
	bm.waitReplies("ack", received)
	com.PrintMsgLog(color.HiGreenString, "All acks received...")

	bm.barrier.Send(bm.barrier.Me, barMS.End{End: true}, false)
	com.PrintMsgLog(color.HiMagentaString, "End Barrier System")

	bm.logger.LogLocalEvent("Stopping barrier", govec.GetDefaultLogOptions())
}
