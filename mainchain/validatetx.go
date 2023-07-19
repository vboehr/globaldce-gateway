package mainchain

import (
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/utility"
	//"math/big"
	//"math"
	"fmt"
	//"os"
)

func (mn *Maincore) ValidateTransaction(tx *utility.Transaction) (bool, uint64) {
	//if !tx.VerifySignatures(){
	//	applog.Trace("invalid signatures")
	//	return false,0
	//}
	//
	txsigninghash, err := tx.ComputeSigningHash() //
	if err != nil {
		return false, 0
	}
	var totalinputamount uint64 = 0
	var totaloutputamount uint64 = 0
	applog.Trace("nb inputs %d", len(tx.Vin))
	for i := 0; i < len(tx.Vin); i++ {
		//applog.Trace("ModuleId %d",utility.DecodeBytecodeId(tx.Vin[i].Bytecode))
		applog.Trace("Validating txin %d", i)
		txinvalue, err := mn.ValidateTxIn(txsigninghash, tx.Vin[i])
		if err != nil {
			//TODO display warning
			applog.Trace("%v", err)
			return false, 0
		}
		totalinputamount += txinvalue
	}
	for j := 0; j < len(tx.Vout); j++ {
		txoutvalue, err := mn.ValidateTxOut(tx.Vout[j])
		if err != nil {
			//TODO display warning
			applog.Trace("%v", err)
			return false, 0
		}
		totaloutputamount += txoutvalue
		//totaloutputamount+=int64(tx.Vout[j].Value)
	}
	if totalinputamount < totaloutputamount {
		return false, 0
	}
	feeamount := totalinputamount - totaloutputamount
	applog.Trace("inputamount %d outputamount %d", totalinputamount, totaloutputamount)

	return true, feeamount
}

