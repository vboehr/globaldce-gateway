package daemon

import (
    "github.com/globaldce/globaldce-gateway/applog"
    "os"
    "os/signal"
    "time"
    "github.com/globaldce/globaldce-gateway/utility"
    "github.com/globaldce/globaldce-gateway/mainchain"
    "github.com/globaldce/globaldce-gateway/wire"
    "log"
    "fmt"
	"github.com/globaldce/globaldce-gateway/content"
	"context"
    "github.com/globaldce/globaldce-gateway/rpc"
)
//

//////////////////////////////
func listenSigInt() chan os.Signal {
    //go func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
    //}()

}
/////////////////////////////

func startmining(){
    fmt.Println("Miningrequested",Miningrequested)
    for {
        //fmt.Printf("****** %v %v %v %v\n",Miningrequested ,Miningaddrressesfileloaded , Walletinstantiated, Wireswarm.Syncingdone )
        //fmt.Printf("?.????? %v\n",Miningrequested && (Miningaddrressesfileloaded || Walletinstantiated) && Wireswarm.Syncingdone )
        if Miningrequested && (Miningaddrressesfileloaded || Walletinstantiated) && Wireswarm.Syncingdone {
            var miningaddress utility.Hash
            if Walletinstantiated {
                miningaddress=Wlt.GenerateKeyPair()
            } else if Miningaddrressesfileloaded {
                miningaddress=MAddresses.GetRandomAddress()
            } else {
                fmt.Println("Not mining")
                continue
            }
            fmt.Println("Working..")
            Miningrunning=true
              //success,mb:=Mn.Mine(Wlt)

              
              success,mb:=Mn.Mine(miningaddress)
              if success {
                  Wireswarm.BroadcastMainblock(mb)
              }
              //if !Wlt.HotWallet{
                  
                  if Walletinstantiated && success{
                    Mn.SyncWallet(Wlt)
                    Wlt.SaveJSONWalletFile(MainwalletFilePath,MainwalletFileKey)
                  }
              //}
              
          //applog.Trace("mainchaillength %d",Mn.GetConfirmedMainchainLength())
            } else {
                Miningrunning=false
        }
        time.Sleep(3*time.Second)
    }
    
}

