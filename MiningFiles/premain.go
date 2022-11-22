package main
import (
	"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	//"time"
	//"math/big"
	"github.com/globaldce/globaldce-gateway/utility"
	"github.com/globaldce/globaldce-gateway/mainchain"
	"os"
	"time"
	"encoding/json"
	"math/big"
	"math/rand"
	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/globaldce/globaldce-gateway/daemon"
)

func main() {
	if len(os.Args)<2{
        	fmt.Println("No mining addresses provided")
        	os.Exit(0)
    	}
    	miningaddressesfilepath:=os.Args[1]
    	fmt.Printf("Loading mining addresses file: %s\n",miningaddressesfilepath)
    	maddresses:=new(daemon.MiningAddresses)
    	maddresses.LoadJSONMiningAddressesFile(miningaddressesfilepath)
	//key:=[]byte("aaveryaveryaveryaveryasecretakey")//455s5s6qsqc886sq6q6s6/*éé*&
	applog.Init()
	applog.EnableDisplayTrace()
	
	mn:=mainchain.NewMaincore()
	mn.LoadMaincore()
	defer mn.UnloadMaincore()

	//level db integration
	//mainstatedb, err := leveldb.OpenFile("Mainstate", nil)

	//defer mainstatedb.Close()
	//err = db.Put([]byte("key"), []byte("value"), nil)
	//data, err := db.Get([]byte("key"), nil)
	//applog.Trace("%s%x ",data,err)

	
	//level db integration
	//mainchain.Mainstatedb, _ = leveldb.OpenFile("Mainstatedb", nil)
	//defer mainchain.Mainstatedb.Close()


	genesisblock:=mainchain.GenesisBlock()
	//--------------------------------------------------
	//genesisblock.Mine(genesisblock.Header.Prevblockhash,genesisblock.Header.Bits)
	applog.Trace("MINED genesis block %x ",genesisblock)
	//--------------------------------------------------
	genesisheader:=genesisblock.Header
	mainheaders := make([]mainchain.Mainheader, 0)
	mainheaders=append(mainheaders, genesisheader )

	////////////////////////////////////////////////////////
	//maincoredirpath:=filepath.Join(mn.path,"Mainblocks")
	//if _, err := os.Stat("./Mainblocks"); os.IsNotExist(err) {
	//	os.Mkdir("./Mainblocks", os.ModePerm)
	//}
	////////////////////////////////////////////////////////

	var mainbf * utility.ChunkStorage
	if _, err := os.Stat("Mainblocks/Mainblocks000"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		mainbf = utility.OpenChunkStorage("Mainblocks/Mainblocks")
		mainbf.AddChunk(genesisblock.Serialize())
		//mainchain.PutTxState(genesisblock.Header.Hash,0,0)
	} else {
		mainbf = utility.OpenChunkStorage("Mainblocks/Mainblocks")  
	}

	for i:=1;i<mainbf.NbChunks();i++{

		blockstring:=mainbf.GetChunkById(i)
		var block mainchain.Mainblock
		_=json.Unmarshal(blockstring,&block)
		mainheaders=append(mainheaders, block.Header )
	}

	rand.Seed(time.Now().UnixNano())
	nbiterations:=rand.Intn(200)+200//1096+rand.Intn(1000)
	fmt.Printf(" Nb itrations %d ",nbiterations)

	for i := 1; i < nbiterations; i++ {
		
	
		mainblock:=mainchain.NewMainblock()
		//pubkeyhash:=utility.ComputeHash([]byte("45644444444444455555555555"))
		pubkeyhash:=maddresses.GetRandomAddress()
		testtx:=utility.NewRewardTransaction(mainchain.GENESIS_BLOCK_REWARD,0,pubkeyhash)
		//applog.Trace("transaction %d ",testtx.Vout[0].Value)
		testtxhash:=testtx.ComputeHash()
		mainblock.Header.Roothash=testtxhash
		mainblock.Transactions=append(mainblock.Transactions,*testtx)
		
		var targetbits uint32
		if (mainbf.NbChunks() % int (mainchain.DIFFICULTY_TUNING_INTERVAL)==0){
			targetbigint:=utility.BigIntFromCompact(mainheaders[mainbf.NbChunks()-1].Bits)
			idealtimeinterval:=int64 (mainchain.DIFFICULTY_TUNING_INTERVAL-1)*600
			realtimeinterval:=int64 (mainheaders[mainbf.NbChunks()-1].Timestamp-mainheaders[mainbf.NbChunks()-int (mainchain.DIFFICULTY_TUNING_INTERVAL)].Timestamp)
			
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
			
			targetbits =mainheaders[mainbf.NbChunks()-1].Bits

		}
		mainblock.Height=uint32 (mainbf.NbChunks()+1)
		success:=PreMine(mainblock,mainheaders[mainbf.NbChunks()-1].Timestamp,mainheaders[mainbf.NbChunks()-1].Hash,targetbits)//TODO Replace with mainblockPreMine
		mainheaders=append(mainheaders, mainblock.Header )
		
		
		//wlt.AddAsset(testtxhash,0,mainchain.GENESIS_BLOCK_REWARD,uint32(len(wlt.Privatekeyarray)-1))//=append(wlt.TransactionsHashes, testtxhash)
		//wlt.Lastknownblock=mainblock.Height
		if mainbf.NbChunks()>3 {
		applog.Trace("target %d time between blocks %d success %v",targetbits,mainheaders[mainbf.NbChunks()-1].Timestamp-mainheaders[mainbf.NbChunks()-2].Timestamp,success)
		}
		if success {
			mainbf.AddChunk(mainblock.Serialize())
			//////////////////////////////////////////////
			for j:=0;j<len(mainblock.Transactions);j++{
				//applog.Trace("block %d %x ",i,mn.GetMainblock(i).Transactions[j].ComputeHash())
				i:=mainblock.Height
				tx:=mainblock.Transactions[j]
				txhash:=tx.ComputeHash()
				mn.PutTxState(txhash,uint32(i),uint32(j))
				mn.UpdateMainstate(tx,uint32(i))
				//
			}
			//////////////////////////////////////////////


		} else {
			applog.Trace("Mining unsuccessful")
			break		
		}
		//applog.Trace("[ %d %d ]",i,)//mainheaders[i].Timestamp,mainheaders[mainbf.NbChunks()-1].Timestamp)
		
		
		//mainheaders[i].Timestamp=mainheaders[i-1].Timestamp+600

		if mainblock.Header.Timestamp>=(time.Now().Unix()){
			fmt.Print("Now is the time !....")
			break
		}
	}
	// 
	//applog.Trace("genesis block %x ",genesisblock)
	//applog.Trace("genesis block hash %d target %d ",utility.BigIntFromHash(&genesisblock.Header.Hash),utility.BigIntFromCompact(genesisblock.Header.Bits))

	//genesisblock.Mine(initialhash,genesisblock.Prevblockhash,genesisblock.Bits)
	
	//applog.Trace("mainheaders %x ",mainheaders)
	lastblocktimestamp:=mainheaders[mainbf.NbChunks()-1].Timestamp
	unixTimeUTC:=time.Unix(lastblocktimestamp, 0)
	applog.Trace("\n last block timestamp %d or %v",lastblocktimestamp,unixTimeUTC)



	

}

