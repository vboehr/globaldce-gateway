
package mainchain

import (
	"github.com/globaldce/globaldce-toolbox/applog"

	"github.com/globaldce/globaldce-toolbox/utility"
	//"fmt"
	//"math/big"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)



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
			case utility.ModuleIdentifierECDSAEngagementPublicPost:
				applog.Trace("Updating ECDSAEngagement Tx")
				eid,pptxhash,pptxindex,claimaddress,_,_:=utility.DecodeECDSAEngagement(tx.Vout[k].Bytecode)
				_,height,number:=mn.GetTxState(*pptxhash)

				ppinput:=mn.GetMainblock(int(height)).Transactions[number].Vin[pptxindex]

				_,name,_,_:=utility.DecodeECDSANameRegistration(ppinput.Bytecode)
				if eid==utility.EngagementIdentifierLikePublicPost {
					mn.AddEngagementLikeName(name,tx.Vout[k].Value)
				}
				if eid==utility.EngagementIdentifierDislikePublicPost {
					mn.AddEngagementDislikeName(name,tx.Vout[k].Value)
				}
				mn.PutEngagementPublicPostState(*pptxhash,pptxindex,*claimaddress,StateValueIdentifierUnclaimedEngagementPublicPost,txhash,uint32(k) )
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
				
				//_,txhash,txindex,_,err:=mn.GetPublicPostState(ed.Hash)//([]byte, uint32, error){
				_,_,_,_,err:=mn.GetPublicPostState(ed.Hash)
				
				if (err!=nil)&&(!mn.IsBannedName(name)){
					mn.PutPublicPostState(ed.Hash,name,txhash,uint32(l),uint32(0))
					mn.AddToMissingDataHashArray(ed.Hash)	
				}
			case utility.ModuleIdentifierECDSAEngagementPublicPostRewardClaim:
				applog.Trace("Updating ECDSAEngagementRewardClaim Tx")
				//_,height,number:=mn.GetTxState(tmptxin.Hash)//URGENT TODO
				//applog.Trace("height%d,number%d",height,number)
				//engagementtxout:=mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]//URGENT TODO
				//_,height,number:=mn.GetTxState(tx.Vin[l].Hash)
				//applog.Trace("height%d,number%d",height,number)
				//engagementtxout:=mn.GetMainblock(int(height)).Transactions[number].Vout[tx.Vin[l].Index]

				//_,publicposttxhash,publicposttxindex,_,_,_:=utility.DecodeECDSAEngagement(engagementtxout.Bytecode)
				//
				//
				//mn.PutEngagementPublicPostState(tx.Vin[l].Hash,tx.Vin[l].Index,StateValueIdentifierClaimedEngagementPublicPost)
				//

			//default:
		}
		
		
	}

}