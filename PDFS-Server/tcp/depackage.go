package tcp

func depackage(byteStream []byte, pg *Package) {
	pg.Op = byteStream[0]
	pg.FileName = string(byteStream[1:])
}
