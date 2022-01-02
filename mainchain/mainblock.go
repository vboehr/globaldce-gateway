package mainchain
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	"fmt"
	"time"
	"encoding/json"
	"encoding/binary"
	"github.com/globaldce/globaldce-toolbox/utility"
)
type Mainblock struct {
	// 
	Height uint32
	Header Mainheader
	Transactions [] utility.Transaction
}

type Mainheader struct {
	// 
	Version int32
	Prevblockhash utility.Hash
	Roothash utility.Hash//
	Timestamp int64
	Bits uint32
	Nonce uint32
	Hash utility.Hash
}
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
func (mb *Mainblock) ComputeRoot() {
	hashes:= make([]utility.Hash, len(mb.Transactions))
	for i:=0;i<len(mb.Transactions);i++{
		hashes[i]=mb.Transactions[i].ComputeHash()
	}
	mb.Header.Roothash=utility.ComputeRoot(&hashes)
}
func (mb *Mainblock) ComputeHash() {
	tmpbuffer := make([]byte, 0)
	bufferVersion := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferVersion, uint32(mb.Header.Version))
	bufferTimestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufferTimestamp, uint64(mb.Header.Timestamp))
	bufferBits := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferBits, uint32(mb.Header.Bits))
	bufferNonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferNonce, uint32(mb.Header.Nonce))

	tmpbuffer=append(tmpbuffer, bufferVersion ...)
	tmpbuffer=append(tmpbuffer, mb.Header.Prevblockhash[:] ...)
	tmpbuffer=append(tmpbuffer, mb.Header.Roothash[:] ...)
	tmpbuffer=append(tmpbuffer, bufferTimestamp ...)
	tmpbuffer=append(tmpbuffer, bufferBits ...)
	tmpbuffer=append(tmpbuffer, bufferNonce ...)
	mb.Header.Hash=utility.ComputeHash(tmpbuffer)
}
func (mb *Mainblock) CheckHash() bool {
	tmpbuffer := make([]byte, 0)
	bufferVersion := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferVersion, uint32(mb.Header.Version))
	bufferTimestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufferTimestamp, uint64(mb.Header.Timestamp))
	bufferBits := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferBits, uint32(mb.Header.Bits))
	bufferNonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferNonce, uint32(mb.Header.Nonce))

	tmpbuffer=append(tmpbuffer, bufferVersion ...)
	tmpbuffer=append(tmpbuffer, mb.Header.Prevblockhash[:] ...)
	tmpbuffer=append(tmpbuffer, mb.Header.Roothash[:] ...)
	tmpbuffer=append(tmpbuffer, bufferTimestamp ...)
	tmpbuffer=append(tmpbuffer, bufferBits ...)
	tmpbuffer=append(tmpbuffer, bufferNonce ...)
	applog.Trace("hash %x computedhash %x",mb.Header.Hash,utility.ComputeHash(tmpbuffer))
	return (mb.Header.Hash==utility.ComputeHash(tmpbuffer))
}
func (mb *Mainblock) Serialize() []byte{
	//TODO switch to binary serialization
	blockstring, err := json.Marshal(mb)
	if err != nil {
		fmt.Println("Mainblock serialize error:", err)
	}
	return blockstring
}
func UnserializeMainblock(bytes []byte)  (*Mainblock,error){
	//TODO binary serialization
	mb:=new(Mainblock)
	err:=json.Unmarshal(bytes,mb)
	if err != nil {
		fmt.Println("Mainblock unserialize error:", err)
		return nil,err
	}
	return mb,nil
}
func (mh *Mainheader) JSONSerialize() []byte{

	blockstring, err := json.Marshal(mh)
	if err != nil {
		fmt.Println("Mainheader serialize error:", err)
	}
	return blockstring
}
func JSONUnserializeMainheader(bytes []byte)  (*Mainheader,error){

	mh:=new(Mainheader)
	err:=json.Unmarshal(bytes,mh)
	if err != nil {
		fmt.Println("Mainheader unserialize error:", err)
		return nil,err
	}
	return mh,nil
}

