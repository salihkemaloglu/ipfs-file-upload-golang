package ipfs

import (	
	"bytes"
	shell "github.com/ipfs/go-ipfs-api"

)

func UploadToIpfs(data []byte) (string,error){
	sh := shell.NewShell("localhost:5001")
	reader := bytes.NewReader(data)
	fileHash, err := sh.Add(reader)
	if err != nil {
		return "",err
	}
	return fileHash,nil
}