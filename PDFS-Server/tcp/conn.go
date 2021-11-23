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
	user := string(buf[:n])
	log.Println("Receive request from:",user,",reply ok")
	_, _ = conn.Write([]byte("ok"))

	n, err = conn.Read(buf)
	if err != nil {
		log.Println("conn.Read err =", err)
		return
	}
	op := string(buf[:n])

	if op == WRITE_OP {
		log.Println("Receive write request from:",conn.RemoteAddr().String(),"Reply ok")
		_, _ = conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		path := string(buf[:n])
		log.Println("Receiving file", path, "from", conn.RemoteAddr().String(), ",reply ok")
		_, _ = conn.Write([]byte("ok"))

		api.RevFile(path, conn)
	} else if op == READ_OP {
		log.Println("Receive read request from:",conn.RemoteAddr().String(),"Reply ok")
		_, _ = conn.Write([]byte("ok"))

		n, err = conn.Read(buf)
		if err != nil {
			log.Println("conn.Read err =", err)
		}

		name := string(buf[:n])
		log.Println("Sending file", name, "from", conn.RemoteAddr(), ",reply ok")
		_, _ = conn.Write([]byte("ok"))

		api.SendFile(name, conn)
	} else {
		log.Println("Reply err to", conn.RemoteAddr().String())
		_, _ = conn.Write([]byte("error"))
		conn.Close()
	}
}
