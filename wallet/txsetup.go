package wallet

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/utility"
)

func (wlt *Wallet) ComputeBalance() uint64 {
	var totalbalance uint64
	nbunspent := 0
	for i := 0; i < len(wlt.Assetarray); i++ {

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		if wlt.Assetarray[i].StateString == "UNSPENT" {
			totalbalance += wlt.Assetarray[i].Value
			nbunspent++
		} else {
			fmt.Printf("Asset state %s %d \n", wlt.Assetarray[i].StateString, wlt.Assetarray[i].Value)
		}

	}
	fmt.Printf("Total Assets %d - nbunspent %d\n", len(wlt.Assetarray), nbunspent)
	return totalbalance
}
func (wlt *Wallet) SetupTransactionForNameUnregistration(name string, fee uint64) (*utility.Transaction, error) {
	a, aerr := wlt.GetAssetFromRegisteredName(name)
	if aerr != nil {
		return nil, aerr
	}
	hash := a.Hash
	index := a.Index
	prvkeyindex := a.Privatekeyindex
	amount := a.Value //TODO if amount less than fee add automatically new assets in the input

	pubkey := wlt.Privatekeyarray[prvkeyindex].PubKey().SerializeCompressed()
	tmptxin := utility.NewECDSANameUnregistration(hash, index, pubkey)
	//tx,err:=wlt.SetupTransactionAmount(amount,fee,&tmptxin,nil)
	tx := new(utility.Transaction)
	tx.Version = 1
	tx.Vin = append(tx.Vin, tmptxin)
	tmptxout := utility.NewECDSATxOut(amount-fee, wlt.GenerateKeyPair())
	tx.Vout = append(tx.Vout, tmptxout)
	newtxhash, herr := tx.ComputeSigningHash() //
	if herr != nil {
		return nil, herr
	}

	//i:=len(tx.Vin)-1
	i := 0
	/////sig:= ecdsa.Sign(wlt.Privatekeyarray[prvkeyindex],newtxhash[:])
	sigBytes := utility.Sign(wlt.Privatekeyarray[prvkeyindex], newtxhash[:])
	//if err != nil {
	//	applog.Trace("error: unable to sign input %d of transaction",i)
	//	return nil,fmt.Errorf("error: unable to sign input %d of transaction",i)
	//}
	//applog.Trace("signature %x", sig)
	tmpbw := utility.NewBufferWriter()
	//tmpbw.PutVarUint(uint64(len(sig)))
	/////tmpbw.PutBytes(sig.Serialize())
	tmpbw.PutBytes(sigBytes)
	//tx.Vin[i].Bytecode=append(tx.Vin[i].Bytecode, tmpbw.GetContent() ...)
	tx.Vin[i].Signature = tmpbw.GetContent() //selectedassetarray[i].Privatekeyindex

	return tx, nil
}

func (wlt *Wallet) SetupTransactionForRegistredNameCommit(name string, commitbytes []byte, fee uint64) (*utility.Transaction, error) {
	a, aerr := wlt.GetAssetFromRegisteredName(name)
	if aerr != nil {
		return nil, aerr
	}
	hash := a.Hash
	index := a.Index
	prvkeyindex := a.Privatekeyindex
	//
	pubkey := wlt.Privatekeyarray[prvkeyindex].PubKey().SerializeCompressed()

	tmptxin := utility.NewECDSARegistredNameCommit(hash, index, pubkey, commitbytes)
	tx, err := wlt.SetupTransactionAmount(0, fee, &tmptxin, nil)
	//tx:=new(utility.Transaction)
	//tx.Version=1
	//tx.Vin=append(tx.Vin, tmptxin )
	//tmptxout:=utility.NewECDSATxOut(0,wlt.GenerateKeyPair())
	//tx.Vout=append(tx.Vout,tmptxout)
	applog.Trace("------------txin %x", tmptxin)
	newtxhash, herr := tx.ComputeSigningHash() //
	if herr != nil {
		return nil, herr
	}

	i := len(tx.Vin) - 1
	//i:=0
	sig, err := ecdsa.SignCompact(wlt.Privatekeyarray[prvkeyindex], newtxhash[:], true)
	if err != nil {
		applog.Trace("error: unable to sign input %d of transaction", i)
		return nil, fmt.Errorf("error: unable to sign input %d of transaction", i)
	}
	//applog.Trace("signature %x", sig)
	tmpbw := utility.NewBufferWriter()
	//tmpbw.PutVarUint(uint64(len(sig)))
	tmpbw.PutBytes(sig)
	//tx.Vin[i].Bytecode=append(tx.Vin[i].Bytecode, tmpbw.GetContent() ...)
	tx.Vin[i].Signature = tmpbw.GetContent() //selectedassetarray[i].Privatekeyindex

	return tx, err
}

