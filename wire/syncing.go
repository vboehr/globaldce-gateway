package wire

import (
	"fmt"
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/utility"
	"os"
	"path/filepath"
	"time"
)

func (sw *Swarm) GetPeersMainchainLength() {
	//applog.Trace("our Mainchainlength is %d", currentmainchainlength)

	for k, p := range sw.Peers {
		//fmt.Println("k:", k, "v  Add:", v.Connection)
		//if p.SyncingMainchainlength>0 {
		//	applog.Trace("peer mainchainlength %d",p.SyncingMainchainlength)
		//	//return
		//	nbrespondingpeers++
		//	continue
		//}

		// requesting the length of the peer chain
		requestlengthmsg := NewMessage(MsgIdentifierRequestMainchainLength)
		//requestlengthmsg.WriteBytes(p.Connection)
		p.WriteMessage(requestlengthmsg)

		p.Connection.SetReadDeadline(time.Now().Add(30 * time.Second))
		//lengthmsg,rerr:=ReadMessage(p.Connection)
		lengthmsg, rerr := p.ReadMessage()
		if rerr != nil || lengthmsg == nil || !lengthmsg.CheckIdentifier(MsgIdentifierReplyMainchainLength) {
			applog.Trace("error %v - lengthmsg %v", rerr, lengthmsg)

			//return
			continue
		}

		tmpbr := utility.NewBufferReader(lengthmsg.GetContent())

		p.SyncingMainchainlength = tmpbr.GetUint32()
		tmpbrerr := tmpbr.GetError()
		if tmpbrerr != nil {
			p.SyncingMainchainlength = 0
		}
		sw.Peers[k] = p
		applog.Trace("got Mainchainlength %d for %v ", sw.Peers[k].SyncingMainchainlength, k)

	}

}

func (sw *Swarm) GetLongestMainchainPeerAddress(currentmainchainlength uint32) string {
	longestmainchainpeeraddress := ""
	longestmainchain := currentmainchainlength
	for paddress, p := range sw.Peers {
		//applog.Trace("checking peer %v with mainchainlength %d",k, p.SyncingMainchainlength)
		if p.SyncingMainchainlength > currentmainchainlength {
			if p.SyncingMainchainlength > longestmainchain {
				longestmainchainpeeraddress = paddress
				longestmainchain = p.SyncingMainchainlength
			}
		}
		applog.Notice("peer with longest mainchain %v it mainchain length %d our current mainchainlength %d", longestmainchainpeeraddress, longestmainchain, currentmainchainlength)
	}
	return longestmainchainpeeraddress
}

