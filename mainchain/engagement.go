package mainchain

import (
	"github.com/globaldce/globaldce-toolbox/applog"

	"github.com/globaldce/globaldce-toolbox/utility"
	//"fmt"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
	"math/big"
)


func (mn *Maincore) AddEngagementLikeName(name []byte,stake uint64)  {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameLike)
	nbengagement++
	//totalstake+=stake
	newtotalstake := new(big.Int)
	newtotalstake.Add(&totalstake,big.NewInt(int64(stake)))
	_=mn.PutEngagementNameState(name,StateKeyIdentifierEngagementNameLike,nbengagement,*newtotalstake)
}
func (mn *Maincore) AddEngagementDislikeName(name []byte,stake uint64)  {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameDislike)
	nbengagement++
	//totalstake+=stake
	newtotalstake := new(big.Int)
	newtotalstake.Add(&totalstake,big.NewInt(int64(stake)))
	_=mn.PutEngagementNameState(name,StateKeyIdentifierEngagementNameDislike,nbengagement,*newtotalstake)
}
func (mn *Maincore) GetEngagementLikeName(name []byte) (uint64,big.Int) {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameLike)
	return nbengagement,totalstake
}
func (mn *Maincore) GetEngagementDislikeName(name []byte) (uint64,big.Int) {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameDislike)
	return nbengagement,totalstake
}
//
func (mn *Maincore) AddEngagementPublicPostRewardLike(publicposttxhash utility.Hash,publicposttxindex uint32,stakelike uint64,height uint32) {

	weight:=stakelike*uint64(height)
	_,liketotalstake,disliketotalstake,liketotalweight,disliketotalweight:=mn.GetEngagementPublicPostRewardState(publicposttxhash,publicposttxindex)
	liketotalstake+=stakelike
	newliketotalweight := new(big.Int)
	//liketotalweight+=weight
	bigweight:=new(big.Int)
	bigweight.SetUint64(weight)
	newliketotalweight.Add(&liketotalweight,bigweight)
	err:=mn.PutEngagementPublicPostRewardState(publicposttxhash,publicposttxindex,liketotalstake,disliketotalstake,*newliketotalweight,disliketotalweight)
	_=err
	applog.Trace("publicposttxhash %x,publicposttxindex %d,liketotalstake %d,disliketotalstake %d,*newliketotalweight %s,disliketotalweight %s",publicposttxhash,publicposttxindex,liketotalstake,disliketotalstake,(*newliketotalweight).String(),disliketotalweight.String())

}


func (mn *Maincore) AddEngagementPublicPostRewardDislike(publicposttxhash utility.Hash,publicposttxindex uint32,stakedislike uint64,height uint32) {

	weight:=stakedislike*uint64(height)
	_,liketotalstake,disliketotalstake,liketotalweight,disliketotalweight:=mn.GetEngagementPublicPostRewardState(publicposttxhash,publicposttxindex)
	disliketotalstake+=stakedislike
	//disliketotalweight+=weight
	newdisliketotalweight:=new(big.Int)
	bigweight:=new(big.Int)
	bigweight.SetUint64(weight)
	newdisliketotalweight.Add(&disliketotalweight,bigweight)
	err:=mn.PutEngagementPublicPostRewardState(publicposttxhash,publicposttxindex,liketotalstake,disliketotalstake,liketotalweight,*newdisliketotalweight)
	_=err
	
}
/*
func (mn *Maincore) AddEngagementLikePublicPost(name []byte,stake uint64)  {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameLike)
	nbengagement++
	totalstake+=stake
	_=mn.PutEngagementNameState(name,StateKeyIdentifierEngagementNameLike,nbengagement,totalstake)
}

func (mn *Maincore) GetEngagementLikePub(name []byte) (uint32,uint64) {
	nbengagement,totalstake:=mn.GetEngagementNameState(name,StateKeyIdentifierEngagementNameLike)
	return nbengagement,totalstake
}
*/