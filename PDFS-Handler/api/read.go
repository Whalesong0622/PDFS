package api

import (
	"PDFS-Handler/common"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type block struct {
	file  []byte
	index int
}
type blocks []block
func (b blocks) Less(i, j int) bool { return b[i].index < b[j].index }
func (b blocks) Swap(i, j int) { b[i],b[j] = b[j],b[i] }
func (b blocks) Len() int { return len(b) }

func Read(fileName string, conn net.Conn, blockNums int) {
	defer conn.Close()
	blockPath = common.GetBlocksPathConfig()
	path := strings.Join([]string{blockPath, fileName}, "")

	buf := make([]byte, 1024*1024*64)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	begin := now.Local().UnixNano() / (1000 * 1000)

	files := make(blocks,0)
	wc := sync.WaitGroup{}

	for i := 0; i < blockNums; i++ {
		go func(idx int) {
			wc.Add(1)
			conn2, err := net.Dial("tcp", "43.132.181.175:11111")
			defer conn2.Close()
			if err != nil {
				fmt.Println("net.Dial err = ", err)
				return
			}
			conn2.Write([]byte(strconv.Itoa(2)))
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("conn.Read err", err)
			}
			if "ok" == string(buf[:n]) {
				fmt.Println("成功连接，请输入需要下载的文件的名字")
				name := path + "-" + strconv.Itoa(idx)
				conn2.Write([]byte(name))
				tmpBlock := ReadFromServer(name, conn)
				files = append(files, block{tmpBlock, idx})
			}
		}(i)
	}
	wc.Wait()
	sort.Sort(files)

	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	if "ok" == string(buf[:n]) {
		for i := 0;i < blockNums;i++{
			conn.Write(files[i].file)
		}
	}
	end := time.Now().Local().UnixNano() / (1000 * 1000)
	log.Printf("Read file %s to %s ended! Timecost: %d ms", fileName, conn.RemoteAddr().String(), end-begin)
}

func ReadFromServer(name string, conn net.Conn) []byte {
	defer conn.Close()

	// 拿到数据
	buf := make([]byte, 1024*1024*64)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	file := make([]byte, 0)
	if "ok" == string(buf[:n]) {
		conn.Write([]byte("ok"))
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("conn.Read err =", err)
				if err == io.EOF {
					log.Println("文件结束了", err)
				}
				break
			}
			if n == 0 {
				log.Println("文件结束了", err)
				break
			}
			file = append(file, buf[:n]...)
			// fmt.Println("sum:", sum)
		}
	}
	return file
}
