package mainchain
import (
	"github.com/globaldce/go-globaldce/applog"
	//"github.com/globaldce/go-globaldce/wire"
	//"github.com/globaldce/go-globaldce/mainchain"
	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/globaldce/go-globaldce/utility"
	"github.com/globaldce/go-globaldce/wallet"
	//"github.com/globaldce/go-globaldce/mainchain"
	//"os"
	"time"
	//"encoding/json"
	"math/big"
	//"net"
	//"log"
)


////////////////////////////////////////////////////////////////////
func (mn *Maincore) Mine(wlt *wallet.Wallet) (bool,*Mainblock){
	applog.Notice("Mining")
	mainblock:=NewMainblock()
	mainblock.Height=uint32 (mn.GetMainchainLength())
	pubkeyhash:=wlt.GetLastAddress()
	selectedtxs,totalfees:= mn.txspool.GetHighestPriorityTxs(mn.GetConfirmedMainchainLength())

	rewardtx:=utility.NewRewardTransaction(GetMainblockReward(mainblock.Height),uint64(totalfees),pubkeyhash)//GENESIS_BLOCK_REWARD
	//applog.Trace("transaction %d ",rewardtx.Vout[0].Value)
	//rewardtxhash:=rewardtx.ComputeHash()
	//mainblock.Header.Roothash=rewardtxhash
	mainblock.Transactions=append(mainblock.Transactions,*rewardtx)
	
	for _, selectedtx := range *selectedtxs {
		mainblock.Transactions=append(mainblock.Transactions,selectedtx)
		//applog.Trace("Add a transaction !!!!!!!!!!!!!!!!!")
		
	}

	mainblock.ComputeRoot()
	//////////////////////
	var targetbits uint32
	mainchainlength:=int(mn.GetMainchainLength())
	//TODO replace code by associated function ??? GetTargetBits()????
	if (mainchainlength % int ( DIFFICULTY_TUNING_INTERVAL)==0){
		
		targetbigint:=utility.BigIntFromCompact( mn.GetLastMainheader().Bits)
		idealtimeinterval:=int64 ( DIFFICULTY_TUNING_INTERVAL-1)*600
		realtimeinterval:=int64 ( mn.GetLastMainheader().Timestamp- mn.GetMainheader(mainchainlength-int ( DIFFICULTY_TUNING_INTERVAL)).Timestamp)
		
		if (realtimeinterval>=int64 (3)*idealtimeinterval) {
			targetbigint.Mul(targetbigint,big.NewInt(3))
			applog.Trace("bigger than 3  ")
		} else if (idealtimeinterval>=int64 (3) *realtimeinterval) {
			targetbigint.Div(targetbigint,big.NewInt(3))
			applog.Trace("smaller that 1/3 ")
		} else {
			targetbigint.Mul(targetbigint,big.NewInt(realtimeinterval))
			targetbigint.Div(targetbigint,big.NewInt(idealtimeinterval))
		}

		targetbits = utility.CompactFromBigInt(targetbigint)

		applog.Trace("\n realtime %d idealtime %d  targetbitint %d ",realtimeinterval,idealtimeinterval,targetbigint)
	} else {
		
		targetbits = mn.GetMainheader(mainchainlength-1).Bits

	}
	///////////////////
	tmplastmainheader:=mn.GetLastMainheader()
	//TODO unlocking maincore and wallet
	time.Sleep(1*time.Second)
	applog.Trace("Lets go .......")
	success:=mainblock.Mine( tmplastmainheader.Timestamp, tmplastmainheader.Hash,targetbits)
	//TODO locking maincore and wallet
	if success && mainblock.Height==mn.GetMainchainLength() {
		if mn.GetMainblock(int(mainblock.Height)-1).Header.Hash==mainblock.Header.Prevblockhash{
			mn.AddInMemoryBlock(mainblock)
			mn.ConfirmBlocks()
			applog.Trace("Mined mainblock added to mainchain")
			applog.Trace("Mainchain length %d Mainchain confirmed length %d",mn.GetMainchainLength(),mn.mainbf.NbChunks())

			wlt.GenerateKeyPair()
			return true,mainblock
		} else {
			applog.Warning("mined mainblock rejected - hash of mainblock found do not match previous block hash")
		}
			
		} else {
			if !success{
				applog.Warning("Mining unsuccessful - no appropiate mainblock found")
			} else {
				applog.Warning("mined mainblock rejected - mainblock.Height %d mn.GetMainchainLength() %d ",mainblock.Height,mn.GetMainchainLength())
			}
		}
		
 return false,nil
}
func (mn *Maincore) ConfirmMainblock(mb*Mainblock) {
	mn.mainbf.AddChunk(mb.Serialize())
	mn.mainheaders=append(mn.mainheaders,mb.Header)
	i:=mb.Height
	if !mn.ValidateMainblockTransactions(uint32 (i), &mb.Transactions){
		applog.Warning("ConfirmMainblock - Invalid mainblock %d",i)
		mn.PutMainblockState(uint32 (i),StateValueIdentifierUnvalidMainblock)
	} else {
		mn.PutMainblockState(uint32 (i),StateValueIdentifierValidMainblock)
	}
	
	for j:=0;j<len(mb.Transactions);j++{
		//applog.Trace("block %d %x ",i,mn.GetMainblock(i).Transactions[j].ComputeHash())
		tx:=mb.Transactions[j]
		txhash:=tx.ComputeHash()
		mn.txspool.DeleteTransaction(&txhash)
		mn.PutTxState(txhash,uint32(i),uint32(j))
		mn.UpdateMainstate(tx,uint32(i))
		//
	}
}
func (mn *Maincore) ConfirmBlocks() {
	if len( mn.inmemorymainblocks)==0 {
		return 
	} else {
		firstminedblockheight:=int(mn.inmemorymainblocks[0].Height)
		if firstminedblockheight>mn.mainbf.NbChunks(){
			applog.Trace("error: inmemory do not correspond to mainchain stored blocks")
			return
		}
		if uint32(len( mn.inmemorymainblocks))>mn.GetConfirmationLayer(){
			tmpmb:=mn.inmemorymainblocks[len( mn.inmemorymainblocks)-int(mn.GetConfirmationLayer())-1]
			if tmpmb.Height==uint32(mn.mainbf.NbChunks()){
				applog.Trace("adding inmemory %d to mainchain stored blocks",tmpmb.Height)
				
				/////////////////////////////////////
				mn.ConfirmMainblock(&tmpmb)
				return
			} else{
				applog.Trace("error: inmemory do not correspond to mainchain stored blocks")
				return
			}
		}
	}

}
func (mn *Maincore) GetLastMainheader() *Mainheader{
 if len( mn.inmemorymainblocks)==0 {
	return &mn.mainheaders[mn.mainbf.NbChunks()-1]
 } else {
	 tmplength:=len(mn.inmemorymainblocks)
	 return &mn.inmemorymainblocks[tmplength-1].Header
 }

}
func (mn *Maincore) GetMainheader(i int) *Mainheader{
	if i<mn.mainbf.NbChunks(){
	   return &mn.mainheaders[i]
	} else if len(mn.inmemorymainblocks)>0 {
		firstminedblockheight:=int(mn.inmemorymainblocks[0].Height)
		tmplength:=len(mn.inmemorymainblocks)
		if i<firstminedblockheight+tmplength{
			return &mn.inmemorymainblocks[i-firstminedblockheight].Header
		}  
	} 
   return nil
   }
func (mn *Maincore) GetMainblock(i int) *Mainblock{
	//
	if (i<mn.mainbf.NbChunks()){
		return mn.GetConfirmedMainblock(i)
	 } else if len( mn.inmemorymainblocks)>0 {
		 firstminedblockheight:=int(mn.inmemorymainblocks[0].Height)
		 tmplength:=len(mn.inmemorymainblocks)
		 if i<firstminedblockheight+tmplength{
			 return &mn.inmemorymainblocks[i-firstminedblockheight]
		 }
	 }
	 return nil

}

////////////////////////////////////////////////////////////////////