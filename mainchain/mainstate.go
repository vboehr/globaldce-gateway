package mainchain

import (
	"github.com/globaldce/globaldce-toolbox/applog"

	"github.com/globaldce/globaldce-toolbox/utility"
	"fmt"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)

const (
	StateValueIdentifierNotFound=0
	StateValueIdentifierUnspentTxOutput=1
	StateValueIdentifierSpentTxOutput=2
	StateValueIdentifierTx=3
	StateValueIdentifierActifNameRegistration=4
	StateValueIdentifierInactifNameRegistration=5
	StateValueIdentifierData=6
	StateValueIdentifierDataFile=7
	StateValueIdentifierEngagement=8
)

const (
	StateKeyIdentifierTxOutput=1
	StateKeyIdentifierTx=3
	StateKeyIdentifierNameRegistration=4
	StateKeyIdentifierData=6
	StateKeyIdentifierDataFile=7
	StateKeyIdentifierEngagement=8
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

func (mn *Maincore) PutPublicPostState(datahash utility.Hash,namebytes []byte,id uint32) error{
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierData)
	tmpkeybw.PutHash(datahash)
	//tmpkeybw.GetContent()

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierData)// 
	tmpbw.PutVarUint(uint64(len(namebytes)))
	tmpbw.PutBytes(namebytes)
	tmpbw.PutUint32(id)
	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}
func (mn *Maincore) GetPublicPostState(datahash utility.Hash)([]byte, uint32, error){
	

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierData)
	tmpkeybw.PutHash(datahash)
	//tmpkeybw.PutUint32(index)

	valuebytes, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return nil,0,err
	}
	tmpbr:=utility.NewBufferReader(valuebytes)
	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierData{
		return nil,0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
	}
	namebyteslen:=tmpbr.GetVarUint()
	namebytes:=tmpbr.GetBytes(uint(namebyteslen))
	id:=tmpbr.GetUint32()
	return namebytes,id,nil
}




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
func (mn *Maincore) PutEngagementState(name []byte,engagementidentifier uint32,value uint32) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagement)
	tmpkeybw.PutRegistredNameKey(name)

	tmpkeybw.PutUint32(engagementidentifier)

	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateValueIdentifierEngagement)
	tmpbw.PutUint32(value)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)

	return err
}

func (mn *Maincore) GetEngagementState(name []byte,engagementidentifier uint32) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutUint32(StateKeyIdentifierEngagement)
	tmpkeybw.PutRegistredNameKey(name)
	tmpkeybw.PutUint32(engagementidentifier)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return StateValueIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateValueIdentifierEngagement{
		//return nil,0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
		return 0
	}
	value:=tmpbr.GetUint32()
	return value
}

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
}*/
////////////////////////////////////////////

func  (mn *Maincore)  RebuildMainstate() {
	for i:=0;i<mn.mainbf.NbChunks();i++{
		//applog.Trace("block %d %d ",i,len(mn.GetMainblock(i).Transactions))
		for j:=0;j<len(mn.GetMainblock(i).Transactions);j++{
			//applog.Trace("block %d %x ",i,mn.GetMainblock(i).Transactions[j].ComputeHash())
			tx:=mn.GetMainblock(i).Transactions[j]
			txhash:=tx.ComputeHash()
			mn.PutTxState(txhash,uint32(i),uint32(j))
			mn.UpdateMainstate(tx)
			//
		}
		//
	}
}

func  (mn *Maincore)  UpdateMainstate(tx utility.Transaction) {
	txhash:=tx.ComputeHash()

	for k:=0;k<len(tx.Vout);k++{
		moduleid:=utility.DecodeBytecodeId(tx.Vout[k].Bytecode)
		switch moduleid {
			case utility.ModuleIdentifierECDSATxOut:
				mn.PutTxOutputState(txhash,uint32(k),StateValueIdentifierUnspentTxOutput)
				applog.Trace("Puttting %x %d  stat %d",txhash,uint32(k),StateValueIdentifierUnspentTxOutput)
			case utility.ModuleIdentifierECDSANameRegistration:
				_,name,_,_:=utility.DecodeECDSANameRegistration(tx.Vout[k].Bytecode) 
				mn.PutTxOutputState(txhash,uint32(k),StateValueIdentifierActifNameRegistration)
				applog.Trace("Puttting %x %d %d",txhash,uint32(k),StateValueIdentifierActifNameRegistration)
				mn.PutNameState(name,StateValueIdentifierActifNameRegistration)
			case utility.ModuleIdentifierEngagement:
				eid,name,_,_:=utility.DecodeEngagement(tx.Vout[k].Bytecode)
				if eid==utility.EngagementIdentifierLikeName {
					mn.AddEngagementLikeName(name)
				}
				if eid==utility.EngagementIdentifierDislikeName {
					mn.AddEngagementDislikeName(name)
				}
			//default:
			//
			//	ModuleIdentifierECDSANameRegistration=3
			//	ModuleIdentifierECDSANameUnregistration=4
		}
		
	}
	//
	for l:=0;l<len(tx.Vin);l++{
		//
		moduleid:=utility.DecodeBytecodeId(tx.Vin[l].Bytecode)
		switch moduleid {
			case utility.ModuleIdentifierECDSATxIn:
				mn.PutTxOutputState(tx.Vin[l].Hash,tx.Vin[l].Index,StateValueIdentifierSpentTxOutput)
			case utility.ModuleIdentifierECDSANamePublicPost:
				_,height,number:=mn.GetTxState(tx.Vin[l].Hash)
				//applog.Trace("height%d,number%d",height,number)
				tmpinpututxo:=mn.GetMainblock(int(height)).Transactions[number].Vout[tx.Vin[l].Index]
				_,name,_,_:=utility.DecodeECDSANameRegistration(tmpinpututxo.Bytecode)
			
				_,ed,_:=utility.DecodeECDSANamePublicPost(tx.Vin[l].Bytecode)
				
				_,_,err:=mn.GetPublicPostState(ed.Hash)//([]byte, uint32, error){
				
				if (err==nil)&&(mn.IsBannedName(name)){
					mn.PutPublicPostState(ed.Hash,name,uint32(0))
					mn.AddToMissingDataHashArray(ed.Hash)	
				}

			//default:
		}
		
		
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