package tcp

import (
	"PDFS-Server/api"
	"log"
	"net"
)

const WRITE_OP = "1"
const READ_OP = "2"

func HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}

	op := string(buf[:n])

	if op == WRITE_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Receiving file", name, "from", conn.RemoteAddr().String(), ",reply ok")
		conn.Write([]byte("ok"))

		api.RevFile(name, conn)
	} else if op == READ_OP {
		log.Println("Reply ok to", conn.RemoteAddr().String())
		conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Sending file", name, "from", conn.RemoteAddr(), ",reply ok")
		conn.Write([]byte("ok"))

		api.SendFile(name, conn)
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		conn.Write([]byte("error"))
		conn.Close()
	}
}
