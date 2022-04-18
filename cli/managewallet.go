
package cli
import (
	"github.com/globaldce/go-globaldce/applog"
	"github.com/globaldce/go-globaldce/daemon"
	"fmt"
	"github.com/globaldce/go-globaldce/mainchain"
	"github.com/globaldce/go-globaldce/utility"
	"github.com/globaldce/go-globaldce/wallet"
	"strings"
	"os"
	"bufio"
    "strconv"
    "encoding/hex"
    "github.com/globaldce/go-globaldce/wire"
	"path/filepath"
)

func GenerateNewMainwalletFilePath(){
	for activewalletfileid:=1;activewalletfileid<=1000;activewalletfileid++{
		filepath:=filepath.Join("WalletFiles",fmt.Sprintf("Wallet%02d.dat",activewalletfileid))//fmt.Sprintf("WalletFiles/Wallet%02d",activewalletfileid)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			daemon.MainwalletFilePath=filepath
			daemon.Usersettings.MainwalletFilePath=filepath
			break
		} 
		if (activewalletfileid==1000){
			applog.Trace(" Cannot generate new wallet file name: Too much wallet files ! ")
		}
	}
}



func managewallet(ws *wire.Swarm,mn *mainchain.Maincore,wlt * wallet.Wallet){
	mn.Lock()
    mn.SyncWallet(wlt)
	mn.Unlock()
	applog.LockDisplay()
	fmt.Printf("\n")
	fmt.Printf("Wallet last known block height: %d\n",wlt.Lastknownblock)
    //fmt.Printf("Wallet ballance: %d\n",wlt.ComputeBalance())
	applog.Trace("Wallet ballance %f",float64 (wlt.ComputeBalance())/1000000.0)
    //fmt.Printf("\nSaving mainwallet file ...")
    //wlt.SaveJSONFile(mainwalletFilePath,mainwalletFileKey)
	
	fmt.Printf("\nEntered manage wallet mode. Here are some options to choose from:\n")
	managewallethelp()
	fmt.Printf("\n")
	
	for {
		//fmt.Printf(">")
		var requeststring string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			requeststring = scanner.Text()
			//fmt.Println(">",requeststring)
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Fatal error",err)
			os.Exit(1)
		}
		//mn.Lock()
		requestarguments := strings.Fields(requeststring)
		if len(requestarguments)==0{
			managewallethelp()
			continue
		}
		switch strings.ToLower(requestarguments[0]) {
			case "generateaddress":
				generateaddress(wlt)
			case "sendtoaddress":
				if len(requestarguments)<3{
					sendtoaddresshelp()
					continue
				}
				
				sendtoaddress(ws,wlt,requestarguments[1],requestarguments[2])
			case "sendnameregistration":
				if len(requestarguments)<3{
					sendnameregistrationhelp()
					continue
				}
				
				Sendnameregistration(ws,mn,wlt,requestarguments[1],requestarguments[2])
			case "sendpublicpost":
				if len(requestarguments)<4{
					sendpublicposthelp()
					continue
				}
				var filepatharray []string
				for i:=4;i<len(requestarguments);i++{
					filepatharray=append(filepatharray,requestarguments[i]) 
				}
				feestring:=fmt.Sprintf("%d",daemon.Usersettings.Publicposttxfee)
				Sendpublicpost(ws,mn,wlt,requestarguments[1],requestarguments[2],requestarguments[3],filepatharray,feestring)
			case "displaybalance":
				displaybalance(mn,wlt)
			case "displayaddresses":
				displayaddresses(wlt)
			case "displayalltransactions":
				displayalltransactions(wlt)
			case "displaybroadcastedtransactions":
				displaybroadcastedtransactions(wlt)
			case "displayregistrednames":
				displayregistrednames(wlt)
			case "scanalltransactions":
				scanalltransactions(mn,wlt)
			case "changepassword":
				changepassword()
			case "help":
				managewallethelp()
			case "h":
				managewallethelp()
			default:
				fmt.Printf("\nUnkown manage wallet option.\n\n")
				managewallethelp()
		}
	}


}
func managewallethelp(){
	fmt.Printf("MANAGE WALLET OPTIONS:\n")
	fmt.Printf("generateaddress			  				\n")
	fmt.Printf("sendtoaddress			  				\n")
	fmt.Printf("sendnameregistration	  				\n")
	fmt.Printf("sendpublicpost			  				\n")
	fmt.Printf("displaybalance  						\n")
	fmt.Printf("displayaddresses	  					\n")
	fmt.Printf("displayalltransactions  				\n")
	fmt.Printf("displaybroadcastedtransactions          \n")
	fmt.Printf("displayregistrednames			        \n")
	fmt.Printf("scanalltransactions				        \n")
	fmt.Printf("changepassword					        \n")
	fmt.Printf("help or h               			provides description of manage wallet mode options")
	fmt.Printf("\n")
}
func changepassword(){
	daemon.MainwalletFileKey=askuserwalletfilekey()
	if len(daemon.MainwalletFileKey)!=0{
		fmt.Printf("\nWallet password set to: %s\n",daemon.MainwalletFileKey)
	}
	
}
func generateaddress(wlt *wallet.Wallet){
	address:=wlt.GenerateKeyPair()
	fmt.Printf("\nGenerated address %x\n",address)
}
func displaybalance(mn *mainchain.Maincore,wlt * wallet.Wallet){
	mn.SyncWallet(wlt)
	fmt.Printf("\nLastknownblock %d\n",wlt.Lastknownblock)
    fmt.Printf("Wallet ballance %f\n",float64 (wlt.ComputeBalance())/1000000.0)
	fmt.Printf("\n")
}
func displayaddresses(wlt * wallet.Wallet){
	fmt.Printf("\nNumber of addresses %d\n",len(wlt.Privatekeyarray))
	for i:=0;i<len(wlt.Privatekeyarray);i++{
	address:=utility.ComputeHash(wlt.Privatekeyarray[i].PubKey().SerializeCompressed())
	fmt.Printf("Address %d %x\n",i,address)
	}

}

