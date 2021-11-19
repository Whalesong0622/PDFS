package tcp

import "PDFS-Handler/common"

var blockPath string
var localIp string

func init(){
	blockPath = common.GetBlocksPathConfig()
	localIp = common.GetBlocksPathConfig()
}