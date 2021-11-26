package mainchain
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	//"github.com/globaldce/globaldce-toolbox/applog"
	//"github.com/globaldce/globaldce-toolbox/wire"
	//"github.com/globaldce/globaldce-toolbox/mainchain"
	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/globaldce/globaldce-toolbox/utility"
	//"github.com/globaldce/globaldce-toolbox/mainchain"
	//"github.com/globaldce/globaldce-toolbox/wallet"
	//"os"
	//"math"
	//"time"
	"encoding/json"
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
	//DataIdentifierDynamicPublicPost

)
type PostInfo struct {
	Name string
	Link string
	Content string

	//user    *user
}

func StringFromPostInfo(p PostInfo) string{

	//json.Unmarshal([]byte(stringData), &data)
	b,_:=json.Marshal(p)
	return string(b)
	
}

func (mn *Maincore) GetPostInfoStringArray(maxposts int)[]string{
	var postsstringarray []string
	nbdata:=int(mn.GetNbData())
	starti:=nbdata-maxposts
	if starti<0{
		starti=0
	}

	for i:=starti;i<nbdata;i++{
		databytes:=mn.GetDataById(i)
		tmpbr:=utility.NewBufferReader(databytes)

		dataidentifier:=tmpbr.GetUint32()
		if dataidentifier==DataIdentifierPublicPost{
			//namebyteslen:=tmpbr.GetVarUint()
			//namebytes:=tmpbr.GetBytes(uint(namebyteslen))
			//namestring:=string(namebytes)
	
			linkbyteslen:=tmpbr.GetVarUint()
			linkbytes:=tmpbr.GetBytes(uint(linkbyteslen))
			linkstring:=string(linkbytes)
			textbyteslen:=tmpbr.GetVarUint()
			textbytes:=tmpbr.GetBytes(uint(textbyteslen))
			textstring:=string(textbytes)
			ed:=utility.NewExtradataFromBytes(databytes)

			namebytes,_,err:=mn.GetPublicPostState(ed.Hash)
			_=err
			//if err!=nil {
			//	applog.Warning("Cannot add data - hash %s - error %v",hash,err)
			//	return
			//}
			namestring:=string(namebytes)
			tmpstring:=StringFromPostInfo(PostInfo{Name:namestring,Link:linkstring,Content:textstring})
			postsstringarray=append(postsstringarray,tmpstring)
		}

		
	}
	return postsstringarray
}

func (mn *Maincore) AddLocalPublicPostData(namestring string,hash utility.Hash,databytes []byte) {
	namebytes:=[]byte(namestring)
	//tmpbw:=utility.NewBufferWriter()
	//tmpbw.PutUint32(DataIdentifierPublicPost)
	//tmpbw.PutRegistredNameKey()
	//tmpbw.PutVarUint(uint64(len(namebytes)))
	//tmpbw.PutBytes(namebytes)
	//tmpbw.PutVarUint(uint64(len(databytes)))
	//tmpbw.PutBytes(databytes)
	// localy generated data can be directly stored
	//mn.dataf.AddChunk(tmpbw.GetContent())
	mn.dataf.AddChunk(databytes)
	mn.PutPublicPostState(hash,namebytes,uint32(mn.dataf.NbChunks()-1))

}

func (mn *Maincore) addData(hash utility.Hash,bytes []byte) {
	//TODO generalize to different data types

	name,id,err:=mn.GetPublicPostState(hash)
	if err!=nil {
		applog.Warning("Cannot add data - hash %s - error %v",hash,err)
		return
	}
	if id!=0 {
		applog.Warning("Cannot add data - hash %s - data already exist stored with id %d",hash,id)
		return
	}
	//name,data,_:=mn.GetPublicPostData(hash)
	//if data!=nil {
	//	applog.Warning("Cannot add data - hash %s - data already exist",hash)
	//	return
	//}
	mn.dataf.AddChunk(bytes)
	mn.PutPublicPostState(hash,name,uint32(mn.dataf.NbChunks()-1))
}
//TODO generalize to different data types
func (mn *Maincore) GetPublicPostData(hash utility.Hash) ([]byte,[]byte,error) {
	name,id,err:=mn.GetPublicPostState(hash)
	if err!=nil{
		return nil,nil,err
	}
	return name,mn.dataf.GetChunkById(int(id)),nil
}
func (mn *Maincore) GetDataById(id int)  []byte {
	return mn.dataf.GetChunkById(id)
}
func (mn *Maincore) GetNbData() uint32{
	//applog.Trace("-------%v",mn)
	if mn.dataf == nil{
		return uint32(0)
	}
	return uint32(mn.dataf.NbChunks())
}