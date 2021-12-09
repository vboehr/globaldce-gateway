package wire
import (
	//"github.com/globaldce/globaldce-toolbox/applog"
	//"net"
	"time"
	//"encoding/json"
	"encoding/binary"
	"github.com/globaldce/globaldce-toolbox/utility"
)
type ConnectionPass struct{
	Version uint32
	Timestamp int64
	Bits uint32
	Nonce uint32
	Hash utility.Hash
}
func GenerateConnectionPass(bits uint32) (bool,ConnectionPass){
	//var bits uint32
	//bits=0xDAFFFFFF//TODO to be stored in usersettings
	targetbigint:=utility.BigIntFromCompact(bits)
	starttime:= time.Now().Unix()
	var pass ConnectionPass
	for starttime-time.Now().Unix()<60{//TODO optimize//
		for loopcounter:=0;loopcounter<10000;loopcounter++{//TODO optimize
			pass.Timestamp=int64 (time.Now().Unix())
			
			pass.Nonce++
			pass.ComputeHash()
	
			//targetbigint=utility.CorrectTargetBigInt(targetbigint,mb.Header.Timestamp,prevtime)
			if (utility.BigIntFromHash(&pass.Hash).Cmp(targetbigint)<0){
				return true,pass
			}
		}
	}

	return false,pass
}
func (cp *ConnectionPass) ComputeHash() {
	tmpbuffer := make([]byte, 0)
	bufferVersion := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferVersion, uint32(cp.Version))
	bufferTimestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufferTimestamp, uint64(cp.Timestamp))
	bufferBits := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferBits, uint32(cp.Bits))
	bufferNonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferNonce, uint32(cp.Nonce))

	tmpbuffer=append(tmpbuffer, bufferVersion ...)
	//tmpbuffer=append(tmpbuffer, mb.Header.Prevblockhash[:] ...)
	//tmpbuffer=append(tmpbuffer, mb.Header.Roothash[:] ...)
	tmpbuffer=append(tmpbuffer, bufferTimestamp ...)
	tmpbuffer=append(tmpbuffer, bufferBits ...)
	tmpbuffer=append(tmpbuffer, bufferNonce ...)
	cp.Hash=utility.ComputeHash(tmpbuffer)
}
/*
func (mb *Mainblock) Mine( prevtime int64, prevblockhash utility.Hash, bits uint32) bool {

	fmt.Println("Minining")
	mb.Header.Bits=bits
	targetbigint:=utility.BigIntFromCompact(bits)
	mb.Header.Prevblockhash=prevblockhash
	starttime:= time.Now().Unix()

	for starttime-time.Now().Unix()<60{//TODO optimize//
	for loopcounter:=0;loopcounter<10000;loopcounter++{//TODO optimize
 	

		mb.Header.Timestamp=int64 (time.Now().Unix())
		
		mb.Header.Nonce++
		mb.ComputeHash()

		//targetbigint=utility.CorrectTargetBigInt(targetbigint,mb.Header.Timestamp,prevtime)
		if (utility.BigIntFromHash(&mb.Header.Hash).Cmp(targetbigint)<0){
			return true
		}
	}
}
	return false
}



*/