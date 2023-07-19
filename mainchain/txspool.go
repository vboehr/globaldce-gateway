package mainchain

import (
	//"time"
	"container/list"
	"fmt"
	"github.com/globaldce/globaldce-gateway/applog"
	"github.com/globaldce/globaldce-gateway/utility"
)

type TxsPool struct {
	txs                map[utility.Hash]utility.Transaction
	txspriority        map[utility.Hash]uint64
	txsfee             map[utility.Hash]uint64
	txsinclusionheight map[utility.Hash]uint32
	//
	txhashsortedlist *list.List
}

func NewTxsPool() *TxsPool {
	txsp := new(TxsPool)

	txsp.txs = make(map[utility.Hash]utility.Transaction)
	txsp.txspriority = make(map[utility.Hash]uint64)
	txsp.txsfee = make(map[utility.Hash]uint64)
	txsp.txsinclusionheight = make(map[utility.Hash]uint32)
	txsp.txhashsortedlist = list.New()
	return txsp
}

func (txsp *TxsPool) AddTransaction(tx *utility.Transaction, txfee uint64, priority uint64) {
	applog.Notice("Add Transaction with fee %d priority %d", txfee, priority)
	txhash := tx.ComputeHash()
	if _, ok := txsp.txs[txhash]; ok {
		return
	}
	txsp.txs[txhash] = *tx
	txsp.txsfee[txhash] = txfee
	txsp.txspriority[txhash] = priority

	for e := txsp.txhashsortedlist.Front(); e != nil; e = e.Next() {

		if txsp.txspriority[e.Value.(utility.Hash)] < priority {
			txsp.txhashsortedlist.InsertBefore(txhash, e)
			return
		}

	}
	//e:= txsp.txhashsortedlist.Back()
	//txsp.txhashsortedlist.InsertAfter(txhash , e)
	txsp.txhashsortedlist.PushBack(txhash)
	return

	/*
		txsp.txsprioritysorted=append(txsp.txsprioritysorted,testtxhash)
			for i:=0;i<len(txsp.txsprioritysorted);i++{
				if priority> txsp.txspriority[txsp.txsprioritysorted[i]]{
					copy(txsp.txsprioritysorted[i+1:], txsp.txsprioritysorted[i:])
					txsp.txsprioritysorted[i] = testtxhash
					return
				}
			}

		return
	*/
}

func (txsp *TxsPool) DeleteTransaction(txhash *utility.Hash) {
	delete(txsp.txs, *txhash)
	delete(txsp.txsfee, *txhash)
	delete(txsp.txspriority, *txhash)

	for e := txsp.txhashsortedlist.Front(); e != nil; e = e.Next() {

		if e.Value.(utility.Hash) == *txhash {
			txsp.txhashsortedlist.Remove(e)
		}

	}

}
func (txsp *TxsPool) GetTransaction(txhash *utility.Hash) (utility.Transaction, bool) {
	tx, ok := txsp.txs[*txhash]
	return tx, ok
}

func (txsp *TxsPool) SetTxInclusionHeight(txhash *utility.Hash, inclusionheight uint32) {
	//priority,ok:=txsp.txspriority[*txhash]
	//if !ok{
	//	return
	//}
	//priority/=2
	_, ok := txsp.txs[*txhash]
	if ok {
		txsp.txsinclusionheight[*txhash] = inclusionheight
	}

}

func (txsp *TxsPool) GetTransactionPriority(txhash *utility.Hash) uint64 {
	priority, ok := txsp.txspriority[*txhash]
	if !ok {
		return 0
	}
	return priority
}

func (txsp *TxsPool) GetHighestPriorityTxs(mainchainlength uint32) (*[]utility.Transaction, uint64) {

	//for key,priority:=range txsp.txspriority{

	//}
	var totalfees uint64 = 0
	totalseize := 0
	var tmptxs []utility.Transaction
	//var tmphashes []utility.Hash
	for e := txsp.txhashsortedlist.Front(); e != nil; e = e.Next() {

		hash := e.Value.(utility.Hash)
		tx, ok := txsp.txs[hash]
		height, okincluded := txsp.txsinclusionheight[hash]
		//if okincluded && height<mainchainlength{
		//	ok=false
		//}
		fmt.Println("considering priority", txsp.txspriority[e.Value.(utility.Hash)], height, okincluded)
		if ok && (!okincluded || height < mainchainlength) {
			totalseize += len(tx.Serialize())
			if uint32(totalseize) > BLOCKTRANSACTIONS_MAX_SEIZE {
				//fmt.Printf("Block seize exceeded %d\n",totalseize)
				break
			} else {
				tmptxs = append(tmptxs, tx)
				totalfees += txsp.txsfee[hash]
				fmt.Printf("Adding tx with fee %d\n", txsp.txsfee[hash])
				//time.Sleep(5*time.Second)
				//tmphashes=append(tmphashes,hash)
			}

		}
	}
	//fmt.Printf("Block seize %d\n",totalseize)
	//for i:=0;i<len(tmphashes);i++{
	//	txsp.DeleteTransaction(&tmphashes[i])
	//	//fmt.Println("deleted %v",tmphashes[i])
	//}

	return &tmptxs, totalfees
}
func (txsp *TxsPool) DisplayTxs() {
	for e := txsp.txhashsortedlist.Front(); e != nil; e = e.Next() {
		fmt.Println("priority", txsp.txspriority[e.Value.(utility.Hash)], e.Value.(utility.Hash))
	}
	fmt.Println(" length ", txsp.txhashsortedlist.Len())
}
