package barMS

import (
	"BarrierTCP/barrier/com"
	"encoding/gob"
	"github.com/DistributedClocks/GoVector/govec"
	"github.com/fatih/color"
	"net"
)

type Message interface{}

type End struct {
	End bool
}

type BarrierSys struct {
	Mbox   chan Message
	peers  []string
	Me     int
	logger *govec.GoLog
}

const (
	MAXMESSAGES = 10000
)

// Pre: pid en {1..n}, el conjunto de procesos del SD
// Post: envía el mensaje msg a pid
func (bs *BarrierSys) Send(pid int, msg Message, checkErr bool) {
	conn, err := net.Dial("tcp", bs.peers[pid-1])
	if checkErr || (!checkErr && err == nil) {
		com.CheckError(err)
		encoder := gob.NewEncoder(conn)
		outBuf := bs.logger.PrepareSend("Sending barrier message", msg, govec.GetDefaultLogOptions())
		err = encoder.Encode(&outBuf)
		com.CheckError(err)
		conn.Close()
	} else {
		com.PrintMsgLog(color.RedString, "Message discarded(not sent)", msg, "to", pid)
	}
}

// Pre: True
// Post: el mensaje msg de algún Proceso P_j se retira del mailbox y se devuelve
//
//	Si mailbox vacío, receive bloquea hasta que llegue algún mensaje
func (bs *BarrierSys) receive() (msg Message) {
	msg = <-bs.Mbox
	return msg
}

//	messageTypes es un slice con tipos de mensajes que los procesos se pueden intercambiar a través de este ms
//
// Hay que registrar un mensaje antes de poder utilizar (enviar o recibir)
// Notar que se utiliza en la función New
func register(messageTypes []Message) {
	for _, msgTp := range messageTypes {
		gob.Register(msgTp)
	}
}

// Pre: whoIam es el pid del proceso que inicializa este ms
//
//	usersFile es la ruta a un fichero de texto que en cada línea contiene IP:puerto de cada participante
//	messageTypes es un slice con todos los tipos de mensajes que los procesos se pueden intercambiar a través de este ms
func New(whoIam int, usersFile string, messageTypes []Message, logger *govec.GoLog) (bs BarrierSys) {
	bs.logger = logger
	bs.Me = whoIam
	bs.peers = com.ParsePeers(usersFile)
	bs.Mbox = make(chan Message, MAXMESSAGES)
	register(messageTypes)
	go func() {
		listener, err := net.Listen("tcp", bs.peers[bs.Me-1])
		com.CheckError(err)
		com.PrintMsgLog(color.HiWhiteString, "Barrier listening at "+bs.peers[bs.Me-1])
		defer close(bs.Mbox)
		for {
			conn, err := listener.Accept()
			com.CheckError(err)
			decoder := gob.NewDecoder(conn)
			var msg Message
			var inBuf []byte
			err = decoder.Decode(&inBuf)
			bs.logger.UnpackReceive("Received barrier message", inBuf, &msg, govec.GetDefaultLogOptions())
			com.CheckError(err)
			conn.Close()
			var endMsg End
			isEnd := com.TryDecode(msg, &endMsg)
			if !isEnd {
				com.CheckError(err)
				bs.Mbox <- msg
			} else {
				com.PrintMsgLog(color.HiWhiteString, "Listening Process End")
				return
			}
		}
	}()
	return bs
}

func (bs *BarrierSys) Size() int {
	return len(bs.peers)
}
