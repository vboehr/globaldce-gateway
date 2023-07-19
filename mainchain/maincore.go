package mainchain

import (
	"github.com/globaldce/globaldce-gateway/applog"
	//"github.com/globaldce/globaldce-gateway/applog"
	//"github.com/globaldce/globaldce-gateway/wire"
	//"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/utility"
	"github.com/syndtr/goleveldb/leveldb"
	//"github.com/globaldce/globaldce-gateway/mainchain"
	"github.com/globaldce/globaldce-gateway/wallet"
	"os"
	//"math"
	//"time"
	"encoding/json"
	//"math/big"
	//"net"
	//"log"
	"fmt"
	"path/filepath"
	"sync"
)

type Maincore struct {
	MissingContentArray []ContentFileInfo
	//MissingDataFileHashArray []utility.Hash
	//BannedNameArray []string
	path         string
	genesisblock Mainblock
	mainheaders  []Mainheader
	mainbf       *utility.ChunkStorage
	dataf        *utility.ChunkStorage
	mainstatedb  *leveldb.DB
	//
	inmemorymainblocks []Mainblock
	confirmationlayer  uint32
	freezingcoef       uint64
	txspool            *TxsPool
	mu                 sync.Mutex
}

func NewMaincore() *Maincore {
	mn := new(Maincore)
	mn.path = ""
	mn.genesisblock = *GenesisBlock()
	genesisheader := mn.genesisblock.Header
	mn.mainheaders = append(mn.mainheaders, genesisheader)

	mn.confirmationlayer = 6
	mn.freezingcoef = 1000000000000
	mn.txspool = NewTxsPool()
	return mn
}
func (mn *Maincore) Lock() {
	//
	mn.mu.Lock()
}
func (mn *Maincore) Unlock() {
	//
	mn.mu.Unlock()
}
func (mn *Maincore) GetConfirmationLayer() uint32 {
	return mn.confirmationlayer
}

func (mn *Maincore) AddInMemoryBlock(mb *Mainblock) {
	mn.inmemorymainblocks = append(mn.inmemorymainblocks, *mb)
	i := mb.Height
	for j := 0; j < len(mb.Transactions); j++ {
		//applog.Trace("block %d %x ",i,mn.GetMainblock(i).Transactions[j].ComputeHash())
		tx := mb.Transactions[j]
		txhash := tx.ComputeHash()
		mn.txspool.SetTxInclusionHeight(&txhash, i)

		//
	}
}
func (mn *Maincore) PutPath(path string) {
	mn.path = path
}
func (mn *Maincore) GetPath() string {
	return mn.path
}
func (mn *Maincore) AddTransactionToTxsPool(tx *utility.Transaction, fee uint64, priority uint64) {
	mn.txspool.AddTransaction(tx, fee, priority)
}
func (mn *Maincore) LoadUnconfirmedBroadcastedTxs(wlt *wallet.Wallet) {
	txs := wlt.GetUnconfirmedBroadcastedTxs()

	for i := 0; i < len(txs); i++ {
		_, fee := mn.ValidateTransaction(txs[i])
		priority := fee + 1000000
		//applog.Notice("fee %d priority %d ",fee,priority)
		mn.txspool.AddTransaction(txs[i], fee, priority)
	}

}
func (mn *Maincore) LoadMaincore() {
	applog.Notice("Maincore loading ...")
	maincoredirpath := filepath.Join(mn.path, "Mainblocks")
	if _, err := os.Stat(maincoredirpath); os.IsNotExist(err) {
		os.Mkdir(maincoredirpath, os.ModePerm)
		applog.Notice("Creating :%s", maincoredirpath)
	}
	tmpwalletfilesdirpath := filepath.Join(mn.path, "WalletFiles")
	if _, err := os.Stat(tmpwalletfilesdirpath); os.IsNotExist(err) {
		os.Mkdir(tmpwalletfilesdirpath, os.ModePerm)
		applog.Notice("Creating :%s", tmpwalletfilesdirpath)
	}
	if _, err := os.Stat(filepath.Join(mn.path, "Mainblocks", "Mainblocks000")); os.IsNotExist(err) {
		// path does not exist
		mn.mainbf = utility.OpenChunkStorage(filepath.Join(mn.path, "Mainblocks", "Mainblocks"))
		mn.mainbf.AddChunk(mn.genesisblock.Serialize())
		//mn.PutTxState(genesisblock.Header.Hash,0,0)
	} else {
		mn.mainbf = utility.OpenChunkStorage(filepath.Join(mn.path, "Mainblocks", "Mainblocks"))
	}
	////////////////////
	/*
		datadirpath:=filepath.Join(mn.path,"Data")
		if _, err := os.Stat(datadirpath); os.IsNotExist(err) {
			os.Mkdir(datadirpath, os.ModePerm)
		}

		if _, err := os.Stat( filepath.Join(mn.path,"Data","Data000")); os.IsNotExist(err) {
			// path does not exist
			mn.dataf = utility.OpenChunkStorage( filepath.Join(mn.path,"Data","Data"))
			mn.dataf.AddChunk([]byte("emptydata"))
		} else {
			mn.dataf = utility.OpenChunkStorage(filepath.Join(mn.path,"Data","Data"))
		}
	*/
	///////////////////

	for i := 1; i < mn.mainbf.NbChunks(); i++ {

		blockstring := mn.mainbf.GetChunkById(i)
		var block Mainblock
		_ = json.Unmarshal(blockstring, &block)
		mn.mainheaders = append(mn.mainheaders, block.Header)
	}
	//

	var mainstatedbempty = false
	if _, err := os.Stat(filepath.Join(mn.path, "Mainstatedb")); os.IsNotExist(err) {
		mainstatedbempty = true
	}
	msdb, dberr := leveldb.OpenFile(filepath.Join(mn.path, "Mainstatedb"), nil)
	if dberr != nil {
		applog.Warning("error: loading maincore state - %v", dberr)
		return
	}
	mn.mainstatedb = msdb
	if mainstatedbempty {
		applog.Notice("Maincore state empty - Rebuilding Mainstate")
		//TODO
		mn.RebuildMainstate()
	}
	//
	applog.Notice("Maincore loading done.")

	//////////////////////////////

}
func (mn *Maincore) CleanMainstate() {
	mn.mainstatedb.Close()
	os.RemoveAll(filepath.Join(mn.path, "Mainstatedb"))

}
func (mn *Maincore) CleanMainblocks() {
	os.RemoveAll(filepath.Join(mn.path, "Mainblocks"))
	os.Mkdir(filepath.Join(mn.path, "Mainblocks"), os.ModePerm)
}
func (mn *Maincore) GetMainchainLength() uint32 {
	if len(mn.inmemorymainblocks) == 0 {
		return uint32(mn.mainbf.NbChunks())
	} else {
		firstminedblockheight := int(mn.inmemorymainblocks[0].Height)
		tmplength := len(mn.inmemorymainblocks)
		if mn.mainbf.NbChunks() > int(firstminedblockheight+tmplength) {
			return uint32(mn.mainbf.NbChunks())
		} else {
			return uint32(firstminedblockheight + tmplength)
		}

	}
	//return uint32(mn.mainbf.NbChunks())
}
func (mn *Maincore) GetConfirmedMainchainLength() uint32 {
	return uint32(mn.mainbf.NbChunks())
}