func scanalltransactions(mn *mainchain.Maincore,wlt * wallet.Wallet){
	var wltPubkeyHash [] utility.Hash
	fmt.Printf("\nNumber of keys %d\n",len(wlt.Privatekeyarray))
	for i:=0;i<len(wlt.Privatekeyarray);i++{
	wltPubkeyHash=append(wltPubkeyHash,utility.ComputeHash(wlt.Privatekeyarray[i].PubKey().SerializeCompressed()))
	}

	for i:=uint32(1);i<mn.GetConfirmedMainchainLength();i++{
		mb:=mn.GetConfirmedMainblock(int(i))
		for j:=0;j<len(mb.Transactions);j++{
			//
			for k:=0;k<len(mb.Transactions[j].Vout);k++{
				for l:=0;l<len(wltPubkeyHash);l++{

					//if mb.Transactions[j].Vout[k].Address==wltPubkeyHash[l]{
					if mb.Transactions[j].Vout[k].CompareWithAddress(wltPubkeyHash[l]){
						//applog.Notice("Adding asset ...")
						//wlt.AddAsset(mb.Transactions[j].ComputeHash(),uint32(k),mb.Transactions[j].Vout[k].Value,uint32(l))
						fmt.Printf("Received amount %f transaction hash %x\n",float64 (mb.Transactions[j].Vout[k].Value)/1000000.0,mb.Transactions[j].ComputeHash())
					}
				}
			}
			//
			for k:=0;k<len(mb.Transactions[j].Vin);k++{
				tmpindex:=mb.Transactions[j].Vin[k].Index
				txstate,height,number:=mn.GetTxState(mb.Transactions[j].Vin[k].Hash)
				if txstate!=mainchain.StateValueIdentifierTx{
					fmt.Printf("Error: Transaction not found - hash %x\n",mb.Transactions[j].Vin[k].Hash)
					continue
				}

				for l:=0;l<len(wltPubkeyHash);l++{
					tmptx:=mn.GetConfirmedMainblock(int(height)).Transactions[number]
					//if tmptx.Vout[tmpindex].Address==wltPubkeyHash[l]{
					if tmptx.Vout[tmpindex].CompareWithAddress(wltPubkeyHash[l]){
						//applog.Notice("Adding asset ...")
						//wlt.AddAsset(mb.Transactions[j].ComputeHash(),uint32(k),mb.Transactions[j].Vout[k].Value,uint32(l))
						fmt.Printf("Sent amount %f transaction hash %x\n",float64 (mb.Transactions[j].Vout[k].Value)/1000000.0,mb.Transactions[j].ComputeHash())
					}
				}
			}
			//
		}
		//wlt.Lastknownblock=i
	}
}
func displaybroadcastedtransactions(wlt * wallet.Wallet){
	if len(wlt.Broadcastedtxarray) == 0 {
		fmt.Printf("\nNo broadcasted transactions stored in wallet.\n")
	}
	fmt.Printf("\n")
	for i:=0;i<len(wlt.Broadcastedtxarray);i++{

		//applog.Trace("hash %x index %d value %d pkindex %d ",wlt.Assetarray[i].Hash,wlt.Assetarray[i].Index,wlt.Assetarray[i].Value,wlt.Assetarray[i].Privatekeyindex)
		//tmpstate,_,_:=mn.GetTxState(wlt.Broadcastedtxarray[i].Tx.ComputeHash())
		//if tmpstate==StateValueIdentifierTx{
			//applog.Trace("unspents")
			fmt.Printf("Transaction %d hash %x state %s\n",i,wlt.Broadcastedtxarray[i].Tx.ComputeHash(),wlt.Broadcastedtxarray[i].ConfirmationString)
			fmt.Printf("Transaction %d - %s\n",i,wlt.Broadcastedtxarray[i].Tx.JSONSerialize())
		//}0

	}
}
//
func displayalltransactions(wlt * wallet.Wallet){
	if len(wlt.Assetarray) == 0 {
		fmt.Printf("\nNo asset transactions stored in wallet.\n")
	}
	for _, asset := range wlt.Assetarray {

		if asset.StateString=="UNSPENT"{
			fmt.Printf("Asset %f Hash %x State %s\n",float64 (asset.Value)/1000000.0,asset.Hash,asset.StateString)

		} else {
			fmt.Printf("Asset %f Hash %x State %s SpendingTxHash %x SpendingIndex %d\n",float64 (asset.Value)/1000000.0,asset.Hash,asset.StateString,asset.SpendingTxHash,asset.SpendingIndex)

		}
		

	}
}
func displayregistrednames(wlt * wallet.Wallet){
	registerednames:=wlt.GetRegisteredNames()
	for i,name :=range registerednames{
		fmt.Printf("Registred name %d : %s\n",i,name)
	}
}
// 
//
func Sendtoaddressarray(ws *wire.Swarm,wlt *wallet.Wallet,addrstringarray []string,amountstringarray []string,feestring string){
    //applog.Trace("Wallet ballance %f",float64 (wlt.ComputeBalance())/1000000.0)
    //applog.Trace("Wallet Lastknownblock %d",wlt.Lastknownblock)

	var addressarray []utility.Hash
	var amountarray []uint64

	for _,addrstring:=range addrstringarray {
	
		//address:=[]byte(addrstring)
		addressbytes,addrerr:=hex.DecodeString(addrstring)
		if addrerr!=nil{
			fmt.Printf("\nError: inappropriate address provided - %v",addrerr)
			return
		}
		if len(addressbytes)!=32{
			fmt.Printf("Error: inappropriate address length")
			return
		}
		addressarray=append(addressarray,*utility.NewHash(addressbytes)) 
		//applog.Notice("Sending to address %x an amount of %d globals",address,amount)

	}
	for _,amountstring:=range amountstringarray {
		amount, err := strconv.ParseInt(amountstring, 10, 64)
		amount*=1000000
		if err!=nil{
			fmt.Printf("Error: inappropriate amount provided - %v",err)
			return
		}
		amountarray=append(amountarray,uint64(amount))
	}
	//
    //amountfee:=amount*1/100//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
	fee, ferr := strconv.ParseInt(feestring, 10, 64)
    fee*=1000000
    if ferr!=nil{
        fmt.Printf("Error: inappropriate fee provided - %v",ferr)
        //os.Exit(0)
		return
    }
	//
	tx,err:=wlt.SetupTransactionToPublicKeyHashArray(addressarray,amountarray,uint64 (fee))
    if err!=nil{
        fmt.Printf("%v",err)
        //os.Exit(0)
		return
    }
	fmt.Printf("new sendtoaddress tx seize %d tx %x",len(tx.Serialize()),tx)
    if tx!=nil{
		_,fee:= daemon.Mn.ValidateTransaction(tx)
		priority:=fee
		daemon.Mn.AddTransactionToTxsPool(tx,fee,priority)
        wlt.AddBroadcastedtx(*tx)
		ws.BroadcastTransaction(tx)
    }
    
}

