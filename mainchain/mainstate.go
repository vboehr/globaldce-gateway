package mainchain

import (
	"github.com/globaldce/globaldce-gateway/applog"

	"github.com/globaldce/globaldce-gateway/utility"
	//"fmt"
	"os"
	//"math/big"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)

const (
	StateKeyIdentifierTxOutput=1
	StateKeyIdentifierTx=2
	StateKeyIdentifierMainblock=3
	StateKeyIdentifierNameRegistration=4

	StateKeyIdentifierAddressBalance=5
	//StateKeyIdentifierAddressNbAssets
	//StateKeyIdentifierAddressAsset

	
	//StateKeyIdentifierData=6
	//StateKeyIdentifierDataFile=7
	//StateKeyIdentifierEngagementName=8
	//StateKeyIdentifierEngagementNameLike=9
	//StateKeyIdentifierEngagementNameDislike=10
	//StateKeyIdentifierEngagementPublicPost=11
	//StateKeyIdentifierEngagementPublicPostLike=12
	//StateKeyIdentifierEngagementPublicPostDislike=13
	//StateKeyIdentifierEngagementPublicPostReward=14
)
const (
	StateValueIdentifierNotFound=0

	StateValueIdentifierUnspentTxOutput=1001
	StateValueIdentifierSpentTxOutput=1002

	StateValueIdentifierTx=2001

	StateValueIdentifierValidMainblock=3001
	StateValueIdentifierInvalidMainblock=3002

	StateValueIdentifierActifNameRegistration=4001
	StateValueIdentifierInactifNameRegistration=4002

	StateValueIdentifierAddressBalance=5001

	//StateValueIdentifierData=6
	//StateValueIdentifierDataFile=7
	//StateValueIdentifierEngagementName=8
	//StateValueIdentifierUnclaimedEngagementPublicPost=9
	//StateValueIdentifierClaimedEngagementPublicPost=10
	//StateValueIdentifierEngagementPublicPostReward=11
)

func (mn *Maincore) PutTxState(txhash utility.Hash,height uint32,number uint32) error{
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierTx)
	tmpkeybw.PutBytes(txhash[:])
	//tmpkeybw.GetContent()

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierTx)// transaction type
	tmpbw.PutUint32(height)
	tmpbw.PutUint32(number)
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetTxState(txhash utility.Hash) (uint32,uint32,uint32) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierTx)
	tmpkeybw.PutBytes(txhash[:])
	//tmpkeybw.GetContent()

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetTxState - txhash %x : %v",txhash, err)
		return StateValueIdentifierNotFound,0,0
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	height:=tmpbr.GetUint32()
	number:=tmpbr.GetUint32()

	applog.Trace("type %d height %d number %d ",stateidentifier,height,number)
	return stateidentifier,height,number
}


func (mn *Maincore) PutTxOutputState(txhash utility.Hash,index uint32,stateidentifier uint32) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierTxOutput)
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(stateidentifier)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetTxOutputState(txhash utility.Hash,index uint32) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierTxOutput)
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetTxOutputState - txhash %x - index %d : %v",txhash,index, err)
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	return stateidentifier
}


func (mn *Maincore) PutNameState(name []byte,stateidentifier uint32) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierNameRegistration)
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(stateidentifier)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetNameState(name []byte) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierNameRegistration)
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	return stateidentifier
}


func (mn *Maincore) PutAddressBalanceState(addr utility.Hash,amount uint64) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierAddressBalance)
	tmpkeybw.PutBytes(addr[:])
	//tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateKeyIdentifierAddressBalance)
	tmpbw.PutUint64(amount)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetAddressBalanceState(addr utility.Hash) (uint32,uint64) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierAddressBalance)
	tmpkeybw.PutBytes(addr[:])
	//tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetAddressBalanceState - address %x : %v",addr, err)
		return StateValueIdentifierNotFound,0
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	amount:=tmpbr.GetUint64()
	return stateidentifier,amount
}