func (mn *Maincore) AddBlockChunck(bytes []byte) {
	mn.mainbf.AddChunk(bytes)
}

func (mn *Maincore) UnloadMaincore() {
	mn.mainstatedb.Close()
}
func (mn *Maincore) GetSerializedMainchainLength() []byte {
	tmpbw := utility.NewBufferWriter()
	tmpbw.PutUint32(mn.GetMainchainLength())
	return tmpbw.GetContent()
}
func (mn *Maincore) GetConfirmedMainblock(i int) *Mainblock {
	blockstring := mn.mainbf.GetChunkById(i)
	var block Mainblock
	_ = json.Unmarshal(blockstring, &block)
	return &block
}
func (mn *Maincore) GetLastConfirmedMainblockTimestamp() int64 {
	return mn.GetConfirmedMainblock(mn.mainbf.NbChunks() - 1).Header.Timestamp
}
func (mn *Maincore) GetSerializedMainheaders(first uint32, last uint32) []byte {
	mainchainlength := mn.GetMainchainLength()
	if mainchainlength <= last {
		last = mainchainlength - 1
		applog.Warning("Request last mainheader %d greater than mainchainlength %d", last, mainchainlength)
	}
	tmpbw := utility.NewBufferWriter()

	applog.Trace("\n length %d", len(mn.mainheaders))
	for i := first; i <= last; i++ {
		tmpbytes := mn.GetMainheader(int(i)).Serialize()

		tmpbw.PutUint32(uint32(len(tmpbytes)))
		tmpbw.PutBytes(tmpbytes)
	}
	return tmpbw.GetContent()
}

func (mn *Maincore) UnserializeMainheaders(bytes []byte) (*[]Mainheader, error) {

	var mhs []Mainheader
	byteslength := uint(len(bytes))
	tmpbr := utility.NewBufferReader(bytes)

	for tmpbr.GetCounter() < byteslength {
		mhbyteslength := tmpbr.GetUint32()
		mh, err := UnserializeMainheader(tmpbr.GetBytes(uint(mhbyteslength)))
		if err != nil {
			return nil, err
		}
		mhs = append(mhs, *mh)
	}
	if !tmpbr.EndOfBytes() {
		return nil, fmt.Errorf("End of bytes not reached")
	}

	return &mhs, nil
}

func (mn *Maincore) GetSerializedMainblockTransactions(requestedblockheight uint32) []byte {
	rblock := mn.GetMainblock(int(requestedblockheight)) //TODO GetMainblock int to be changed to uint32
	tmpbw := utility.NewBufferWriter()
	for i := 0; i < len(rblock.Transactions); i++ {

		tmpbytes := rblock.Transactions[i].Serialize()
		tmpbw.PutUint32(uint32(len(tmpbytes)))

		tmpbw.PutBytes(tmpbytes)
	}
	return tmpbw.GetContent()
}
func (mn *Maincore) UnserializeMainblockTransactions(bytes []byte) (*[]utility.Transaction, error) {

	var mbtxs []utility.Transaction
	byteslength := uint(len(bytes))
	tmpbr := utility.NewBufferReader(bytes)

	for tmpbr.GetCounter() < byteslength {
		txbyteslength := tmpbr.GetUint32()
		tx, serr := utility.UnserializeTransaction(tmpbr.GetBytes(uint(txbyteslength)))
		if serr != nil {
			return nil, serr
		}
		mbtxs = append(mbtxs, *tx)
	}
	if !tmpbr.EndOfBytes() {
		return nil, fmt.Errorf("End of bytes not reached")
	}

	return &mbtxs, nil
}
