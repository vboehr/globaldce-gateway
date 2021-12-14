package mainchain
/*
import (
	//"github.com/globaldce/globaldce-toolbox/applog"

	//"github.com/globaldce/globaldce-toolbox/utility"
	"fmt"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)
*/
const (
	EngagementIdentifierLikeName=1
	EngagementIdentifierDislikeName=2
)

func (mn *Maincore) AddEngagementLikeName(name []byte)  {
	value:=mn.GetEngagementState(name,EngagementIdentifierLikeName)
	value++
	_=mn.PutEngagementState(name,EngagementIdentifierLikeName,value)
}
func (mn *Maincore) AddEngagementDislikeName(name []byte)  {
	value:=mn.GetEngagementState(name,EngagementIdentifierDislikeName)
	value++
	_=mn.PutEngagementState(name,EngagementIdentifierDislikeName,value)
}
func (mn *Maincore) GetEngagementLikeName(name []byte) uint32 {
	value:=mn.GetEngagementState(name,EngagementIdentifierLikeName)
	return value
}
func (mn *Maincore) GetEngagementDislikeName(name []byte) uint32 {
	value:=mn.GetEngagementState(name,EngagementIdentifierDislikeName)
	return value
}