/*
func (wlt * Wallet) SetupTransactionForEngagementPublicPost(eid uint32,txhash utility.Hash,index uint32,claimaddress utility.Hash,stakedamount uint64, fee uint64)  (*utility.Transaction,error){
	tmptxout:=utility.NewECDSAEngagementPublicPost(eid,stakedamount,txhash,index,claimaddress)
	tx,err:=wlt.SetupTransactionAmount(stakedamount,fee,nil,&tmptxout)
	return tx,err
}


func (wlt * Wallet) SetupTransactionForEngagementRewardClaim(etxhash utility.Hash,etxindex uint32,claimaddress utility.Hash,claimedamount uint64,fee uint64)  (*utility.Transaction,error){


	prvkeyindex:=wlt.GetPrivatekeyindexFromAddress(claimaddress)
	pubkey:=wlt.Privatekeyarray[prvkeyindex].PubKey().SerializeCompressed()
	tmptxin:=utility.NewECDSAEngagementRewardClaim(etxhash,etxindex,pubkey)

	//rewardamount,err:=GetEngagementPublicPostRewardValue(publicposttxhash,publicposttxindex,publicpostheight,engagementid,engagementstake,height)

	tx:=new(utility.Transaction)
	tx.Version=1
	tx.Vin=append(tx.Vin, tmptxin )
	tmptxout:=utility.NewECDSATxOut(claimedamount-fee,wlt.GenerateKeyPair())
	tx.Vout=append(tx.Vout,tmptxout)

	applog.Trace("------------txin %x",tmptxin)
	newtxhash,herr:=tx.ComputeSigningHash()//
	if herr!=nil{
		return nil,herr
	}

	i:=len(tx.Vin)-1
	//i:=0
		sig, err := ecdsa.SignCompact(wlt.Privatekeyarray[prvkeyindex],newtxhash[:],true)
		if err != nil {
			applog.Trace("error: unable to sign input %d of transaction",i)
			return nil,fmt.Errorf("error: unable to sign input %d of transaction",i)
		}
		//applog.Trace("signature %x", sig)
		tmpbw:=utility.NewBufferWriter()
		//tmpbw.PutVarUint(uint64(len(sig)))
		tmpbw.PutBytes(sig)
		//tx.Vin[i].Bytecode=append(tx.Vin[i].Bytecode, tmpbw.GetContent() ...)
		tx.Vin[i].Signature=tmpbw.GetContent()//selectedassetarray[i].Privatekeyindex


	return tx,err
}
*/

func (wlt *Wallet) SetupTransactionForNameRegistration(name []byte, pubkeyhash utility.Hash, amount uint64, fee uint64) (*utility.Transaction, error) {
	CommPubkey := wlt.GenerateCommKey(name)
	tmptxout := utility.NewECDSANameRegistration(amount, name, pubkeyhash, utility.CommKeyIdentifierRSAOAEP, CommPubkey)
	tx, err := wlt.SetupTransactionAmount(amount, fee, nil, &tmptxout)
	return tx, err
}
func (wlt *Wallet) SetupTransactionToPublicKeyHash(pubkeyhash utility.Hash, amount uint64, fee uint64) (*utility.Transaction, error) {

	tmptxout := utility.NewECDSATxOut(amount, pubkeyhash)
	tx, err := wlt.SetupTransactionAmount(amount, fee, nil, &tmptxout)
	return tx, err
}