//////////////////////////////
func (mn *Maincore) AddToAddressBalance(addr utility.Hash,amount uint64) error {

	_,balance:=mn.GetAddressBalanceState(addr)
	balance+=amount
	_=mn.PutAddressBalanceState(addr,balance)

	return nil
}
func (mn *Maincore) SubtractFromAddressBalance(addr utility.Hash,amount uint64) error {

	_,balance:=mn.GetAddressBalanceState(addr)
	balance-=amount
	_=mn.PutAddressBalanceState(addr,balance)

	return nil
}
//////////////////////////////
func (mn *Maincore) PutMainblockState(height uint32,state uint32) error{
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierMainblock)
	tmpkeybw.PutUint32(height)
	//tmpkeybw.GetContent()

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(state)
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetMainblockState(height uint32) (uint32) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierTx)
	tmpkeybw.PutUint32(height)
	//tmpkeybw.GetContent()

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetMainblockStat height %d: %v",height, err)
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()


	applog.Trace("type %d ",stateidentifier)
	return stateidentifier
}
////////////////////////////////////////

/*
func (mn *Maincore) PutPublicPostState(datahash utility.Hash,namebytes []byte,txhash utility.Hash,index uint32,id uint32) error{
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierData)
	tmpkeybw.PutHash(datahash)
	//tmpkeybw.GetContent()

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierData)// 
	tmpbw.PutVarUint(uint64(len(namebytes)))
	tmpbw.PutBytes(namebytes)
	tmpbw.PutHash(txhash)
	tmpbw.PutUint32(index)
	tmpbw.PutUint32(id)
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}
func (mn *Maincore) GetPublicPostState(datahash utility.Hash)([]byte,*utility.Hash,uint32, uint32, error){
	

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierData)
	tmpkeybw.PutHash(datahash)
	//tmpkeybw.PutUint32(index)

	valuebytes, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return nil,nil,0,0,err
	}
	tmpbr:=utility.NewBufferReader(valuebytes)
	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierData{
		return nil,nil,0,0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
	}
	namebyteslen:=tmpbr.GetVarUint()
	namebytes:=tmpbr.GetBytes(uint(namebyteslen))
	txhash:=tmpbr.GetHash()
	index:=tmpbr.GetUint32()
	id:=tmpbr.GetUint32()
	return namebytes,&txhash,index,id,nil
}

*/
/*

func (mn *Maincore) PutDataFileState(datafilehash utility.Hash,size uint64) error{
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierDataFile)
	tmpkeybw.PutHash(datafilehash)
	//tmpkeybw.GetContent()

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierData)// transaction type
	tmpbw.PutVarUint(size)
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetDataFileState(datafilehash utility.Hash)(uint64, error){
	

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierDataFile)
	tmpkeybw.PutHash(datafilehash)
	//tmpkeybw.PutUint32(index)

	valuebytes, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return 0,err
	}
	tmpbr:=utility.NewBufferReader(valuebytes)
	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierDataFile{
		return 0,fmt.Errorf("Found an incorrect stateidentifier associated with data file hash %x - identifier found")
	}
	size:=tmpbr.GetVarUint()

	return size,nil
}


//
func (mn *Maincore) PutEngagementNameState(name []byte,engagementidentifier uint32,nbengagement uint64,totalstake big.Int) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementName)
	tmpkeybw.PutRegistredNameKey(name)

	tmpkeybw.PutUint32(engagementidentifier)

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierEngagementName)
	tmpbw.PutUint64(nbengagement)
	tmpbw.PutBigInt(&totalstake)
	

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)

	return err
}

func (mn *Maincore) GetEngagementNameState(name []byte,engagementidentifier uint32) (uint64,big.Int) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementName)
	tmpkeybw.PutRegistredNameKey(name)
	tmpkeybw.PutUint32(engagementidentifier)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return 0,*big.NewInt(0)
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierEngagementName{
		//return nil,0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
		return 0,*big.NewInt(0)
	}
	nbengagement:=tmpbr.GetUint64()
	totalstake:=tmpbr.GetBigInt()
	return nbengagement,*totalstake // returns nbengagement and totalstake
}
*/
/*
func (mn *Maincore) PutEngagementPublicPostState(pptxhash utility.Hash,pptxindex uint32,claimaddress utility.Hash,engagementidentifier uint32,etxhash utility.Hash,etxindex uint32) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementPublicPost)
	tmpkeybw.PutHash(pptxhash)
	tmpkeybw.PutUint32(pptxindex)
	tmpkeybw.PutHash(claimaddress)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(engagementidentifier)
	tmpbw.PutHash(etxhash)
	tmpbw.PutUint32(etxindex)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}*/
