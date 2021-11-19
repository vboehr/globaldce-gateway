package mainchain
import (
	//"github.com/globaldce/globaldce/applog"
	//"github.com/globaldce/globaldce/applog"
	//"github.com/globaldce/globaldce/wire"
	//"github.com/globaldce/globaldce/mainchain"
	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/globaldce/globaldce/utility"
	//"github.com/globaldce/globaldce/mainchain"
	//"github.com/globaldce/globaldce/wallet"
	//"os"
	//"math"
	//"time"
	//"encoding/json"
	//"math/big"
	//"net"
	//"log"
	//"fmt"
	//"path/filepath"
	//"sync"
	//"strings"
)

const (
	DataIdentifierPublicPost=1

)

func (mn *Maincore) GetPosts(maxposts int)[]string{
	var postsstringarray []string
	nbdata:=int(mn.GetNbData())
	starti:=nbdata-maxposts
	if starti<0{
		starti=0
	}

	for i:=starti;i<nbdata;i++{
		
		tmpbr:=utility.NewBufferReader(mn.GetDataById(i))

		dataidentifier:=tmpbr.GetUint32()
		if dataidentifier==DataIdentifierPublicPost{
			namebyteslen:=tmpbr.GetVarUint()
			namebytes:=tmpbr.GetBytes(uint(namebyteslen))
			databyteslen:=tmpbr.GetVarUint()
			databytes:=tmpbr.GetBytes(uint(databyteslen))
			tmpstring:=string(namebytes)+" "+string(databytes)
			postsstringarray=append(postsstringarray,tmpstring)
		}

		
	}
	return postsstringarray
}

func (mn *Maincore) AddPublicPostData(namestring string,hash utility.Hash,databytes []byte) {
	namebytes:=[]byte(namestring)
	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(DataIdentifierPublicPost)
	//tmpbw.PutRegistredNameKey()
	tmpbw.PutVarUint(uint64(len(namebytes)))
	tmpbw.PutBytes(namebytes)
	tmpbw.PutVarUint(uint64(len(databytes)))
	tmpbw.PutBytes(databytes)

	mn.dataf.AddChunk(tmpbw.GetContent())
	mn.PutDataState(hash,uint32(mn.dataf.NbChunks()-1))
}

func (mn *Maincore) addData(hash utility.Hash,bytes []byte) {
	mn.dataf.AddChunk(bytes)
	mn.PutDataState(hash,uint32(mn.dataf.NbChunks()-1))
}
func (mn *Maincore) GetData(hash utility.Hash)  []byte {
	id,err:=mn.GetDataState(hash)
	if err!=nil{
		return nil
	}
	return mn.dataf.GetChunkById(int(id))
}
func (mn *Maincore) GetDataById(id int)  []byte {
	return mn.dataf.GetChunkById(id)
}
func (mn *Maincore) GetNbData() uint32{
	if mn.dataf == nil{
		return uint32(0)
	}
	return uint32(mn.dataf.NbChunks())
}