package wire
import(
	"github.com/globaldce/globaldce-toolbox/utility"
	"github.com/globaldce/globaldce-toolbox/applog"
	"github.com/globaldce/globaldce-toolbox/mainchain"
)
const (
	RequestMainheadersMax=uint32(100)
)

//
func EncodeBroadcastTransaction(NbHops uint32,tx *utility.Transaction ) (*Message){
	msg:=NewMessage(MsgIdentifierBroadcastTransaction)
	contentbw:=utility.NewBufferWriter()
	contentbw.PutUint32(NbHops)
	tmpbytes:=tx.Serialize()
	contentbw.PutUint32(uint32(len(tmpbytes)))
	contentbw.PutBytes(tmpbytes)
	msg.PutContent(contentbw.GetContent())
	return msg
}
func DecodeBroadcastTransaction(msg *Message)(bool,uint32,uint32,*utility.Transaction){
	tmpbr:=utility.NewBufferReader(msg.GetContent())
	var tx *utility.Transaction
	nbhops:=tmpbr.GetUint32()

		txbyteslength:=	tmpbr.GetUint32()
			tx,serr:=utility.UnserializeTransaction(tmpbr.GetBytes(uint(txbyteslength)))
			if serr!=nil {
				return false,uint32(0),uint32(0),nil
			}
	
	

	return true,nbhops,txbyteslength,tx
}
//

func EncodeBroadcastMainblock(NbHops uint32,mb *mainchain.Mainblock ) (*Message){
	msg:=NewMessage(MsgIdentifierBroadcastMainblock)
	contentbw:=utility.NewBufferWriter()
	//contentbw.PutUint32(uint32(4+4))//content length

	contentbw.PutUint32(NbHops)
	contentbw.PutUint32(mb.Height)
	headertmpbytes:=mb.Header.Serialize()

	contentbw.PutUint32(uint32(len(headertmpbytes)))
	contentbw.PutBytes(headertmpbytes)

	for i:=0;i<len(mb.Transactions);i++{
		
		tmpbytes:=mb.Transactions[i].Serialize()
		contentbw.PutUint32(uint32(len(tmpbytes)))

		contentbw.PutBytes(tmpbytes)
	}
	msg.PutContent(contentbw.GetContent())
	return msg
}
func DecodeBroadcastMainblock(msg *Message)(bool,uint32,uint32,*mainchain.Mainblock){
	tmpbr:=utility.NewBufferReader(msg.GetContent())
	var mbtxs []utility.Transaction
	nbhops:=tmpbr.GetUint32()
	height:=tmpbr.GetUint32()
	var mb mainchain.Mainblock
	byteslength:=uint(len(msg.GetContent()))
	mhbyteslength:=	tmpbr.GetUint32()
		mh,err:=mainchain.UnserializeMainheader(tmpbr.GetBytes(uint(mhbyteslength)))
		if err!=nil {
			return false,uint32(0),uint32(0),nil
		}
	for tmpbr.GetCounter() < byteslength{
		txbyteslength:=	tmpbr.GetUint32()
			tx,serr:=utility.UnserializeTransaction(tmpbr.GetBytes(uint(txbyteslength)))
			if serr!=nil {
				return false,uint32(0),uint32(0),nil
			}
			mbtxs=append(mbtxs,*tx)
	}
	mb.Header=*mh
	mb.Height=height
	mb.Transactions=mbtxs
	return true,nbhops,height,&mb
}


func EncodeRequestMainheaders(first uint32, last uint32) (*Message){
	msg:=NewMessage(MsgIdentifierRequestMainheaders)
	contentbw:=utility.NewBufferWriter()
	//contentbw.PutUint32(uint32(4+4))//content length
	contentbw.PutUint32(first)
	contentbw.PutUint32(last)
	
	
	//data:=[]byte("coool")
	
	//tmpbw:=utility.NewBufferWriter()
	//tmpbw.PutUint32(uint32(len(data)))
	//tmpbw.PutBytes(data)

	msg.PutContent(contentbw.GetContent())
	return msg
}
func DecodeRequestMainheaders(msg *Message)(bool,uint32,uint32){
	tmpbr:=utility.NewBufferReader(msg.GetContent())

	first:=tmpbr.GetUint32()
	last:=tmpbr.GetUint32()
	if (last<first)||(last-first>RequestMainheadersMax){
		applog.Trace("error first %d last %d",first,last)
		return false,0,0
	}
	
	return true,first,last
}

//
func EncodeRequestMainblockTransactions(blockheight uint32) (*Message){
	msg:=NewMessage(MsgIdentifierRequestMainblockTransactions)
	contentbw:=utility.NewBufferWriter()
	//contentbw.PutUint32(uint32(4+4))//content length
	contentbw.PutUint32(blockheight)
	
	
	//data:=[]byte("coool")
	
	//tmpbw:=utility.NewBufferWriter()
	//tmpbw.PutUint32(uint32(len(data)))
	//tmpbw.PutBytes(data)

	msg.PutContent(contentbw.GetContent())
	return msg
}
func DecodeRequestMainblockTransactions(msg *Message)(bool,uint32){
	tmpbr:=utility.NewBufferReader(msg.GetContent())
	// TODO Check min length of buffer reader
	requestedblockheight:=tmpbr.GetUint32()
	/*
	last:=tmpbr.GetUint32()
	if (last<first)||(last-first>RequestMainheadersMax){
		applog.Trace("error first %d last %d",first,last)
		return false,0,0
	}
	*/
	
	return true,requestedblockheight
}