/*
func (mn *Maincore) GetEngagementPublicPostState(pptxhash utility.Hash,pptxindex uint32,claimaddress utility.Hash) (uint32,*utility.Hash,uint32) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementPublicPost)
	tmpkeybw.PutHash(pptxhash)
	tmpkeybw.PutUint32(pptxindex)
	tmpkeybw.PutHash(claimaddress)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetTxOutputState - txhash %x - index %d : %v",pptxhash,pptxindex, err)
		return StateValueIdentifierNotFound,nil,0
	}
	tmpbr:=utility.NewBufferReader(data)

	engagementidentifier:=tmpbr.GetUint32()
	etxhash:=tmpbr.GetHash()
	etxindex:=tmpbr.GetUint32()
	return engagementidentifier,&etxhash,etxindex
}*/


/*
func (mn *Maincore) PutEngagementPublicPostRewardState(publicposttxhash utility.Hash, publicposttxindex uint32,liketotalstake uint64,disliketotalstake uint64,liketotalweight big.Int,disliketotalweight big.Int) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementPublicPostReward)
	//tmpkeybw.PutRegistredNameKey(name)
	tmpkeybw.PutHash(publicposttxhash)
	tmpkeybw.PutUint32(publicposttxindex)

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierEngagementPublicPostReward)
	tmpbw.PutUint64(liketotalstake)
	tmpbw.PutUint64(disliketotalstake)
	tmpbw.PutBigInt(&liketotalweight)
	tmpbw.PutBigInt(&disliketotalweight)
	
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)

	return err
}

func (mn *Maincore) GetEngagementPublicPostRewardState(publicposttxhash utility.Hash, publicposttxindex uint32) (uint32,uint64,uint64,big.Int,big.Int) {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagementPublicPostReward)
	//tmpkeybw.PutRegistredNameKey(name)
	tmpkeybw.PutHash(publicposttxhash)
	tmpkeybw.PutUint32(publicposttxindex)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return 0,0,0,*big.NewInt(0),*big.NewInt(0)
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierEngagementPublicPostReward{
		//return nil,0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
		return 0,0,0,*big.NewInt(0),*big.NewInt(0)
	}

	//
	liketotalstake:=tmpbr.GetUint64()
	disliketotalstake:=tmpbr.GetUint64()
	liketotalweight:=tmpbr.GetBigInt()
	disliketotalweight:=tmpbr.GetBigInt()
	
	return stateidentifier,liketotalstake,disliketotalstake,*liketotalweight,*disliketotalweight // 
}




//engagementreward,err:=mn.GetEngagementPublicPostRewardValue(publicposttxhash,publicposttxindex,engagementid,engagementtxout.Value,height)
//
func (mn *Maincore) GetEngagementPublicPostRewardValue(publicposttxhash utility.Hash,publicposttxindex uint32,publicpostheight uint32,
													engagementid uint32,engagementstake uint64,stakeheight uint32) (uint64,error){

	var claimreward uint64
	var totalweight big.Int
	applog.Trace("GetEngagementPublicPostRewardValue hash %x index %d",publicposttxhash,publicposttxindex)
	stateidentifier,liketotalstake,disliketotalstake,liketotalweight,disliketotalweight:=mn.GetEngagementPublicPostRewardState(publicposttxhash,publicposttxindex)
	applog.Trace("stateidentifier %d,liketotalstake %d,disliketotalstake %d,liketotalweight %s,disliketotalweight %s",stateidentifier,liketotalstake,disliketotalstake,liketotalweight.String(),disliketotalweight.String())
	//_=stateidentifier
	if stateidentifier!=StateValueIdentifierEngagementPublicPostReward{
		return 0,fmt.Errorf("Found an incorrect stateidentifier associated with engagement public post reward state - identifier found %d",stateidentifier)
	}
	weight:=uint64(engagementstake)*uint64(stakeheight)
	//claimreward=totalreward*weight/totalweight
	if liketotalstake>=disliketotalstake{
		totalweight=liketotalweight
		if engagementid==utility.EngagementIdentifierDislikePublicPost {
			return 0,nil
		}
	} else {
		totalweight=disliketotalweight
		if engagementid==utility.EngagementIdentifierLikePublicPost {
			return 0,nil
		}
	} 
	bigtmp:=new(big.Int)
	bigtmp.SetUint64(liketotalstake+disliketotalstake)
 	//bigtmp:=big.NewInt(int64(liketotalstake+disliketotalstake))
	bigweight:=new(big.Int)
	bigweight.SetUint64(weight)
	bigtmp2:=bigtmp.Mul(bigtmp,bigweight)
	bigclaimreward:=bigtmp2.Div( bigtmp2 , &totalweight)
	claimreward=uint64(bigclaimreward.Int64())
	return claimreward,nil
}
*/
/*
func (mn *Maincore) GetNameState(name []byte) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	return stateidentifier
}
*/
////////////////////////////////////////////

