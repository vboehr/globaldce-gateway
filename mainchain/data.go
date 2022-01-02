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
	"time"
	"encoding/json"
	//"math/big"
	//"net"
	//"log"
	"fmt"
	//"path/filepath"
	//"sync"
	"strings"
	"math/rand"
)

const (
	DataIdentifierPublicPost=1
	//DataIdentifierDynamicPublicPost

)
type PostInfo struct {
	Name string
	Link string
	Content string
	AttachmentSizeArray []int
	AttachmentHashArray []utility.Hash
	//user    *user
}

func StringFromPostInfo(p PostInfo) string{

	//json.Unmarshal([]byte(stringData), &data)
	b,_:=json.Marshal(p)
	return string(b)
	
}

func (mn *Maincore) GetPostInfoStringArray(keywords string,maxposts int)[]string{
	var keywordarray []string
	var postsstringarray []string
	keywordarray=strings.Split(keywords, " ")
	nbdata:=int(mn.GetNbData())
	starti:=nbdata-maxposts
	if starti<0{
		starti=0
	}

	for i:=nbdata-1;i>=starti;i--{
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
			nbattachement:=tmpbr.GetVarUint()
			var tmpAttachmentSizeArray []int
			var tmpAttachmentHashArray []utility.Hash
			for j:=0;j<int(nbattachement);j++ {
				tmpAttachmentSize:=int(tmpbr.GetVarUint())
				tmpAttachmentHash:=tmpbr.GetHash()
				tmpAttachmentSizeArray=append(tmpAttachmentSizeArray,tmpAttachmentSize)
				tmpAttachmentHashArray=append(tmpAttachmentHashArray,tmpAttachmentHash)
			}
			//
			ed:=utility.NewExtradataFromBytes(databytes)

			namebytes,_,err:=mn.GetPublicPostState(ed.Hash)
			_=err
			//if err!=nil {
			//	applog.Warning("Cannot add data - hash %s - error %v",hash,err)
			//	return
			//}
			
			//Block namestring should not be displayed
			if mn.IsBannedName(namebytes) {
				continue
			}
			namestring:=string(namebytes)
			tmpstring:=StringFromPostInfo(PostInfo{
				Name:namestring,
				Link:linkstring,
				Content:textstring,
				AttachmentSizeArray:tmpAttachmentSizeArray,
				AttachmentHashArray:tmpAttachmentHashArray,
			})
			if len(keywordarray)!=0{
				for _,k:=range keywordarray{
					if strings.Index(tmpstring,k)>=0{
						postsstringarray=append(postsstringarray,tmpstring)
					}
				}
			} else {
				postsstringarray=append(postsstringarray,tmpstring)
			}

			
		}

		
	}
	return postsstringarray
}
func (mn *Maincore) IsBannedName(namebytes []byte) bool{
	name:=string(namebytes)
	for _ , bn:= range mn.BannedNameArray{
		if bn==name{
			return true
		}
	}
	return false
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

func (mn *Maincore) AddToMissingDataHashArray(hash utility.Hash) {
	var alreadyexist=false

	for _,h:=range mn.MissingDataHashArray{
		if h == hash{
			alreadyexist=true
			applog.Warning("Data already existes %x",hash)
			return
		}
	}
	if !alreadyexist{
		mn.MissingDataHashArray=append(mn.MissingDataHashArray,hash)
	}
}
func (mn *Maincore) AddToMissingDataFileHashArray(hash utility.Hash) {
	var alreadyexist=false

	for _,h:=range mn.MissingDataFileHashArray{
		if h == hash{
			alreadyexist=true
			applog.Warning("Data already existes %x",hash)
			return
		}
	}
	if !alreadyexist{
		mn.MissingDataFileHashArray=append(mn.MissingDataFileHashArray,hash)
	}
}
func (mn *Maincore) RemoveMissingDataHash(hash utility.Hash) {
	for i,h:=range mn.MissingDataHashArray{
		if h == hash{
			mn.MissingDataHashArray=append(mn.MissingDataHashArray[:i], mn.MissingDataHashArray[i+1:]...)
			return
		}
	}
}
func (mn *Maincore) RemoveMissingDataFileHash(hash utility.Hash) {
	for i,h:=range mn.MissingDataFileHashArray{
		if h == hash{
			mn.MissingDataFileHashArray=append(mn.MissingDataFileHashArray[:i], mn.MissingDataFileHashArray[i+1:]...)
			return
		}
	}
}

func (mn *Maincore) AddData(hash utility.Hash,bytes []byte) {
	//TODO generalize to different data types

	name,id,err:=mn.GetPublicPostState(hash)
	if err!=nil {
		applog.Warning("Cannot add data - hash %x - error %v",hash,err)
		return
	}
	if id!=0 {
		applog.Warning("Cannot add data - hash %x - data already exist stored with id %d",hash,id)
		return
	}
	//name,data,_:=mn.GetPublicPostData(hash)
	//if data!=nil {
	//	applog.Warning("Cannot add data - hash %s - data already exist",hash)
	//	return
	//}
	mn.dataf.AddChunk(bytes)
	mn.PutPublicPostState(hash,name,uint32(mn.dataf.NbChunks()-1))
	mn.RemoveMissingDataHash(hash)
	mn.UpdateMissingDataFileHashArray(bytes)
}
func (mn *Maincore) GetData(hash utility.Hash) ([]byte,error) {
	name,data,err:=mn.GetPublicPostData(hash)
	if err!=nil{
		return ([]byte("")),err
	}
	if mn.IsBannedName(name) {
		return ([]byte("")),fmt.Errorf("Data associated with banned name")
	}
	return data,nil
}
func (mn *Maincore) UpdateMissingDataFileHashArray(databytes []byte){

		tmpbr:=utility.NewBufferReader(databytes)

		dataidentifier:=tmpbr.GetUint32()
		if dataidentifier!=DataIdentifierPublicPost{
			return //,fmt.Errorf("Databytes is not Public Post Data")
		}
			//namebyteslen:=tmpbr.GetVarUint()
			//namebytes:=tmpbr.GetBytes(uint(namebyteslen))
			//namestring:=string(namebytes)
	
			linkbyteslen:=tmpbr.GetVarUint()
			linkbytes:=tmpbr.GetBytes(uint(linkbyteslen))
			_=linkbytes

			textbyteslen:=tmpbr.GetVarUint()
			textbytes:=tmpbr.GetBytes(uint(textbyteslen))
			_=textbytes
			nbattachement:=tmpbr.GetVarUint()
			//var tmpAttachmentSizeArray []int
			//var tmpAttachmentHashArray []utility.Hash
			for j:=0;j<int(nbattachement);j++ {
				tmpAttachmentSize:=int(tmpbr.GetVarUint())
				tmpAttachmentHash:=tmpbr.GetHash()
				_=tmpAttachmentSize
				mn.AddToMissingDataFileHashArray(tmpAttachmentHash)
				//tmpAttachmentSizeArray=append(tmpAttachmentSizeArray,tmpAttachmentSize)
				//tmpAttachmentHashArray=append(tmpAttachmentHashArray,tmpAttachmentHash)
			}

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
func (mn *Maincore) GetRandomMissingDataHash() *utility.Hash{
	rand.Seed(time.Now().UnixNano())
	if len(mn.MissingDataHashArray)==0{
		return nil
	}
	i:=rand.Intn(len(mn.MissingDataHashArray))//TODO concentrate on the last Data Hashes
	return &mn.MissingDataHashArray[i]
}
func (mn *Maincore) GetRandomMissingDataFileHash() *utility.Hash{
	rand.Seed(time.Now().UnixNano())
	if len(mn.MissingDataFileHashArray)==0{
		return nil
	}
	i:=rand.Intn(len(mn.MissingDataFileHashArray))//TODO concentrate on the last DataFile Hashes
	return &mn.MissingDataFileHashArray[i]
}
func (mn *Maincore) IsMissingData(hash utility.Hash) bool{
	for _,missingdatahash:=range mn.MissingDataHashArray{
		if missingdatahash==hash {
			return true
		}
	}
	return false
}
func (mn *Maincore) IsMissingDataFile(hash utility.Hash) bool{
	for _,missingdatafilehash:=range mn.MissingDataFileHashArray{
		if missingdatafilehash==hash {
			return true
		}
	}
	return false
}