func (wlt *Wallet) SetupTransactionAmount(amount uint64, fee uint64, txin *utility.TxIn, txout *utility.TxOut) (*utility.Transaction, error) {
	tx := new(utility.Transaction)
	tx.Version = 1
	/*if amount==0{
		applog.Trace("error: can not setup transaction with no ammount")
		return nil,fmt.Errorf("error: can not setup transaction with no ammount")
	}*/

	var selectedassetarray []Asset
	var selectedassetindexarray []int
	var selectedamount uint64 = 0
	var i int = 0
	for {
		//applog.Trace("selected amout %d i %d ",selectedamount,i)
		if i == len(wlt.Assetarray) {
			applog.Trace("error: insufficient funds to setup transaction")
			return nil, fmt.Errorf("error: insufficient funds to setup transaction")
		}

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		if wlt.Assetarray[i].StateString == "UNSPENT" {
			selectedassetarray = append(selectedassetarray, wlt.Assetarray[i])
			selectedamount += wlt.Assetarray[i].Value
			//wlt.Assetarray[i].StateString="BROADCASTED"
			selectedassetindexarray = append(selectedassetindexarray, i)
		}
		if selectedamount >= amount+fee {
			break
		}
		i++
	}
	for k := 0; k < len(selectedassetindexarray); k++ {
		wlt.Assetarray[selectedassetindexarray[k]].StateString = "BROADCASTED"
	}
	applog.Trace("selected amout %d", selectedamount)

	//tx.Vin=make([]utility.TxIn,len(selectedassetarray))
	for i := 0; i < len(selectedassetarray); i++ {

		pubkeycompressedbytes := wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex].PubKey().SerializeCompressed()
		if len(pubkeycompressedbytes) != 33 {
			applog.Trace("error: serialize compressed public key is not 33 length")
			return nil, fmt.Errorf("error: serialize compressed public key is not 33 length")
		}

		//tx.Vin[i].Signature=//selectedassetarray[i].Privatekeyindex
		tmptxin := utility.NewECDSATxIn(selectedassetarray[i].Hash, selectedassetarray[i].Index, pubkeycompressedbytes)
		tx.Vin = append(tx.Vin, tmptxin)

	}

	//tx.Vout=make([]utility.TxOut,2)

	tx.Vout = append(tx.Vout, utility.NewECDSATxOut(selectedamount-amount-fee, wlt.GenerateKeyPair()))
	//tx.Vout[0].Value=uint64(selectedamount-amount)
	//tx.Vout[0].Bytecode=wlt.GenerateKeyPair()
	if txin != nil {
		tx.Vin = append(tx.Vin, *txin)
		applog.Trace("////////////////")
	}
	if txout != nil {
		tx.Vout = append(tx.Vout, *txout)
	}

	newtxhash, err := tx.ComputeSigningHash() //
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(selectedassetarray); i++ {
		///sig:= ecdsa.Sign(wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex],newtxhash[:])
		sigBytes := utility.Sign(wlt.Privatekeyarray[selectedassetarray[i].Privatekeyindex], newtxhash[:])
		//if err != nil {
		//	applog.Trace("error: unable to sign input %d of transaction",i)
		//	return nil,fmt.Errorf("error: unable to sign input %d of transaction",i)
		//}
		//applog.Trace("signature %x", sig)
		tmpbw := utility.NewBufferWriter()
		//tmpbw.PutVarUint(uint64(len(sig)))
		////tmpbw.PutBytes(sig.Serialize())
		tmpbw.PutBytes(sigBytes)
		//tx.Vin[i].Bytecode=append(tx.Vin[i].Bytecode, tmpbw.GetContent() ...)
		tx.Vin[i].Signature = tmpbw.GetContent() //selectedassetarray[i].Privatekeyindex
	}

	return tx, nil
}

func (wlt *Wallet) GetUnconfirmedBroadcastedTxs() []*utility.Transaction {
	var tmptxs []*utility.Transaction
	for i := 0; i < len(wlt.Broadcastedtxarray); i++ {

		if wlt.Broadcastedtxarray[i].ConfirmationString != "CONFIRMED" {
			tmptxs = append(tmptxs, &wlt.Broadcastedtxarray[i].Tx)
		}

	}
	return tmptxs
}
