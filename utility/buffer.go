package utility

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type BufferWriter struct {
	content []byte
}

func NewBufferWriter() *BufferWriter {
	bw:=new(BufferWriter)
	//br.content= make([]byte, 0)
	return bw
}

func (bw * BufferWriter) PutVarUint(i uint64){
	if i<253 {
		bw.PutUint8(uint8(i))
	} else if i<= math.MaxUint16 {
		bw.PutUint8(253)
		bw.PutUint16(uint16(i))
	} else if i<= math.MaxUint32 {
		bw.PutUint8(254)
		bw.PutUint32(uint32(i))
	} else  {
		bw.PutUint8(255)
		bw.PutUint64(uint64(i))
	}
}
func (bw * BufferWriter) PutUint8(i uint8){
	tmpbufferuint := make([]byte, 1)
	tmpbufferuint[0] = byte(i)
	bw.content=append(bw.content, tmpbufferuint ...)
}
func (bw * BufferWriter) PutUint16(i uint16){
	tmpbufferuint := make([]byte, 2)
	binary.LittleEndian.PutUint16(tmpbufferuint, i)
	bw.content=append(bw.content, tmpbufferuint ...)
}
func (bw * BufferWriter) PutUint32(i uint32){
	tmpbufferuint := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpbufferuint, i)
	bw.content=append(bw.content, tmpbufferuint ...)
}
func (bw * BufferWriter) PutUint64(i uint64){
	tmpbufferuint := make([]byte, 8)
	binary.LittleEndian.PutUint64(tmpbufferuint, i)
	bw.content=append(bw.content, tmpbufferuint ...)
}

func (bw * BufferWriter) PutBigInt(i * big.Int){

	//tmpbufferuint := make([]byte, 8)
	//binary.LittleEndian.PutUint64(tmpbufferuint, i)
	buf:=i.Bytes()
	bw.PutVarUint(uint64(len(buf)))
	bw.content=append(bw.content, buf ...)
}

func (bw * BufferWriter) PutBytes(buf []byte){
	//tmpbufferuint := make([]byte, 4)
	//binary.LittleEndian.PutUint32(tmpbufferuint, i)
	bw.content=append(bw.content, buf ...)
}
func (bw * BufferWriter) PutHash(h Hash){
	bw.PutBytes(h[:])
}
func (bw * BufferWriter) PutRegistredNameKey(name []byte){
	var newname [RegistredNameMaxSize]byte
	copy(newname[:],name)
	bw.PutBytes(newname[:])
}
func (bw * BufferWriter) GetContent() []byte{
	return bw.content
}
//
//

type BufferReader struct {
	content []byte
	counter int
	err error 
}

func NewBufferReader(content []byte) *BufferReader {
	br:=new(BufferReader)
	br.counter=0
	br.content= content
	return br
}
func (br * BufferReader) GetError() error{
	return br.err
}
func (br * BufferReader) GetCounter() uint{
	return uint(br.counter)
}
func (br * BufferReader) EndOfBytes() bool{
	return (br.counter==len(br.content))
}
////////

func (br * BufferReader) GetVarUint() uint64{
	i:=br.GetUint8()
	if i<253 {
		return uint64(i)
		//bw.PutUint8(uint8(i))
	} else if i== 253 {
		return uint64(br.GetUint16())
	} else if i== 254 {
		return uint64(br.GetUint32())
	} else {
		return uint64(br.GetUint64())
	}
}

func (br * BufferReader) GetUint8() uint8{
	br.counter+=1
	if br.counter>len(br.content){
		br.err=fmt.Errorf("end of bytes reached")
		return uint8(0)
	}
	return uint8(br.content[br.counter-1])
}
//////////////////

func (br * BufferReader) GetUint16() uint16{
	br.counter+=2
	if br.counter>len(br.content){
		br.err=fmt.Errorf("end of bytes reached")
		return uint16(0)
	}
	return binary.LittleEndian.Uint16(br.content[br.counter-2:br.counter])
}

func (br * BufferReader) GetUint32() uint32{
	br.counter+=4
	if br.counter>len(br.content){
		br.err=fmt.Errorf("end of bytes reached")
		return uint32(0)
	}
	return binary.LittleEndian.Uint32(br.content[br.counter-4:br.counter])
}

func (br * BufferReader) GetUint64() uint64{
	br.counter+=8
	if br.counter>len(br.content){
		br.err=fmt.Errorf("End of bytes reached")
		return uint64(0)
	}
	return binary.LittleEndian.Uint64(br.content[br.counter-8:br.counter])
}
func (br * BufferReader) GetBigInt()  (* big.Int){
	l:=br.GetVarUint()
	buf:=br.GetBytes(uint(l))
	i := new(big.Int)
	i.SetBytes(buf)
	return i
}
func (br * BufferReader) GetBytes(length uint) []byte{
	br.counter+=int(length)
	if br.counter>len(br.content){
		br.err=fmt.Errorf("End of bytes reached")
		return nil
	}
	return br.content[br.counter-int(length):br.counter]
}
func (br * BufferReader) GetHash() Hash{
hb:= br.GetBytes(HashSize)
 return *NewHash(hb)
}

const ExtrabytesMaxSize=100

func (br * BufferReader) GetExtrabytes() []byte{
	//var extrabytes []byte
	extrabyteslen:=br.GetVarUint()
	if extrabyteslen>ExtrabytesMaxSize{
		br.err=fmt.Errorf("ExtrabytesMaxSize exceeded")
		return nil
	}

	return br.GetBytes(uint(extrabyteslen))
}

