
package mainchain

import (
	"github.com/globaldce/globaldce-toolbox/applog"

	"github.com/globaldce/globaldce-toolbox/utility"
	//"fmt"
	//"math/big"
	//leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"//
	//leveldbutil "github.com/syndtr/goleveldb/leveldb/util"//
)



func  (mn *Maincore)  UpdateMainstate(tx utility.Transaction,blockheight uint32) {
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
				applog.Trace("Updating ECDSAEngagement TxHash %x %d",txhash,uint32(k))
				mn.PutTxOutputState(txhash,uint32(k),StateValueIdentifierUnclaimedEngagementPublicPost)
				eid,pptxhash,pptxindex,_,_,_:=utility.DecodeECDSAEngagement(tx.Vout[k].Bytecode)
				_,height,number:=mn.GetTxState(*pptxhash)

				ppinput:=mn.GetMainblock(int(height)).Transactions[number].Vin[pptxindex]
				stakeheight:=ENGAGEMENT_REWARD_FINALIZATION_INTERVAL-(blockheight-height)
				_,name,_,_:=utility.DecodeECDSANameRegistration(ppinput.Bytecode)
				if eid==utility.EngagementIdentifierLikePublicPost {
					applog.Trace("EngagementIdentifierLikePublicPost hash %x index %d stakeheight %d height %d blockheight %d",*pptxhash,pptxindex,stakeheight,height,blockheight)
					mn.AddEngagementLikeName(name,tx.Vout[k].Value)
					mn.AddEngagementPublicPostRewardLike(*pptxhash,pptxindex,tx.Vout[k].Value,stakeheight)
				}
				if eid==utility.EngagementIdentifierDislikePublicPost {
					applog.Trace("EngagementIdentifierDislikePublicPost hash %x index %d",*pptxhash,pptxindex)
					mn.AddEngagementDislikeName(name,tx.Vout[k].Value)
					mn.AddEngagementPublicPostRewardDislike(*pptxhash,pptxindex,tx.Vout[k].Value,stakeheight)
				}
				
				
				//mn.PutEngagementPublicPostState(*pptxhash,pptxindex,*claimaddress,StateValueIdentifierUnclaimedEngagementPublicPost,txhash,uint32(k) )
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
				
				if int(number)>=len(mn.GetMainblock(int(height)).Transactions){
					return
				}
				if int(tx.Vin[l].Index)>=len(mn.GetMainblock(int(height)).Transactions[number].Vout){
					return
				}
				
				tmpinpututxo:=mn.GetMainblock(int(height)).Transactions[number].Vout[tx.Vin[l].Index]
				_,name,_,_:=utility.DecodeECDSANameRegistration(tmpinpututxo.Bytecode)
			
				_,ed,_:=utility.DecodeECDSANamePublicPost(tx.Vin[l].Bytecode)
				
				//_,txhash,txindex,_,err:=mn.GetPublicPostState(ed.Hash)//([]byte, uint32, error){
				if ed!=nil{
					oldname,_,_,_,err:=mn.GetPublicPostState(ed.Hash)
					if (len(oldname)==0)&&(!mn.IsBannedName(name)){
						mn.PutPublicPostState(ed.Hash,name,txhash,uint32(l),uint32(0))
					}
					if (err!=nil)&&(!mn.IsBannedName(name)){
						mn.PutPublicPostState(ed.Hash,name,txhash,uint32(l),uint32(0))
						mn.AddToMissingDataHashArray(ed.Hash)	
					}
				}

			case utility.ModuleIdentifierECDSAEngagementPublicPostRewardClaim:
				applog.Trace("Updating ECDSAEngagementRewardClaim Tx")
				//tmptxin:=tx.Vin[l]
				//pubkey,_,_:=utility.DecodeECDSAEngagementRewardClaim(tmptxin.Bytecode)
				//rewardclaimaddress:=utility.ComputeHash(pubkey)
				//_,etxhash,etxindex:=mn.GetEngagementPublicPostState(tmptxin.Hash,tmptxin.Index,rewardclaimaddress)
				//_,height,number:=mn.GetTxState(tmptxin.Hash)//URGENT TODO
				//applog.Trace("height%d,number%d",height,number)
				//engagementtxout:=mn.GetMainblock(int(height)).Transactions[number].Vout[tmptxin.Index]//URGENT TODO
				//_,pptxhash,pptxindex,_,_,_:=utility.DecodeECDSAEngagement(engagementtxout.Bytecode)
				//mn.PutEngagementPublicPostState(*pptxhash,pptxindex,rewardclaimaddress,StateValueIdentifierClaimedEngagementPublicPost,*etxhash,etxindex )
				mn.PutTxOutputState(tx.Vin[l].Hash,tx.Vin[l].Index,StateValueIdentifierClaimedEngagementPublicPost)
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