func (sw *Swarm) InitiateSyncing(mn *mainchain.Maincore, lmpeeraddress string) error {

	// requesting and checking all obtained headers
	lmpeermainchainlength := sw.Peers[lmpeeraddress].SyncingMainchainlength
	lmpeer := sw.Peers[lmpeeraddress]

	peermhs := make([]mainchain.Mainheader, 0)
	first := uint32(mn.GetConfirmedMainchainLength())
	last := first + RequestMainheadersMax - 1
	if last > lmpeermainchainlength {
		last = lmpeermainchainlength - 1
	}
	for first < lmpeermainchainlength {
		//time.Sleep(1 * time.Second)
		newmsg := EncodeRequestMainheaders(first, last)
		//applog.Trace("re %x",newmsg.GetContent())
		applog.Trace("sending message %x", newmsg.GetContent())
		//newmsg.WriteBytes(lmpeer.Connection)
		lmpeer.WriteMessage(newmsg)

		//time.Sleep(5 * time.Second)
		//rmsg,rerr:= ReadMessage(lmpeer.Connection)
		rmsg, rerr := lmpeer.ReadMessage()
		if rerr != nil {
			applog.Trace("error", rerr)
			// sync stopping
			return rerr
		}
		applog.Trace("got message %x", rmsg.GetContent())
		///////////////////////////////
		mhs, serr := mn.UnserializeMainheaders(rmsg.GetContent())
		if rerr != nil || mhs == nil {
			applog.Trace("unserialization error", serr)
			// sync stopping
			return serr
		}
		applog.Trace("mainheaders %x", mhs)

		//if !mainchain.CheckHeaderChain(first,mhs){
		//	applog.Trace("mainheaders check error ")
		// sync stopping
		//}
		peermhs = append(peermhs, *mhs...)

		///////////////////////////////

		applog.Trace("first %d last %d", first, last)
		first = last + 1
		last = first + RequestMainheadersMax - uint32(1)
		if last > lmpeermainchainlength-1 {
			last = lmpeermainchainlength - 1
		}
	}
	applog.Trace("got from peer a mainchain of %d length", len(peermhs))
	if !mn.CheckHeaderChain(&peermhs) {
		applog.Warning("mainheaders check error while syncing with %v ", lmpeeraddress)
		// sync stopping
		return fmt.Errorf("mainheaders check error")
	}

	// preparing for sync with peer
	lmpeerpath := filepath.Join(mn.GetPath(), "Sync", lmpeeraddress) //"Tmp/Sync/"+lmpeeraddress+"/"
	os.RemoveAll(filepath.Join(mn.GetPath(), "Sync"))
	os.Mkdir(filepath.Join(mn.GetPath(), "Sync"), os.ModePerm)
	//if _, err := os.Stat(lmpeerpath); os.IsNotExist(err) {
	os.Mkdir(lmpeerpath, os.ModePerm)
	//}
	tmpsyncstoragepath := filepath.Join(lmpeerpath, lmpeeraddress+"_")
	tmpsyncstorage := utility.OpenChunkStorage(tmpsyncstoragepath)

	// requesting and checking all obtained blocks
	startingblockheight := uint32(mn.GetConfirmedMainchainLength())
	applog.Trace("Starting syncing with %d", startingblockheight)
	//lmpeermainchainlength=10
	for requestedblockheight := startingblockheight; requestedblockheight < lmpeermainchainlength; requestedblockheight++ {
		newmsg := EncodeRequestMainblockTransactions(requestedblockheight)
		applog.Trace("sending message %x", newmsg.GetContent())
		applog.Trace("requesting mainblock %d", requestedblockheight)
		//newmsg.WriteBytes( lmpeer.Connection)
		lmpeer.WriteMessage(newmsg)

		//time.Sleep(5 * time.Second)
		lmpeer.Connection.SetReadDeadline(time.Now().Add(5 * time.Second))
		//rmsg,rerr:= ReadMessage( lmpeer.Connection)
		rmsg, rerr := lmpeer.ReadMessage()
		if rerr != nil {
			applog.Trace("error", rerr)
			// sync stopping
		}
		//applog.Trace("got message %x",rmsg.GetContent())
		///////////////////////////////
		mbtxs, serr := mn.UnserializeMainblockTransactions(rmsg.GetContent())
		if rerr != nil {
			applog.Trace("unserialization error", serr)
			// sync stopping
			return serr
		}
		applog.Trace("mainblock transactions %x", mbtxs)
		if !mainchain.CheckMainblockTransactions(mbtxs, peermhs[requestedblockheight-startingblockheight].Roothash) {
			applog.Trace("empty transactions array")
			// sync stopping
			return fmt.Errorf("error: empty transactions array")
		}
		//
		newmb := mainchain.NewMainblock()

		newmb.Height = requestedblockheight
		newmb.Header = peermhs[requestedblockheight-startingblockheight]
		newmb.Transactions = *mbtxs

		applog.Trace("received block %d", newmb.Height)
		//applog.Trace("***** %s ***** %d",newmb.Serialize(),len(newmb.Serialize()))
		tmpsyncstorage.AddChunk(newmb.Serialize())
		//
		//requestedblockheight++
	}
	////////////////////////////////////////////

	////////////////////////////////////////////
	//applog.Trace("position %d size %d file %d", tmpsyncstorage.Chunkposition[0]+4, tmpsyncstorage.Chunksize[0],tmpsyncstorage.Chunkfileid[0])
	//applog.Trace("chunkfileid  %v",tmpsyncstorage.Chunkfileid)

	for i := 0; i < int(lmpeermainchainlength-startingblockheight-mn.GetConfirmationLayer()); i++ {
		//applog.Trace("add to maichain block %d",i+int(startingblockheight))
		//applog.Trace("%v",tmpsyncstorage.GetChunkById(i))
		//applog.Trace("%v", tmpsyncstorage.GetChunkById(i))
		//applog.Trace("position %d size %d file %d", tmpsyncstorage.Chunkposition[i], tmpsyncstorage.Chunksize[i],tmpsyncstorage.Chunkfileid[i])

		tmpblock, _ := mainchain.UnserializeMainblock(tmpsyncstorage.GetChunkById(i))
		//applog.Trace("Validating block%d",i)
		//if mn.ValidateMainblockTransactions(uint32 (i), &tmpblock.Transactions){
		/*
			//TODO rebuild mainstate
			applog.Warning("error: invalide transactions")
			mn.CleanMainstate()
			mn.RebuildMainstate()
			//TODO ban peer
			return fmt.Errorf("error: invalide transactions")
		*/
		//mn.AddBlockChunck(tmpsyncstorage.GetChunkById(i))
		applog.Trace("Confirm mainblock of height %d", tmpblock.Height)
		mn.ConfirmMainblock(tmpblock)

		applog.Trace("Confirm mainblock length %d ", mn.GetMainchainLength())
		//}

	}
	for i := int(lmpeermainchainlength - startingblockheight - mn.GetConfirmationLayer()); i < int(lmpeermainchainlength-startingblockheight); i++ {
		//applog.Trace("add to maichain block %d",i+int(startingblockheight))
		//applog.Trace("%v",tmpsyncstorage.GetChunkById(i))
		//applog.Trace("%v", tmpsyncstorage.GetChunkById(i))
		//applog.Trace("position %d size %d file %d", tmpsyncstorage.Chunkposition[i], tmpsyncstorage.Chunksize[i],tmpsyncstorage.Chunkfileid[i])
		//mn.AddBlockChunck(tmpsyncstorage.GetChunkById(i))
		inmemoryblock, _ := mainchain.UnserializeMainblock(tmpsyncstorage.GetChunkById(i))
		//if mn.ValidateMainblockTransactions(uint32 (i), &inmemoryblock.Transactions){
		mn.AddInMemoryBlock(inmemoryblock)
		applog.Trace("Adding inmemorymainblock %d", inmemoryblock.Height)
		//}
	}
	//applog.Trace("chunkposition %v ",tmpsyncstorage.Chunkposition)

	//applog.Trace("chunk position %v size %v file %v",tmpsyncstorage.Chunkposition,tmpsyncstorage.Chunksize,tmpsyncstorage.Chunkfileid)
	///////////////////////////////
	//os.RemoveAll(tmpsyncstoragepath)
	applog.Notice("syncing done successfully - MainchainLength %d ConfirmedMainchainLength %d", mn.GetMainchainLength(), mn.GetConfirmedMainchainLength())
	return nil
}
