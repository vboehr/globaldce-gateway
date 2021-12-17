package mainchain
import (
	"github.com/globaldce/globaldce-toolbox/applog"
	"github.com/globaldce/globaldce-toolbox/utility"
	"github.com/globaldce/globaldce-toolbox/wallet"
	"strings"
	//"fmt"
	//"os"
)

func  (mn *Maincore) SyncWallet (wlt *wallet.Wallet){
	var wltPubkeyHash [] utility.Hash
	for i:=0;i<len(wlt.Privatekeyarray);i++{
	wltPubkeyHash=append(wltPubkeyHash,utility.ComputeHash(wlt.Privatekeyarray[i].PubKey().SerializeCompressed()))
	}
	applog.Trace("Maincore Syncing Wallet ...")

	applog.Trace("Maincore Syncing Wallet Adding Assets ...")
	for i:=wlt.Lastknownblock+uint64(1);i<uint64(mn.GetConfirmedMainchainLength());i++{
		mb:=mn.GetConfirmedMainblock(int(i))
		for j:=0;j<len(mb.Transactions);j++{
			for k:=0;k<len(mb.Transactions[j].Vout);k++{
				for l:=0;l<len(wltPubkeyHash);l++{
					if mb.Transactions[j].Vout[k].CompareWithAddress(wltPubkeyHash[l]){
						assetstate:=mb.Transactions[j].Vout[k].GetAssetState()
						wlt.AddAsset(mb.Transactions[j].ComputeHash(),uint32(k),mb.Transactions[j].Vout[k].Value,uint32(l),assetstate)
					}
				}
			}
			//
			//
			for k:=0;k<len(mb.Transactions[j].Vin);k++{
						for m:=0;m<len(wlt.Assetarray);m++{
							if wlt.Assetarray[m].Hash==mb.Transactions[j].Vin[k].Hash && wlt.Assetarray[m].Index==mb.Transactions[j].Vin[k].Index {
								//
								//utility.DecodeBytecodeId(mb.Transactions[j].Vin[k].Bytecode)
								if wlt.Assetarray[m].StateString=="UNSPENT" {
									wlt.Assetarray[m].StateString="SPENT"
								} else {
									moduleid:=utility.DecodeBytecodeId(mb.Transactions[j].Vin[k].Bytecode)
									
									if strings.Index(wlt.Assetarray[m].StateString,"NAMEREGISTERED_")==0 && moduleid==utility.ModuleIdentifierECDSANameUnregistration{
										r := strings.NewReplacer("NAMEREGISTERED_", "NAMEUNREGISTERED_")
										wlt.Assetarray[m].StateString=r.Replace(wlt.Assetarray[m].StateString)
										//fmt.Println("Unregistration occured",j,k)
										//os.Exit(0)
									}
								}
								
	
								
								//
								wlt.Assetarray[m].SpendingTxHash=mb.Transactions[j].ComputeHash()
								wlt.Assetarray[m].SpendingIndex=uint32(k)

							}
						}
				//
			}
			//
		}
		wlt.Lastknownblock=i
	}
	applog.Trace("Maincore Syncing Wallet Refreshing Broadcasted Transactions Confirmation ...")	
	for i:=0;i<len(wlt.Broadcastedtxarray);i++{

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		tmpstate,_,_:=mn.GetTxState(wlt.Broadcastedtxarray[i].Tx.ComputeHash())
		if tmpstate==StateValueIdentifierTx{
			//applog.Trace("unspents")
			wlt.Broadcastedtxarray[i].ConfirmationString="CONFIRMED"
		}

	}
	
	for i:=0;i<len(wlt.Broadcastedtxarray);i++{
			if wlt.Broadcastedtxarray[i].ConfirmationString=="CONFIRMED"{
				txhash:=wlt.Broadcastedtxarray[i].Tx.ComputeHash()
				mn.txspool.DeleteTransaction(&txhash)
			}			
	}
	applog.Notice("Maincore Syncing Wallet done.")
}


