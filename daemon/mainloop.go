package daemon

import (
    "github.com/globaldce/globaldce/applog"
    "os"
    "os/signal"
    "time"
    "github.com/globaldce/globaldce/mainchain"
    "github.com/globaldce/globaldce/wire"
    "log"
    "fmt"
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

        for  Wireswarm.Syncingdone {
              fmt.Println("Working..")
                success,mb:=Mn.Mine(Wlt)
                if success {
                    Wireswarm.BroadcastMainblock(mb)
                }
                Mn.SyncWallet(Wlt)
            //applog.Trace("mainchaillength %d",Mn.GetConfirmedMainchainLength())
        }

}

/////////////////////////////
func Mainloop(){


    applog.Init()
    quitChan := listenSigInt()
    
    Mn=mainchain.NewMaincore()
    Mn.PutPath( AppPath)
	Mn.LoadMaincore()
	defer Mn.UnloadMaincore()

    Wireswarm=wire.NewSwarm()

    if Seed {
        Wireswarm.Syncingdone=true
    }


	
	
	
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

    //
    for {

        //

        if  Miningrequested && Walletloaded && Wireswarm.Syncingdone{
           
            Mn.SyncWallet(Wlt)
            Mn.LoadUnconfirmedBroadcastedTxs(Wlt)
            go startmining()
             Miningrequested=false
             Miningrunning=true
       
        }

        
        
        select{
        /////////////////////////////////////////////////
        
        case <-time.After(7 * time.Second):

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
        case <-time.After(180 * time.Minute):
            // (re)broadcasting wallet transactions that have not been included in the mainchain
            Mn.SyncWallet(Wlt)
            broadcastingtxs:=Wlt.GetUnconfirmedBroadcastedTxs()
            //fmt.Println("broadcastime",broadcastingtxs)
            for _, broadcastingtx := range broadcastingtxs {
                Wireswarm.BroadcastTransaction(broadcastingtx)
            }
            //applog.Trace("Wallet ballance %f",float64 (Wlt.ComputeBalance())/1000000.0)
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
            Wlt.SaveJSONFile(MainwalletFilePath,MainwalletFileKey)
            os.Exit(0)
        }
    }    
}