//
func sendtoaddress(ws *wire.Swarm,wlt *wallet.Wallet,addrstring string,amountstring string){
    //applog.Trace("Wallet ballance %f",float64 (wlt.ComputeBalance())/1000000.0)
    //applog.Trace("Wallet Lastknownblock %d",wlt.Lastknownblock)

    amount, err := strconv.ParseInt(amountstring, 10, 64)
    amount*=1000000
    if err!=nil{
        fmt.Printf("Error: inappropriate amount provided - %v",err)
        //os.Exit(0)
		return
    }
    //address:=[]byte(addrstring)
    address,addrerr:=hex.DecodeString(addrstring)
    if addrerr!=nil{
        fmt.Printf("\nError: inappropriate address provided - %v",addrerr)
        //os.Exit(0)
		return
    }
    applog.Notice("Sending to address %x an amount of %d globals",address,amount)
    if len(address)!=32{
        fmt.Printf("Error: inappropriate address length")
        //os.Exit(0)
		return
    }
    //amountfee:=amount*1/100//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
	amountfee:=daemon.Usersettings.Sendtoaddressarraytxfee//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
	amountfee*=1000000
    tx,err:=wlt.SetupTransactionToPublicKeyHash(*utility.NewHash(address),uint64 (amount),uint64 (amountfee))
    if err!=nil{
        fmt.Printf("%v",err)
        //os.Exit(0)
		return
    }
	fmt.Printf("new sendtoaddress tx seize %d tx %x",len(tx.Serialize()),tx)
    if tx!=nil{
		_,fee:= daemon.Mn.ValidateTransaction(tx)
		priority:=fee
		daemon.Mn.AddTransactionToTxsPool(tx,fee,priority)
        wlt.AddBroadcastedtx(*tx)
		ws.BroadcastTransaction(tx)
    }
    
}
func sendtoaddresshelp(){
    fmt.Printf("\nError: sendtoaddress inappropiate usage\n")
    fmt.Printf("In order to send an address X an amount of N globals, enter as follows:\n")
    fmt.Printf("sendtoaddress X N \n")
    //
}
func sendnameregistrationhelp(){
    fmt.Printf("\nError: sendnameregistration inappropiate usage\n")
    fmt.Printf("In order to proceed with registration of name X with an amount of N globals, enter as follows:\n")
    fmt.Printf("sendnameregistration X N \n")
    //
}
func sendpublicposthelp(){
    fmt.Printf("\nError: sendpublicpost inappropiate usage\n")
    fmt.Printf("In order to proceed with a public post for the registred name X with a web link of Y, a text Z and attached files F1 F2 ... Fn, enter as follows:\n")
    fmt.Printf("sendpublicpost X Y Z F1 F2 ... Fn\n")
    //
}
func Sendpublicpost(ws *wire.Swarm,mn *mainchain.Maincore,wlt *wallet.Wallet,namestring string,linkstring string,textstring string,filepatharray []string,amountfeestring string){
	/*
	a,aerr:=wlt.GetAssetFromRegisteredName(namestring)
    if aerr!=nil{
        fmt.Printf("%v",aerr)
        return
    }*/
	
	amountfee, ferr := strconv.ParseInt(amountfeestring, 10, 64)
    if ferr!=nil{
        fmt.Printf("Error: inappropriate fee provided - %v",ferr)
        //os.Exit(0)
		return
    }
	//amountfee:=amountfeestring//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
	amountfee*=1000000
	
	tmpbw:=utility.NewBufferWriter()
	tmpbw.PutUint32(mainchain.DataIdentifierPublicPost)
	tmpbw.PutVarUint(uint64(len([]byte(linkstring))))
	tmpbw.PutBytes([]byte(linkstring))

	tmpbw.PutVarUint(uint64(len([]byte(textstring))))
	tmpbw.PutBytes([]byte(textstring))
	/////////////////////////////////
	tmpbw.PutVarUint(uint64(len(filepatharray)))
	for _, filepath := range filepatharray {
		cfed,cerr:=mn.CacheExistingFile(filepath)
		if cerr!=nil{
			fmt.Printf("%v",cerr)
			return
		}
		tmpbw.PutVarUint((*cfed).Size)
		tmpbw.PutHash((*cfed).Hash)
	}
	/////////////////////////////////

	//databytes:=[]byte(datastring)//
	databytes:=tmpbw.GetContent()
	ed:=utility.NewExtradataFromBytes(databytes)

	tx,err:=wlt.SetupTransactionForNamePublicPost(namestring,ed,uint64 (amountfee))
    if err!=nil{
        fmt.Printf("%v",err)
        return
    }
    fmt.Printf("new public post seize %d tx %x",len(tx.Serialize()),tx)
    if tx!=nil{
		mn.AddLocalPublicPostData(namestring,ed.Hash,databytes)
		_,fee:= mn.ValidateTransaction(tx)
		priority:=fee
		mn.AddTransactionToTxsPool(tx,fee,priority)
        wlt.AddBroadcastedtx(*tx)
		ws.BroadcastTransaction(tx)
    }
}
func Sendnameregistration(ws *wire.Swarm,mn *mainchain.Maincore,wlt *wallet.Wallet,namestring string,amountstring string) error{
    //applog.Trace("Wallet ballance %f",float64 (wlt.ComputeBalance())/1000000.0)
    //applog.Trace("Wallet Lastknownblock %d",wlt.Lastknownblock)

    amount, err := strconv.ParseInt(amountstring, 10, 64)
    amount*=1000000
    if err!=nil{
        fmt.Printf("Error: inappropriate amount provided - %v",err)
		return fmt.Errorf("Error: inappropriate amount provided - %v",err)
    }

    name:=[]byte(namestring)

    fmt.Printf("Sending name registration for %x with amount of %d globals\n",name,amount)
	validationerr:=mn.ValidateNameRegistration(name)
	if validationerr!=nil{
		fmt.Printf("Validation error: %v\n",validationerr)
		return fmt.Errorf("Validation error: %v\n",validationerr)
	}

    amountfee:=daemon.Usersettings.Nameregistrationtxfee//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
	amountfee*=1000000
	pubkeyhash:=wlt.GenerateKeyPair()
    tx,err:=wlt.SetupTransactionForNameRegistration(name,pubkeyhash,uint64 (amount),uint64 (amountfee))
    if err!=nil{
        fmt.Printf("%v",err)
        return fmt.Errorf("%v",err)
    }
	fmt.Printf("new nameregistration seize %d tx %x",len(tx.Serialize()),tx)
    if tx!=nil{
		_,fee:= mn.ValidateTransaction(tx)
		priority:=fee
		mn.AddTransactionToTxsPool(tx,fee,priority)
        wlt.AddBroadcastedtx(*tx)
		ws.BroadcastTransaction(tx)
    }
	return nil
}
func Sendnameunregistration(ws *wire.Swarm,mn *mainchain.Maincore,wlt *wallet.Wallet,namestring string) error{
    //applog.Trace("Wallet ballance %f",float64 (wlt.ComputeBalance())/1000000.0)
    //applog.Trace("Wallet Lastknownblock %d",wlt.Lastknownblock)


    //name:=[]byte(namestring)
	
    fmt.Printf("Sending name unregistration for %x \n",namestring)
	amountfee:=100*1000000//TODO customizable fees based on bytes - fee = 1 to 10 globals * transaction bytes
    tx,err:=wlt.SetupTransactionForNameUnregistration(namestring,uint64 (amountfee))
    if err!=nil{
        fmt.Printf("%v",err)
        return fmt.Errorf("%v",err)
    }
	fmt.Printf("new nameunregistration seize %d tx %x",len(tx.Serialize()),tx)
    if tx!=nil{
		_,fee:= mn.ValidateTransaction(tx)
		priority:=fee
		mn.AddTransactionToTxsPool(tx,fee,priority)
        wlt.AddBroadcastedtx(*tx)
		ws.BroadcastTransaction(tx)
    }
	return nil
}