/////////////////////////////
func Mainloop(){
    applog.Init()
    
    
    Mn=mainchain.NewMaincore()
    Mn.PutPath( AppPath)
	Mn.LoadMaincore()
	go mainchain.InitLocalhost()
    //
	// Create a context with a cancellation function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Mncc:=content.Newcontentclient(ctx,"./Cache")
    
	go Mncc.Initcontentclient()
    rpc.RPCInit(Mncc)
    //
    //Mn.BannedNameArray=Usersettings.BannedNameArray
    Wireswarm=wire.NewSwarm()

    if Seed {
        fmt.Println("This is a seed")
        Wireswarm.Syncingdone=true
    }


	quitChan := listenSigInt()
	defer Mn.UnloadMaincore()
    err:=Wireswarm.SetupListener()
	if err != nil {
		log.Fatal(err)
	}
    
    //app_start_time:=time.Now().Unix()
    syncingdelaylength:=(time.Now().Unix()-Mn. GetLastConfirmedMainblockTimestamp())/int64(mainchain.BLOCK_TIME)
    applog.Notice("MainchainLength is %d",Mn.GetMainchainLength())
    applog.Notice("Mainchain is about %d blocks behind",syncingdelaylength)

    
    

    //

   
    Wireswarm.Bootstrap()
    ticker1 := time.Tick(time.Second * 7)
    //ticker2 := time.Tick(time.Second * 10)
    //ticker3 := time.Tick(time.Second * 66)

    //
    go startmining()
    //
    for {
        //fmt.Println("Main loop")
        if Wlt!=nil{
            if Wlt.Walletloaded{ 
                Walletinstantiated=true
            }
        }

        //
        /*
        if  Miningrequested && Walletinstantiated && Wireswarm.Syncingdone{
           
            Mn.SyncWallet(Wlt)
            Mn.LoadUnconfirmedBroadcastedTxs(Wlt)
            go startmining()
             Miningrequested=false
             Miningrunning=true
       
        }*/
        //
        //if  Miningrequested && ((Wlt.HotWallet) || (Walletinstantiated)) && Wireswarm.Syncingdone{
        if Walletinstantiated && Wireswarm.Syncingdone{
           //if !Wlt.HotWallet{
                Mn.SyncWallet(Wlt)
                Mn.LoadUnconfirmedBroadcastedTxs(Wlt)   
           //}
            //if Miningrequested {
                //go startmining()
                //Miningrequested=false
            //}
        }

        
        
        select{
        /////////////////////////////////////////////////

        case <-ticker1://time.After(7 * time.Second):

            if (Wireswarm.NbPeers()>SyncingMinNbPeers) && (!Wireswarm.Syncingdone){
 
                applog.Notice("nb peers %d",Wireswarm.NbPeers())
                Wireswarm.GetPeersMainchainLength()
                longestmainchainpeeraddress:=Wireswarm.GetLongestMainchainPeerAddress(Mn.GetMainchainLength())
                
                if longestmainchainpeeraddress==""{
                    applog.Notice("Syncing stopped - no peer with significantly longer mainchain was found")
                    
                    Wireswarm.Syncingdone=true
                    for _, p := range Wireswarm.Peers {
                        go Wireswarm.ListenPeerMessages(&p) 
                    }
                } else{
                    applog.Notice("longest peer %v",longestmainchainpeeraddress)
                    //os.Exit(2)
                    syncerr:=Wireswarm.InitiateSyncing(Mn,longestmainchainpeeraddress)
                    
                    if syncerr==nil{
                        applog.Trace("Mainchainlength %d Confirmedmainchainlength %d",Mn.GetMainchainLength(),Mn.GetConfirmedMainchainLength())
                        //applog.Notice("Syncing done.")
                        
                        Wireswarm.Syncingdone=true
                        for _, p := range Wireswarm.Peers {
                            go Wireswarm.ListenPeerMessages(&p) 
                        }

                    } else{
                        Wireswarm.RemovePeerByAddress(longestmainchainpeeraddress)
                        //TODO ban peer that caused the error
                    }  
                }

        
            }
        /*
        case <-ticker2 ://time.After(60 * time.Second):
            applog.Trace("Checking missing data")
            hash:=Mn.GetRandomMissingDataHash()
            if hash!=nil {
                Wireswarm.RequestData(*hash)
            } else{
                applog.Trace("No missing data")
            }
        
        case <-ticker3 ://time.After(60 * time.Second):
            applog.Trace("Checking missing data file")
            hash:=Mn.GetRandomMissingDataFileHash()
            if hash!=nil {
                Wireswarm.RequestDataFile(*hash)
            }
        */
        case <-time.After(180 * time.Minute):
            if Walletinstantiated {
                // (re)broadcasting wallet transactions that have not been included in the mainchain
                Mn.SyncWallet(Wlt)
                broadcastingtxs:=Wlt.GetUnconfirmedBroadcastedTxs()
                //fmt.Println("broadcastime",broadcastingtxs)
                for _, broadcastingtx := range broadcastingtxs {
                    Wireswarm.BroadcastTransaction(broadcastingtx)
                }
                //applog.Trace("Wallet ballance %f",float64 (Wlt.ComputeBalance())/1000000.0) 
            }

        case newmsg :=<-Wireswarm.PeersmsgChan:
            applog.Trace("New peer message channel entry %x",newmsg)
            Wireswarm.HandlePeerMessage(Mn,newmsg)
        case newaddr := <-Wireswarm.IpaddrChan:
            applog.Notice("Got new peer address: %s",newaddr)
            go Wireswarm.HintNewPeer(newaddr)
            /*
            newpeer,err:=Wireswarm.HintNewPeer(newaddr)
            if err!=nil{
                applog.Trace("error: %v",err)
            } else {
                Wireswarm.AddPeer(newpeer)
                if Wireswarm.Syncingdone{
                    go Wireswarm.ListenPeerMessages(newpeer,peersmsgChan) 
                }   
            }
            */
        case newpeer := <-Wireswarm.NewpeersChan:
            //
            applog.Notice("New peer connection %v", newpeer.Address)
            //
            Wireswarm.AddPeer(newpeer)
            applog.Notice("nb peers %d",Wireswarm.NbPeers())
            if Wireswarm.Syncingdone{
                go Wireswarm.ListenPeerMessages(newpeer) 
            }   
        /////////////////////////////////////////////////
        case <-quitChan:
            AppIsClosing = true 
            applog.UnlockDisplay()
            applog.Notice("Quitting ...")
            if Walletinstantiated {
                Wlt.SaveJSONWalletFile(MainwalletFilePath,MainwalletFileKey)
            }
            _=SaveUsersettingsFile()
            os.Exit(0)
        }
    }    
}
