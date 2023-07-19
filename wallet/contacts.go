package wallet

//import
//(
//"github.com/globaldce/globaldce-gateway/utility"
//"encoding/json"
//"encoding/binary"
//"bytes"

//"github.com/globaldce/globaldce-gateway/applog"
//"fmt"
//"os"
//"sync"
//)

type Contact struct {
	Name         string
	AddrString   string
	GroupIdArray []uint32
}

func (wlt *Wallet) AddContact(tmpname string, tmpaddrstring string, tmpgroupidarray []uint32) {
	//var emptytxhash utility.Hash
	tmpcontact := Contact{
		Name:         tmpname,
		AddrString:   tmpaddrstring,
		GroupIdArray: tmpgroupidarray,
	}
	wlt.Contactarray = append(wlt.Contactarray, tmpcontact)
}