func  (mn *Maincore)  RebuildMainstate() {
	/*
	for l:=0;l<mn.dataf.NbChunks();l++{
		data:=mn.dataf.GetChunkById(int(l))
		hash:=utility.ComputeHash(data)
		mn.PutPublicPostState(hash,[]byte(""),*utility.NewHash(nil),0,uint32(l))
	}
	*/

	//
	for i:=0;i<mn.mainbf.NbChunks();i++{
		//applog.Trace("block %d %d ",i,len(mn.GetMainblock(i).Transactions))
		if !mn.ValidateMainblockTransactions(uint32 (i), &mn.GetMainblock(i).Transactions){
			applog.Warning("Rejected block - incorrect transactions - block height %d",i)
			applog.Warning("RebuildMainstate Failed")
			os.Exit(0)
			//return false
		}
		for j:=0;j<len(mn.GetMainblock(i).Transactions);j++{
			//applog.Trace("block %d %x ",i,mn.GetMainblock(i).Transactions[j].ComputeHash())
			tx:=mn.GetMainblock(i).Transactions[j]
			txhash:=tx.ComputeHash()
			mn.PutTxState(txhash,uint32(i),uint32(j))
			mn.UpdateMainstate(tx,uint32(i))
			//
		}
		//
	}
}


/*
func (mn *Maincore) PutUnconfirmedTxOutputState(txhash utility.Hash,index uint32,stateidentifier uint32) error{

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutBytes([]byte("UNCONFIRMED"))
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(stateidentifier)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)

	//applog.Trace("%x %x  ",txhash,height,number)
	return err
}
func (mn *Maincore) GetUnconfirmedTxOutputState(txhash utility.Hash,index uint32) (uint32){
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutBytes([]byte("UNCONFIRMED"))
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		if err!=leveldberrors.ErrNotFound {
			applog.Trace("GetUnconfirmedTxOutputState - txhash %x - index %d : %v",txhash,index, err)
		}
		
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()

	applog.Trace("unconfirmed txoutput stateidentifier %d ",stateidentifier)
	return stateidentifier
}
func (mn *Maincore) DeleteAllUnconfirmedTxOutputState() {
	iter := mn.mainstatedb.NewIterator(leveldbutil.BytesPrefix([]byte("UNCONFIRMED")), nil)
	for iter.Next() {
			// contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			key := iter.Key()
			value := iter.Value()
			applog.Trace("unconfirmed txoutput hash %x index %d  ",key,value)
			err := mn.mainstatedb.Delete(key, nil)
			if err!=nil{
				applog.Trace("error",err)
			}

	}
	iter.Release()
	err:= iter.Error()
	if err!=nil{
		applog.Trace("error",err)
	}
	
}
*/