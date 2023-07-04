package utility

import (
	"encoding/json"
	//"encoding/binary"
	"fmt"
)

func (tx *Transaction) JSONSerialize() []byte {

	txbytes, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("Mainheader serialize error:", err)
	}
	return txbytes
}
func (tx *Transaction) Serialize() []byte {
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(uint32(tx.Version))

	tmpbw.PutVarUint(uint64(len(tx.Vin)))
	for i := 0; i < len(tx.Vin); i++ {
		tmpbw.PutBytes(tx.Vin[i].Hash[:])
		tmpbw.PutUint32(tx.Vin[i].Index)
		tmpbw.PutVarUint(uint64(len(tx.Vin[i].Bytecode)))
		tmpbw.PutBytes(tx.Vin[i].Bytecode)
		tmpbw.PutVarUint(uint64(len(tx.Vin[i].Signature)))
		tmpbw.PutBytes(tx.Vin[i].Signature)
	}

	tmpbw.PutVarUint(uint64(len(tx.Vout)))
	for j := 0; j < len(tx.Vout); j++ {
		tmpbw.PutUint64(uint64(tx.Vout[j].Value))
		//tmpbw.PutBytes(tx.Vout[j].Address[:])
		tmpbw.PutVarUint(uint64(len(tx.Vout[j].Bytecode)))
		tmpbw.PutBytes(tx.Vout[j].Bytecode)
	}

	return tmpbw.GetContent()

}
func UnserializeTransaction(bytes []byte) (*Transaction, error) {
	/*
		tmptx:=new(Transaction)
		err:=json.Unmarshal(bytes,tmptx)
		if err != nil {
			fmt.Println("error:", err)
			return nil,err
		}
		return tmptx,nil
	*/
	tmpbr := NewBufferReader(bytes)

	tx := new(Transaction)
	tx.Version = int32(tmpbr.GetUint32())
	lenvin := tmpbr.GetVarUint()

	for i := uint64(0); i < lenvin; i++ {
		var tmpin TxIn
		tmpin.Hash = *NewHash(tmpbr.GetBytes(32))
		tmpin.Index = tmpbr.GetUint32()
		bytecodelen := tmpbr.GetVarUint()
		//primitivemoduleid:=tmpbr.GetUint32()
		//if primitivemoduleid==ModuleIdentifierECDSATxIn {
		//pubkeylen:=uint(tmpbr.GetUint32())
		//tmpin.Pubkeycompressed=tmpbr.GetBytes(pubkeylen)
		//siglen:=uint(tmpbr.GetUint32())
		//tmpin.Signature=tmpbr.GetBytes(siglen)
		//!!	tmpbw.PutBytes(tx.Vin[i].Hash[:])
		//!!	tmpbw.PutUint32(tx.Vin[i].Index)
		//!!	tmpbw.PutBytes(tx.Vin[i].Pubkeycompressed)
		//!!	tmpbw.PutBytes(tx.Vin[i].Signature)
		tmpin.Bytecode = tmpbr.GetBytes(uint(bytecodelen))
		signaturelen := tmpbr.GetVarUint()
		tmpin.Signature = tmpbr.GetBytes(uint(signaturelen))
		tx.Vin = append(tx.Vin, tmpin)
		//}

	}
	lenvout := tmpbr.GetVarUint()

	//for j:=0;j<len(tx.Vout);j++{
	//	tmpbw.PutUint64(uint64(tx.Vout[j].Value))
	//	tmpbw.PutBytes(tx.Vout[j].Pubkeyhash[:])
	//}
	for j := uint64(0); j < lenvout; j++ {
		var tmpout TxOut

		tmpout.Value = uint64(tmpbr.GetUint64())

		//tmpout.Address=*NewHash(tmpbr.GetBytes(32))
		bytecodelen := tmpbr.GetVarUint()
		//if primitivemoduleid!=ModuleIdentifierECDSATxOut{
		//	return nil,fmt.Errorf("output %d serialization error - unsupported output type",j)
		//}
		tmpout.Bytecode = tmpbr.GetBytes(uint(bytecodelen))
		tx.Vout = append(tx.Vout, tmpout)
	}

	//applog.Trace("result %v",tx)
	tmpbrerr := tmpbr.GetError()
	if tmpbrerr != nil {
		return nil, tmpbrerr
	}
	return tx, nil
}
