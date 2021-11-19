package mainchain
import (
	"github.com/globaldce/globaldce/utility"
	"github.com/globaldce/globaldce/applog"
	"math/big"
)
func (mn *Maincore)  AutoCheckMainblocks() bool {
	for i := 1; i < int (mn.GetMainchainLength()); i++ {
		// check the transactions
		mb:=mn.GetMainblock(40)
		if !CheckMainblockTransactions(&(mb.Transactions),mb.Header.Roothash){
			applog.Warning("Rejected Propagating block - incorrect transactions")
			return false
		}
	} 
	return true
}

func (mn *Maincore)  AutoCheckMainheaders() bool {
	headers:=mn.mainheaders

	//checking genesis block
	testgenesisblock:=GenesisBlock()
	if testgenesisblock.Header!=headers[0]{
		applog.Trace("Error corrupt genesis block ")
		return false
	}
	//TODO check inmemoryblocks headers
	for i := 1; i < len(headers); i++ {
		applog.Trace("header %d %x",i,headers[i])
		if (headers[i].Prevblockhash!=headers[i-1].Hash){
			applog.Trace("\n error: blockheader %d - Prevblockhash do not match previous block hash ",i)
			return false
		}
		if (headers[i].Timestamp<headers[i-1].Timestamp){
			applog.Trace("\n error: blockheader %d - Timestamp precede previous block timestamp ",i)
			return false
		}
		//targetbigint:=utility.CorrectTargetBigInt(utility.BigIntFromCompact(headers[i].Bits),headers[i].Timestamp,headers[i-1].Timestamp)
		targetbigint:=utility.BigIntFromCompact(headers[i].Bits)
		if (!(utility.BigIntFromHash(&headers[i].Hash).Cmp(targetbigint)<0)){
			applog.Trace("\n error: hash of header %d do not fall into its own target ",i)
			return false
		}
		//TODO replace code by associated function ??? GetTargetBits()????
		if ((i ) % int (DIFFICULTY_TUNING_INTERVAL)!=0) {
			if (headers[i].Bits!=headers[i-1].Bits){
				applog.Trace("\n error: blockheader %d - Block target do not match previous block target ",i)
				return false
			}
		} else {
			var targetbits uint32
			targetbigint:=utility.BigIntFromCompact(headers[i-1].Bits)
			idealtimeinterval:=int64 (DIFFICULTY_TUNING_INTERVAL-1)*600
			realtimeinterval:=int64 (headers[i-1].Timestamp-headers[i-int (DIFFICULTY_TUNING_INTERVAL)].Timestamp)
			
			if (realtimeinterval>=int64 (3)*idealtimeinterval) {
				targetbigint.Mul(targetbigint,big.NewInt(3))
				//applog.Trace("bigger than 3  ")
			} else if (idealtimeinterval>=int64 (3) *realtimeinterval) {
				targetbigint.Div(targetbigint,big.NewInt(3))
				//applog.Trace("smaller that 1/3 ")
			} else {
				targetbigint.Mul(targetbigint,big.NewInt(realtimeinterval))
				targetbigint.Div(targetbigint,big.NewInt(idealtimeinterval))
			}

			targetbits = utility.CompactFromBigInt(targetbigint)

			applog.Trace("\n realtime %d idealtime %d  targetbitint %d ",realtimeinterval,idealtimeinterval,targetbigint)
			if (targetbits!=headers[i].Bits){
				applog.Trace("\n error: blockheader %d - Block target do not match computed block target ",i)
				return false
			}
		}
	}
	applog.Trace("\n notice: headers verified and found correct ")
	return true
}