//////////////////////
/*
func PreMine(mb *Mainblock, prevtime int64, prevblockhash utility.Hash, bits uint32) bool {

	fmt.Println("Minining")
	mb.Header.Bits=bits
	targetbigint:=utility.BigIntFromCompact(bits)
	mb.Header.Prevblockhash=prevblockhash
	starttime:= time.Now().Unix()

	for starttime-time.Now().Unix()<60{//TODO optimize//
	for loopcounter:=0;loopcounter<10000;loopcounter++{//TODO optimize
 	
		mb.Header.Timestamp=int64 (time.Now().Unix())
		
		mb.Header.Nonce++
		mb.ComputeHash()

		//targetbigint=utility.CorrectTargetBigInt(targetbigint,mb.Header.Timestamp,prevtime)
		if (utility.BigIntFromHash(&mb.Header.Hash).Cmp(targetbigint)<0){
			return true
		}
	}
}
	return false
}*/
///////////////////////
func PreMine(mb *mainchain.Mainblock, prevtime int64, prevblockhash utility.Hash, bits uint32) bool {
	rand.Seed(time.Now().UnixNano())//TODO FINALIZE
	relaxationcoef:=int64 (100-(int(mb.Height)/800))//(100-(mb.Height/500))
	if relaxationcoef<1  {
		relaxationcoef=1
	}
	fmt.Printf("height %d relax coef %d  -",mb.Height,relaxationcoef)

	fmt.Println("Minining")
	mb.Header.Bits=bits
	targetbigint:=utility.BigIntFromCompact(bits)
	mb.Header.Prevblockhash=prevblockhash
	starttime:= time.Now().Unix()

	for starttime-time.Now().Unix()<600{//TODO optimize//
	for loopcounter:=0;loopcounter<10000;loopcounter++{//TODO optimize
 	
		//mb.Header.Timestamp=int64 (time.Now().Unix())
		
		mb.Header.Timestamp=relaxationcoef*(( (time.Now().Unix()))-starttime)+prevtime+350+int64 (rand.Intn(100))
		//mb.Header.Timestamp=prevtime+400

		mb.Header.Nonce++
		mb.ComputeHash()

		//targetbigint=utility.CorrectTargetBigInt(targetbigint,mb.Header.Timestamp,prevtime)
		if (utility.BigIntFromHash(&mb.Header.Hash).Cmp(targetbigint)<0){
			fmt.Printf("miningtime %d  -",(time.Now().Unix()-starttime))
			return true
		}


	}
}
	return false
}