func (mh *Mainheader) Serialize() []byte{

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(uint32(mh.Version))
	tmpbw.PutBytes(mh.Prevblockhash[:])
	tmpbw.PutBytes(mh.Roothash[:])
	tmpbw.PutUint64(uint64(mh.Timestamp))
	tmpbw.PutUint32(mh.Bits)
	tmpbw.PutUint32(mh.Nonce)
	tmpbw.PutBytes(mh.Hash[:])
	return tmpbw.GetContent()
}
func UnserializeMainheader(bytes []byte)  (*Mainheader,error){

	mh:=new(Mainheader)
	tmpbr:=utility.NewBufferReader(bytes)
	mh.Version= int32(tmpbr.GetUint32())
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Version unserialize error -", err)
		return nil,err
	}

	mh.Prevblockhash=*utility.NewHash(tmpbr.GetBytes(32))
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Prevblockhash unserialize error:", err)
		return nil,err
	}

	mh.Roothash=*utility.NewHash(tmpbr.GetBytes(32))
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Roothash unserialize error:", err)
		return nil,err
	}

	mh.Timestamp= int64(tmpbr.GetUint64())
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Timestamp unserialize error:", err)
		return nil,err
	}

	mh.Bits=tmpbr.GetUint32()
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Bits unserialize error:", err)
		return nil,err
	}

	mh.Nonce=tmpbr.GetUint32()
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Nonce unserialize error:", err)
		return nil,err
	}

	mh.Hash=*utility.NewHash(tmpbr.GetBytes(32))
	if tmpbr.GetError() != nil {
		err:=tmpbr.GetError()
		fmt.Println("Mainheader Hash unserialize error:", err)
		return nil,err
	}
	if !tmpbr.EndOfBytes(){
		return nil,fmt.Errorf("End of bytes not reached")
	}
	
	return mh,nil
}
func GenesisBlock() *Mainblock{
	mb:=new(Mainblock)
	mb.Height=0
	mb.Header.Version=1
	/*
	Version int32 1
	Prevblockhash utility.Hash 0000000000000000000000000000000000000000000000000000000000000000
	
	Timestamp uint32 5f60ddf8
	Bits uint32 1d0fffff
	Nonce uint32 1adc477e
	//
	Hash utility.Hash 0015bd4ac1eb37dc25b9d9ff89bf22015350deb9b7ac5af463b57fbddc0a0100

	MINED success true genesis block &{0 {1 0000000000000000000000000000000000000000000000000000000000000000 583e404bad9b9db80bdf6add27a6ec8b040a0b9d7b03fff0ce644284e4866654 5f60ddf8 daffffff 1adc477e 0015bd4ac1eb37dc25b9d9ff89bf22015350deb9b7ac5af463b57fbddc0a0100} [{1 [] [{2d79883d2000 0100000064931a6cb62b8d4518740dd6656cc3b5d1664ca844f46978c9eccb1b6f50596400}]}]}
	*/

	mb.Header.Timestamp=0x5f60ddf8//
	mb.Header.Bits=0xDAFFFFFF//
	mb.Header.Nonce=0x1adc477e//

	//
	initialhash:=utility.Hash([utility.HashSize]byte{
		0x64, 0x93, 0x1a, 0x6c, 0xb6, 0x2b, 0x8d, 0x45, 
		0x18, 0x74, 0xd, 0xd6, 0x65, 0x6c, 0xc3, 0xb5, 
		0xd1, 0x66, 0x4c, 0xa8, 0x44, 0xf4, 0x69, 0x78, 
		0xc9, 0xec, 0xcb, 0x1b, 0x6f, 0x50, 0x59, 0x64,
	})
	testtx:=utility.NewRewardTransaction(GENESIS_BLOCK_REWARD,0,initialhash)
	mb.Header.Roothash,_= testtx.ComputeSigningHash()//
	
	mb.Header.Hash=utility.Hash([utility.HashSize]byte{
		0x0, 0x15, 0xbd, 0x4a, 0xc1, 0xeb, 0x37, 0xdc, 
		0x25, 0xb9, 0xd9, 0xff, 0x89, 0xbf, 0x22, 0x1, 
		0x53, 0x50, 0xde, 0xb9, 0xb7, 0xac, 0x5a, 0xf4, 
		0x63, 0xb5, 0x7f, 0xbd, 0xdc, 0xa, 0x1, 0x0,
	})
	mb.Transactions=append(mb.Transactions,*testtx)
	return mb
}
func NewMainblock() *Mainblock {
	mb:=new(Mainblock)
	mb.Header.Version=1
	mb.Header.Nonce=0
	//mb.Timestamp=time.Now()/1000
	return mb
}
