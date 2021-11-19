package mainchain

import (
	"github.com/globaldce/globaldce-toolbox/applog"

	"github.com/globaldce/globaldce-toolbox/utility"
	"fmt"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)

const (
	StateIdentifierNotFound=0
	StateIdentifierUnspentTxOutput=1
	StateIdentifierSpentTxOutput=2
	StateIdentifierTx=3
	StateIdentifierActifNameRegistration=4
	StateIdentifierInactifNameRegistration=5
	StateIdentifierData=6
	StateIdentifierDataFile=7
)



func (mn *Maincore) PutTxState(txhash utility.Hash,height uint32,number uint32) error{
	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateIdentifierTx)// transaction type
	tmpbw.PutUint32(height)
	tmpbw.PutUint32(number)
	err := mn.mainstatedb.Put(txhash[:], tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetTxState(txhash utility.Hash) (uint32,uint32,uint32) {
	data, err := mn.mainstatedb.Get(txhash[:], nil)
	if err != nil {
		applog.Trace("GetTxState - txhash %x : %v",txhash, err)
		return StateIdentifierNotFound,0,0
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
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(stateidentifier)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetTxOutputState(txhash utility.Hash,index uint32) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutBytes(txhash[:])
	tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		applog.Trace("GetTxOutputState - txhash %x - index %d : %v",txhash,index, err)
		return StateIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	return stateidentifier
}


func (mn *Maincore) PutNameState(name []byte,stateidentifier uint32) error {

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	tmpbw:=utility.NewBufferWriter()

	tmpbw.PutUint32(stateidentifier)

	err := mn.mainstatedb.Put(tmpkeybw.GetContent(), tmpbw.GetContent(), nil)
	return err
}

func (mn *Maincore) GetNameState(name []byte) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return StateIdentifierNotFound
	}
	tmpbr:=utility.NewBufferReader(data)

	stateidentifier:=tmpbr.GetUint32()
	return stateidentifier
}

func (mn *Maincore) PutDataState(datahash utility.Hash,id uint32) error{
	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(StateIdentifierData)// transaction type
	tmpbw.PutUint32(id)
	err := mn.mainstatedb.Put(datahash[:], tmpbw.GetContent(), nil)
	return err
}
func (mn *Maincore) GetDataState(datahash utility.Hash)( uint32, error){
	

	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutHash(datahash)
	//tmpkeybw.PutUint32(index)

	valuebytes, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return 0,err
	}
	tmpbr:=utility.NewBufferReader(valuebytes)
	stateidentifier:=tmpbr.GetUint32()
	if stateidentifier!=StateIdentifierData{
		return 0,fmt.Errorf("Found an incorrect stateidentifier associated with data hash %x - identifier found")
	}
	id:=tmpbr.GetUint32()
	return id,nil
}
/*
func (mn *Maincore) GetNameState(name []byte) uint32 {
	tmpkeybw:=utility.NewBufferWriter()
	tmpkeybw.PutRegistredNameKey(name)
	//tmpkeybw.PutUint32(index)

	data, err := mn.mainstatedb.Get(tmpkeybw.GetContent(), nil)
	if err != nil {
		//applog.Trace("GetNameState - txhash %x - index %d : %v",txhash,index, err)
		return StateIdentifierNotFound
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
				mn.PutTxOutputState(txhash,uint32(k),StateIdentifierUnspentTxOutput)
				applog.Trace("Puttting %x %d  stat %d",txhash,uint32(k),StateIdentifierUnspentTxOutput)
			case utility.ModuleIdentifierECDSANameRegistration:
				_,name,_,_:=utility.DecodeECDSANameRegistration(tx.Vout[k].Bytecode) 
				mn.PutTxOutputState(txhash,uint32(k),StateIdentifierActifNameRegistration)
				applog.Trace("Puttting %x %d %d",txhash,uint32(k),StateIdentifierActifNameRegistration)
				mn.PutNameState(name,StateIdentifierActifNameRegistration)
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
				mn.PutTxOutputState(tx.Vin[l].Hash,tx.Vin[l].Index,StateIdentifierSpentTxOutput)
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
		
		return StateIdentifierNotFound
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