func (mn *Maincore) ValidateTxIn(signinghash utility.Hash, tmptxin utility.TxIn) (uint64, error) {
	moduleid := utility.DecodeBytecodeId(tmptxin.Bytecode)
	applog.Trace("moduleid %d", moduleid)
	switch moduleid {
	case utility.ModuleIdentifierECDSATxIn:
		if mn.GetTxOutputState(tmptxin.Hash, tmptxin.Index) != StateValueIdentifierUnspentTxOutput {
			return 0, fmt.Errorf("invalide GetTxOutputState %d for %x %d ", mn.GetTxOutputState(tmptxin.Hash, tmptxin.Index), tmptxin.Hash, tmptxin.Index)
		}
		_, height, number := mn.GetTxState(tmptxin.Hash)
		inpututxo := mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]
		pubkeycompressedbytes, _, err := utility.DecodeECDSATxInBytecode(tmptxin.Bytecode)
		if err != nil {
			return 0, fmt.Errorf("DecodeECDSATxInBytecode Error %v", err)
		}
		verr := utility.VerifySignature(signinghash, tmptxin.Signature, pubkeycompressedbytes)
		if verr != nil {
			return 0, verr
		}
		if !(inpututxo.CompareWithAddress(utility.ComputeHash(pubkeycompressedbytes))) {
			//if (utility.ComputeHash (pubkeycompressedbytes)!=inpututxo.Address){//TODO
			return 0, fmt.Errorf("Corrupt transaction input - input public key do not match its associated output hash")
		}
		//applog.Trace("Value",inpututxo.Value)
		return inpututxo.Value, nil
	case utility.ModuleIdentifierECDSANameUnregistration:
		if mn.GetTxOutputState(tmptxin.Hash, tmptxin.Index) != StateValueIdentifierActifNameRegistration {
			return 0, fmt.Errorf("invalide ** GetTxOutputState %d for %x %d ", mn.GetTxOutputState(tmptxin.Hash, tmptxin.Index), tmptxin.Hash, tmptxin.Index)
		}
		_, height, number := mn.GetTxState(tmptxin.Hash) //URGENT TODO error handling
		////////////////////////////////////////////
		//if (mn.GetConfirmedMainchainLength()-height)>500000{
		//	return 0,fmt.Errorf("Unregistration is too soon")
		//}
		////////////////////////////////////////////
		inpututxo := mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]       //URGENT TODO error handling
		pubkeycompressedbytes, _, err := utility.DecodeECDSANameUnregistration(tmptxin.Bytecode) //DecodeECDSANameUnregistration(bytecode []byte) ([]byte,*Extradata,error){
		if err != nil {
			return 0, fmt.Errorf("DecodeECDSANameUnregistration Error %v", err)
		}
		verr := utility.VerifySignature(signinghash, tmptxin.Signature, pubkeycompressedbytes)
		if verr != nil {
			return 0, verr
		}
		if !(inpututxo.CompareWithAddress(utility.ComputeHash(pubkeycompressedbytes))) {
			//if (utility.ComputeHash (pubkeycompressedbytes)!=inpututxo.Address){//TODO
			return 0, fmt.Errorf("Corrupt transaction input - input public key do not match its associated output hash")
		}
		//applog.Trace("Value",inpututxo.Value)
		//fmt.Println("UNregistration **********")
		//os.Exit(0)
		/*
			_,name,_,_,_:=utility.DecodeECDSANameRegistration(inpututxo.Bytecode)
			_,totalstakedislike:=mn.GetEngagementDislikeName(name)
			_,totalstakelike:=mn.GetEngagementLikeName(name)
			if inpututxo.Value<mn.freezingcoef*uint64(totalstakedislike-totalstakelike){
				return 0,fmt.Errorf("Deposit frozen")
			}
		*/
		return inpututxo.Value, nil
	/*
		case utility.ModuleIdentifierECDSANamePublicPost:
				if (mn.GetTxOutputState(tmptxin.Hash,tmptxin.Index)!=StateValueIdentifierActifNameRegistration){
					return 0,fmt.Errorf("invalide ** GetTxOutputState %d for %x %d ",mn.GetTxOutputState(tmptxin.Hash,tmptxin.Index),tmptxin.Hash,tmptxin.Index)
				}
				_,height,number:=mn.GetTxState(tmptxin.Hash)
				//applog.Trace("height%d,number%d",height,number)
				inpututxo:=mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]
				pubkeycompressedbytes,_,err:=utility.DecodeECDSANamePublicPost(tmptxin.Bytecode)
				if err!=nil{
					return 0,fmt.Errorf("DecodeECDSANamePublicPost Error %v",err)
				}
				verr:=utility.VerifySignature(signinghash,tmptxin.Signature,pubkeycompressedbytes)
				if verr!=nil{
					return 0,verr
				}
				if !(inpututxo.CompareWithAddress(utility.ComputeHash (pubkeycompressedbytes))){
				//if (utility.ComputeHash (pubkeycompressedbytes)!=inpututxo.Address){//TODO
					return 0,fmt.Errorf("Corrupt transaction input - input public key do not match its associated output hash")
				}
				//applog.Trace("Value",inpututxo.Value)
				//
				//_,name,_,_,_:=utility.DecodeECDSANameRegistration(inpututxo.Bytecode)
				//_,totalstakedislike:=mn.GetEngagementDislikeName(name)
				//_,totalstakelike:=mn.GetEngagementLikeName(name)
				//if inpututxo.Value<mn.freezingcoef*uint64(totalstakedislike-totalstakelike){
				//	return 0,fmt.Errorf("Deposit frozen")
				//}
				//
				return 0,nil

		case utility.ModuleIdentifierECDSAEngagementPublicPostRewardClaim:
			//pubkey,_,_:=utility.DecodeECDSAEngagementRewardClaim(tmptxin.Bytecode)
			//rewardclaimaddress:=utility.ComputeHash(pubkey)
			//stateid,_,_:=mn.GetEngagementPublicPostState(tmptxin.Hash,tmptxin.Index,rewardclaimaddress)
			//if stateid!=StateValueIdentifierUnclaimedEngagementPublicPost{
			//	return 0,fmt.Errorf("Engagement reward claim points to engagement already claimed or inexisting")
			//}
			if (mn.GetTxOutputState(tmptxin.Hash,tmptxin.Index)!=StateValueIdentifierUnclaimedEngagementPublicPost){
				return 0,fmt.Errorf("Engagement reward claim points to engagement already claimed or inexisting - GetTxOutputState %d for %x %d ",mn.GetTxOutputState(tmptxin.Hash,tmptxin.Index),tmptxin.Hash,tmptxin.Index)
			}
			_,height,number:=mn.GetTxState(tmptxin.Hash)//URGENT TODO
			//applog.Trace("height%d,number%d",height,number)
			engagementtxout:=mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]//URGENT TODO
			engagementid,publicposttxhash,publicposttxindex,claimaddress,_,_:=utility.DecodeECDSAEngagement(engagementtxout.Bytecode)
			_,publicpostheight,_:=mn.GetTxState(*publicposttxhash)//URGENT TODO
			if (mn.GetConfirmedMainchainLength()-publicpostheight)<ENGAGEMENT_REWARD_FINALIZATION_INTERVAL {
				return 0,fmt.Errorf("Engagement reward claim was made too soon")
			}

			//
			//tmptxin.Bytecode

				pubkeycompressedbytes,_,err:=utility.DecodeECDSAEngagementRewardClaim(tmptxin.Bytecode)
				if err!=nil{
					return 0,fmt.Errorf("DecodeECDSAEngagementRewardClaim Error %v",err)
				}
				verr:=utility.VerifySignature(signinghash,tmptxin.Signature,pubkeycompressedbytes)
				if verr!=nil{
					return 0,verr
				}
				//if !(inpututxo.CompareWithAddress(utility.ComputeHash (pubkeycompressedbytes))){
				if (utility.ComputeHash (pubkeycompressedbytes)!=*claimaddress){//TODO
					return 0,fmt.Errorf("Corrupt transaction input of engagement reward claim- input public key do not match its associated engagement claimaddress")
				}
				//applog.Trace("Value",inpututxo.Value)
				//return inpututxo.Value,nil

			//
			stakeheight:= ENGAGEMENT_REWARD_FINALIZATION_INTERVAL- (height-publicpostheight)
			engagementreward,err:=mn.GetEngagementPublicPostRewardValue(*publicposttxhash,publicposttxindex,publicpostheight,engagementid,engagementtxout.Value,stakeheight)
			if err!=nil {
				return 0,err
			}
			applog.Trace("Adding valid GetEngagementPublicPostRewardValue - %d",engagementreward)
			return engagementreward,nil

	*/
	default:
		return 0, fmt.Errorf("Unknown moduleid of TxIn")
	}
}
func (mn *Maincore) ValidateTxOut(tmptxout utility.TxOut) (uint64, error) {
	moduleid := utility.DecodeBytecodeId(tmptxout.Bytecode)
	//applog.Trace("txout moduleid %d",moduleid)
	switch moduleid {
	case utility.ModuleIdentifierECDSATxOut:
		_, _, err := utility.DecodeECDSATxOutBytecode(tmptxout.Bytecode)
		if err != nil {
			return 0, err
		}
		return tmptxout.Value, nil

	case utility.ModuleIdentifierECDSANameRegistration:
		_, name, _, _, err := utility.DecodeECDSANameRegistration(tmptxout.Bytecode)
		if err != nil {
			return 0, err
		}
		verr := mn.ValidateNameRegistration(name)
		if verr != nil {
			return 0, verr
		}
		return tmptxout.Value, nil
		/*
			case utility.ModuleIdentifierECDSAEngagementPublicPost:
				_,_,_,_,_,err:=utility.DecodeECDSAEngagement(tmptxout.Bytecode)
				//eid,&hash,index,extradata,nil
				if err!=nil {
					return 0,err
				}
				if tmptxout.Value>ENGAGEMENT_PUBLICPOST_MAXSTAKE{
					return 0,fmt.Errorf("Engagement stake of tx exceeds ENGAGEMENT_PUBLICPOST_MAXSTAKE")
				}
				return tmptxout.Value,nil
		*/
		//default:
		//	return 0,fmt.Errorf("Unknown moduleid of TxOut")
	}
	return 0, fmt.Errorf("Unknown module")
}
