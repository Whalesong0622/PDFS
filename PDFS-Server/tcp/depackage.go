package tcp

func depackage(byteStream []byte,pg *Package){
	pg.Op = string(byteStream[0])
	pg.FileName = string(byteStream[1:])
}