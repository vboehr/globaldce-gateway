package wallet
import (
	"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	"github.com/globaldce/globaldce-gateway/utility"
	//"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)



func (wlt * Wallet) SetupTransactionToPublicKeyHashArray(pubkeyhasharray []utility.Hash,amountarray []uint64,fee uint64) (*utility.Transaction,error){

	//tmptxout:=utility.NewECDSATxOut(amount,pubkeyhash)
	//tx,err:=wlt.SetupTransactionAmountArray(totalamount,fee,nil,&tmptxout)
	//return tx,err
//}
var amount uint64 =0
for _,partialamount:=range amountarray{
	amount+=partialamount
}

//func (wlt * Wallet) SetupTransactionAmountArray(amount uint64,fee uint64,txin *utility.TxIn,txout *utility.TxOut) (*utility.Transaction,error){
	tx:=new(utility.Transaction)
	tx.Version=1
	/*if amount==0{
		applog.Trace("error: can not setup transaction with no ammount")
		return nil,fmt.Errorf("error: can not setup transaction with no ammount")
	}*/

	var selectedassetarray [] Asset
	var selectedassetindexarray [] int
	var selectedamount uint64 =0
	var i int=0
	for {
		//applog.Trace("selected amout %d i %d ",selectedamount,i)
		if i==len(wlt.Assetarray){
			applog.Trace("error: insufficient funds to setup transaction")
			return nil,fmt.Errorf("error: insufficient funds to setup transaction")
		}

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		if wlt.Assetarray[i].StateString=="UNSPENT"{
			selectedassetarray=append(selectedassetarray,wlt.Assetarray[i])
			selectedamount+=wlt.Assetarray[i].Value
			//wlt.Assetarray[i].StateString="BROADCASTED"
			selectedassetindexarray=append(selectedassetindexarray,i)
		}
		if selectedamount>=amount+fee{
			break
		}
		i++
	}
	for k:=0;k<len(selectedassetindexarray);k++{
		wlt.Assetarray[selectedassetindexarray[k]].StateString="BROADCASTED"
	}
	applog.Trace("selected amout %d",selectedamount)

	//tx.Vin=make([]utility.TxIn,len(selectedassetarray))
	for i:=0;i<len(selectedassetarray);i++{

		pubkeycompressedbytes:=wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex].PubKey().SerializeCompressed()
		if len(pubkeycompressedbytes)!=33{
			applog.Trace("error: serialize compressed public key is not 33 length")
			return nil,fmt.Errorf("error: serialize compressed public key is not 33 length")
		}

		
		//tx.Vin[i].Signature=//selectedassetarray[i].Privatekeyindex 
		tmptxin:=utility.NewECDSATxIn(selectedassetarray[i].Hash,selectedassetarray[i].Index,pubkeycompressedbytes)
		tx.Vin=append(tx.Vin,tmptxin)
		
	}



	//tx.Vout=make([]utility.TxOut,2)

	tx.Vout=append(tx.Vout,utility.NewECDSATxOut(selectedamount-amount-fee,wlt.GenerateKeyPair()))
	//tx.Vout[0].Value=uint64(selectedamount-amount) 
	//tx.Vout[0].Bytecode=wlt.GenerateKeyPair()
	/*if txin!=nil{
		tx.Vin=append(tx.Vin,*txin)
		applog.Trace("////////////////")
	}
	*/
	//

	//if txout!=nil{
	for i,pubkeyhash:=range pubkeyhasharray{
		tmptxout:=utility.NewECDSATxOut(amountarray[i],pubkeyhash)
		tx.Vout=append(tx.Vout,tmptxout) 
	}

	newtxhash,err:=tx.ComputeSigningHash()//
	if err!=nil{
		return nil,err
	}

	for i:=0;i<len(selectedassetarray);i++{
		////sig := ecdsa.Sign(wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex],newtxhash[:])
		sigBytes:=utility.Sign(wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex],newtxhash[:])
		//if err != nil {
		//	applog.Trace("error: unable to sign input %d of transaction",i)
		//	return nil,fmt.Errorf("error: unable to sign input %d of transaction",i)
		//}
		//applog.Trace("signature %x", sig)
		tmpbw:=utility.NewBufferWriter()
		//tmpbw.PutVarUint(uint64(len(sig)))
		////tmpbw.PutBytes(sig.Serialize())
		tmpbw.PutBytes(sigBytes)
		//tx.Vin[i].Bytecode=append(tx.Vin[i].Bytecode, tmpbw.GetContent() ...)
		tx.Vin[i].Signature=tmpbw.GetContent()//selectedassetarray[i].Privatekeyindex 
	}

	return tx,nil
}