/*
func  (mn *Maincore) SyncWallet (wlt *wallet.Wallet){
	//wlt.Lastknownblock=0
	var wltPubkeyHash [] utility.Hash
	for i:=0;i<len(wlt.Privatekeyarray);i++{
	wltPubkeyHash=append(wltPubkeyHash,utility.ComputeHash(wlt.Privatekeyarray[i].PubKey().SerializeCompressed()))
	}
	applog.Trace("Maincore Syncing Wallet ...")
	//applog.Trace("wltPubkeyHash %x",wltPubkeyHash)
	applog.Trace("Maincore Syncing Wallet Adding Assets ...")
	for i:=wlt.Lastknownblock+uint64(1);i<uint64(mn.GetConfirmedMainchainLength());i++{
		mb:=mn.GetConfirmedMainblock(int(i))
		for j:=0;j<len(mb.Transactions);j++{
			for k:=0;k<len(mb.Transactions[j].Vout);k++{
				for l:=0;l<len(wltPubkeyHash);l++{
					if mb.Transactions[j].Vout[k].CompareWithAddress(wltPubkeyHash[l]){
					
						//applog.Notice("Adding asset ...")
						wlt.AddAsset(mb.Transactions[j].ComputeHash(),uint32(k),mb.Transactions[j].Vout[k].Value,uint32(l))
					}
				}
			}
			//
			//
			for k:=0;k<len(mb.Transactions[j].Vin);k++{
				//tmpindex:=mb.Transactions[j].Vin[k].Index
				//txstate,height,number:=mn.GetTxState(mb.Transactions[j].Vin[k].Hash)
				//if txstate!=mainchain.StateValueIdentifierTx{
				//	applog.Warning("Error: Transaction not found - hash %x",mb.Transactions[j].Vin[k].Hash)
				//	continue
				//}

				//for l:=0;l<len(wltPubkeyHash);l++{
					//tmptx:=mn.GetConfirmedMainblock(int(height)).Transactions[number]
					//if tmptx.Vout[tmpindex].Pubkeyhash==wltPubkeyHash[l]{
						//applog.Notice("Adding asset ...")
						//wlt.AddAsset(mb.Transactions[j].ComputeHash(),uint32(k),mb.Transactions[j].Vout[k].Value,uint32(l))
						//applog.Trace("Spent amount %d transaction hash %x\n",mb.Transactions[j].Vout[k].Value,mb.Transactions[j].ComputeHash())
						for m:=0;m<len(wlt.Assetarray);m++{
							if wlt.Assetarray[m].Hash==mb.Transactions[j].Vin[k].Hash && wlt.Assetarray[m].Index==mb.Transactions[j].Vin[k].Index {
								// assuming that the transaction signature is already verified
								wlt.Assetarray[m].StateString="SPENT"
								wlt.Assetarray[m].SpendingTxHash=mb.Transactions[j].ComputeHash()
								wlt.Assetarray[m].SpendingIndex=uint32(k)
							}
						}
					//}
				//}
				//
			}
			//
		}
		wlt.Lastknownblock=i
	}
	//applog.Trace("Maincore Syncing Wallet Refreshing Assets State ...")
	//for i:=0;i<len(wlt.Assetarray);i++{

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		//if mn.GetTxOutputState(wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index)==StateValueIdentifierUnspentTxOutput{
		//	//applog.Trace("unspents")
		//	wlt.Assetarray[i].StateString="UNSPENT"
		//}


	//}
	applog.Trace("Maincore Syncing Wallet Refreshing Broadcasted Transactions Confirmation ...")	
	for i:=0;i<len(wlt.Broadcastedtxarray);i++{

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		tmpstate,_,_:=mn.GetTxState(wlt.Broadcastedtxarray[i].Tx.ComputeHash())
		if tmpstate==StateValueIdentifierTx{
			//applog.Trace("unspents")
			wlt.Broadcastedtxarray[i].ConfirmationString="CONFIRMED"
		}

	}
	
	for i:=0;i<len(wlt.Broadcastedtxarray);i++{
			//applog.Trace("unspents")
			//if wlt.Broadcastedtxarray[i].ConfirmationString==""{
			//	mn.txspool.AddTransaction(&wlt.Broadcastedtxarray[i].Tx,&wlt.Broadcastedtxarray[i].Fee,50000000000000)
			//}
			if wlt.Broadcastedtxarray[i].ConfirmationString=="CONFIRMED"{
				txhash:=wlt.Broadcastedtxarray[i].Tx.ComputeHash()
				mn.txspool.DeleteTransaction(&txhash)
			}			
	}

	applog.Notice("Maincore Syncing Wallet done.")

}
*/