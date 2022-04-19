package utility
import
(
	"golang.org/x/crypto/sha3"
	"math/big"
	//"fmt"
	//"github.com/btcsuite/btcd/btcec/v2"
	//"github.com/globaldce/go-globaldce/applog"
)
 


const HashSize=32
type Hash[HashSize]byte
func NewHash(bytes []byte) (*Hash) {
	var nh Hash
	//if len(bytes) != HashSize {
	//	return nil,fmt.Errorf("error: invalid bytes length")
	//}
	copy(nh[:], bytes)
	return &nh
}

func ComputeHash (b []byte) Hash {
	stepone:=sha3.Sum256(b)
	return Hash(sha3.Sum256(stepone[:]))
}
func ComputeHashBytes (b[]byte) []byte {
	stepone:=sha3.Sum256(b)
	steptwo:=sha3.Sum256(stepone[:])
	return (steptwo)[:]
}
/*
func CorrectTargetBigInt(targetbigint *big.Int,timestamp int64,prevtimestamp int64) *big.Int{
	
	obstructioncoeff:=(timestamp-prevtimestamp)/(60*600)//OBSTRUCTED_MINING_TIME// about ten hours
	if obstructioncoeff==0 {
		//obstructioncoeff=1
		return targetbigint
	}

	obstructioncoeffbigint:=big.NewInt(obstructioncoeff)
	targetbigint.Mul(targetbigint,obstructioncoeffbigint)
	
	return targetbigint
}
*/
func BigIntFromHash(h *Hash) *big.Int {

	buf := *h
	blen := len(buf)
	for i := 0; i < blen/2; i++ {
		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}

	return new(big.Int).SetBytes(buf[:])
}

func CompactFromBigInt(value *big.Int) uint32 {
	bigTwo:=big.NewInt(2)
	bigLim:=big.NewInt(16777216)
	var exponent uint32=0
	tvalue:=new(big.Int).Set(value)

	for {
	
		if (tvalue.Cmp(bigLim)==-1){
			mantissa:=uint32(tvalue.Uint64())
			compact := uint32(exponent<<24) | mantissa
			return compact
		}else{
			tvalue.Div(tvalue,bigTwo)
			exponent++
		}
	}


}
func BigIntFromCompact(compact uint32) *big.Int {
	//compact:=0x22fff1ff
	mantissa:=compact&0x00ffffff
	exponent:=uint(compact>>24)
	
	var value * big.Int
	value=big.NewInt(int64(mantissa))
	value.Lsh(value,exponent)
	//applog.Trace(" %x %x",mantissa,exponent)
	return value

}
func ComputeRoot(hashes *[]Hash) Hash{
	var right Hash
	var left Hash
	nodes:= make([]Hash, len(*hashes))
	copy(nodes,*hashes)

	length:=len(nodes)
	for (length>1){
		for i:=0 ; i<length ; i+=2 {
			
			left=nodes[i]
			
			if i+1<length {
				right=nodes[i+1]
			} else {
				right=nodes[i]
			}

			nodes[i/2]=*ComputeHashTreeBranche(&left,&right)
			//applog.Notice("hashes %d %d left %x right %x %x",length,i/2, left,right,computeHashTreeBranche(&left,&right))
		}
		length=length/2+length%2
	}
	return nodes[0]
}

func ComputeHashTreeBranche(left *Hash,right *Hash) *Hash{ 
	var hashconct [HashSize*2]byte 
	copy(hashconct[:HashSize],left[:]) 
	copy(hashconct[HashSize:],right[:])
	result:=ComputeHash(hashconct[:])
	return &result
}
