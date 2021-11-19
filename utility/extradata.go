package utility

//import (
//	"fmt"
//)
const ExtradataMaxSize=200

type Extradata struct {
	Size uint64
	Hash Hash
}
func NewExtradataFromBytes(data []byte) Extradata {
	ed:=new(Extradata)
	ed.Size=uint64(len(data))
	ed.Hash=ComputeHash(data)
	return *ed
}
func NewExtradata() Extradata {
	ed:=new(Extradata)
	ed.Size=uint64(0)
	return *ed
}