package common

const WRITE_OP = "1"
const READ_OP = "2"

func IsLegal(RemoteAddr string, user string, op string, path string, filename string) bool {
	if user == "" {

	} else if op == "" {

	} else if op == WRITE_OP {

	} else if op == READ_OP {

	} else {

	}
	return true
}
