package utility

//import (
//"github.com/globaldce/globaldce-gateway/applog"
//"fmt"
//"encoding/binary"
//"encoding/json"

//)
const (
	ModuleIdentifierECDSATxOut               = 1
	ModuleIdentifierECDSATxIn                = 2
	ModuleIdentifierECDSANameRegistration    = 3
	ModuleIdentifierECDSANameUnregistration  = 4
	ModuleIdentifierECDSARegistredNameCommit = 5
	//ModuleIdentifierECDSAEngagementPublicPost=6
	//ModuleIdentifierECDSAEngagementPublicPostRewardClaim=7
)

//const (
//	EngagementIdentifierLikePublicPost=1
//	EngagementIdentifierDislikePublicPost=2
//)

type TxOut struct {
	Value uint64 // in milion globals
	//Address Hash
	Bytecode []byte
}
type TxIn struct {
	Hash      Hash
	Index     uint32
	Bytecode  []byte
	Signature []byte
}

type Transaction struct {
	Version int32
	//TODO locktime
	//Timestamp int64
	Vin  []TxIn
	Vout []TxOut
}

func NewECDSANameRegistration(Value uint64, Name []byte, Pubkeyhash Hash, CommKeyIdentifier uint32, CommPubkey []byte) TxOut {
	var tmptxout TxOut
	tmptxout.Value = Value
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSANameRegistration) // ECDSA Name Registration Module id
	tmpbw.PutVarUint(uint64(len(Name)))
	//applog.Trace("****** %d",len(Name))
	tmpbw.PutBytes(Name)
	tmpbw.PutHash(Pubkeyhash)
	tmpbw.PutUint32(CommKeyIdentifier)
	if CommKeyIdentifier != CommKeyIdentifierUndefined {
		tmpbw.PutVarUint(uint64(len(CommPubkey)))
		tmpbw.PutBytes(CommPubkey)
	}
	tmpbw.PutVarUint(0) // No extradata
	tmptxout.Bytecode = tmpbw.GetContent()
	return tmptxout
}

func NewECDSATxOut(Value uint64, Pubkeyhash Hash) TxOut {
	var tmptxout TxOut
	tmptxout.Value = Value
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSATxOut) // ECDSA Module id
	tmpbw.PutHash(Pubkeyhash)
	tmpbw.PutVarUint(0) // No extradata
	tmptxout.Bytecode = tmpbw.GetContent()
	return tmptxout
}

func NewECDSATxIn(inhash Hash, index uint32, pubkeycompressedbytes []byte) TxIn {
	var tmptxin TxIn
	tmptxin.Hash = inhash
	tmptxin.Index = index
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSATxIn)
	tmpbw.PutVarUint(uint64(len(pubkeycompressedbytes)))
	tmpbw.PutBytes(pubkeycompressedbytes)
	tmpbw.PutVarUint(uint64(0)) //no extradata
	tmptxin.Bytecode = append(tmptxin.Bytecode, tmpbw.GetContent()...)
	return tmptxin
}

func NewECDSARegistredNameCommit(inhash Hash, index uint32, pubkeycompressedbytes []byte, commitbytes []byte) TxIn {
	var tmptxin TxIn
	tmptxin.Hash = inhash
	tmptxin.Index = index
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSARegistredNameCommit)
	tmpbw.PutVarUint(uint64(len(pubkeycompressedbytes)))
	tmpbw.PutBytes(pubkeycompressedbytes)
	//tmpbw.PutVarUint(uint64 (len(recipientregisteredname))) //recipientregisteredname
	//tmpbw.PutBytes(recipientregisteredname)
	tmpbw.PutVarUint(uint64(len(commitbytes))) //
	tmpbw.PutBytes(commitbytes)
	tmptxin.Bytecode = append(tmptxin.Bytecode, tmpbw.GetContent()...)
	return tmptxin
}

func NewECDSANameUnregistration(inhash Hash, index uint32, pubkeycompressedbytes []byte) TxIn {
	var tmptxin TxIn
	tmptxin.Hash = inhash
	tmptxin.Index = index
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSANameUnregistration)
	tmpbw.PutVarUint(uint64(len(pubkeycompressedbytes)))
	tmpbw.PutBytes(pubkeycompressedbytes)
	tmpbw.PutVarUint(uint64(0)) //no extradata
	tmptxin.Bytecode = append(tmptxin.Bytecode, tmpbw.GetContent()...)
	return tmptxin
}

/*
func NewECDSAEngagementRewardClaim(engagementtxhash Hash,engagementtxindex uint32,pubkeycompressedbytes []byte) TxIn{
	var tmptxin TxIn
	tmptxin.Hash=engagementtxhash
	tmptxin.Index=engagementtxindex
	tmpbw:=NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSAEngagementPublicPostRewardClaim)//
	tmpbw.PutVarUint(uint64(len(pubkeycompressedbytes)))
	tmpbw.PutBytes(pubkeycompressedbytes)
	tmpbw.PutVarUint(uint64 (0)) //no extradata
	tmptxin.Bytecode=append(tmptxin.Bytecode,tmpbw.GetContent()...)
	return tmptxin
}

func NewECDSAEngagementPublicPost(eid uint32,stakedamount uint64,txhash Hash,index uint32,claimaddress Hash) TxOut{
	var tmptxout TxOut
	tmptxout.Value=stakedamount
	tmpbw:=NewBufferWriter()
	tmpbw.PutUint32(ModuleIdentifierECDSAEngagementPublicPost)//
	tmpbw.PutUint32(eid)//EngagementIdentifierLikePublicPost or EngagementIdentifierDislikePublicPost
	//tmpbw.PutVarUint(uint64(len(name)))
	tmpbw.PutHash(txhash)
	tmpbw.PutUint32(index)
	tmpbw.PutHash(claimaddress)
	tmpbw.PutVarUint(0)// No extradata
	tmptxout.Bytecode=tmpbw.GetContent()
	return tmptxout
}
*/
func (txout *TxOut) CompareWithAddress(addr Hash) bool {
	//
	primitivemoduleid := DecodeBytecodeId(txout.Bytecode)
	switch primitivemoduleid {
	case ModuleIdentifierECDSATxOut:
		pubkeyhash, _, err := DecodeECDSATxOutBytecode(txout.Bytecode) // (*Hash,*Extradata,error){
		if err != nil {
			return false
		}
		if *pubkeyhash == addr {
			return true
		}
	case ModuleIdentifierECDSANameRegistration:
		pubkeyhash, _, _, _, err := DecodeECDSANameRegistration(txout.Bytecode)
		if err != nil {
			return false
		}
		if *pubkeyhash == addr {
			return true
		}
	}
	return false
}
func (txout *TxOut) GetAssetState() string {
	//
	primitivemoduleid := DecodeBytecodeId(txout.Bytecode)
	switch primitivemoduleid {
	case ModuleIdentifierECDSATxOut:
		return "UNSPENT"
	case ModuleIdentifierECDSANameRegistration:
		_, name, _, _, _ := DecodeECDSANameRegistration(txout.Bytecode)
		return "NAMEREGISTERED_" + string(name) //
	}
	return "UNKNOWNASSET"
}

func NewRewardTransaction(Value uint64, Fee uint64, Pubkeyhash Hash) *Transaction {
	tx := new(Transaction)
	tx.Version = 1
	tx.Vout = make([]TxOut, 1)
	tx.Vout[0] = NewECDSATxOut(Value+Fee, Pubkeyhash)
	return tx
}
func (tx *Transaction) ComputeSigningHash() (Hash, error) {
	// a signing hash do not include a signature
	tmpbw := NewBufferWriter()
	tmpbw.PutUint32(uint32(tx.Version))
	tmpbw.PutVarUint(uint64(len(tx.Vin)))

	for i := 0; i < len(tx.Vin); i++ {
		tmpbw.PutBytes(tx.Vin[i].Hash[:])
		tmpbw.PutUint32(tx.Vin[i].Index)
		/*
			moduleid:=DecodeBytecodeId(tx.Vin[i].Bytecode)
			if moduleid==ModuleIdentifierECDSATxIn{
				extradata,pubkeycompressedbytes,_,err:=DecodeECDSATxInBytecode(tx.Vin[i].Bytecode)
				if err!=nil{
					return ComputeHash(nil),err
				}

				tmpbytecodebw:=NewBufferWriter()


				tmpbytecodebw.PutVarUint(uint64 (len(pubkeycompressedbytes)))
				tmpbytecodebw.PutBytes(pubkeycompressedbytes)

				signingbytecode:=tmpbytecodebw.GetContent()

				tmpbw.PutVarUint(uint64(len(signingbytecode)))
				tmpbw.PutBytes(signingbytecode)

			} else {*/
		tmpbw.PutVarUint(uint64(len(tx.Vin[i].Bytecode)))
		tmpbw.PutBytes(tx.Vin[i].Bytecode)
		//}
	}
	tmpbw.PutVarUint(uint64(len(tx.Vout)))
	for j := 0; j < len(tx.Vout); j++ {
		tmpbw.PutUint64(uint64(tx.Vout[j].Value))
		//tmpbw.PutBytes(tx.Vout[j].Address[:])
		tmpbw.PutVarUint(uint64(len(tx.Vout[j].Bytecode)))
		tmpbw.PutBytes(tx.Vout[j].Bytecode)
	}

	return ComputeHash(tmpbw.GetContent()), nil

}
func (tx *Transaction) ComputeHash() Hash {
	return ComputeHash(tx.Serialize())
}

/*
func (tx * Transaction) VerifySignatures() bool{
	testtxhash,err:=tx.ComputeSigningHash()//
	if err!=nil{
		return false
	}
	for i:=0;i<len(tx.Vin);i++{
		//tmpbr:=NewBufferReader(tx.Vin[i].Bytecode)
		//primitivemoduleid:=DecodeBytecodeId(tx.Vin[i].Bytecode)
		//if primitivemoduleid==ModuleIdentifierECDSATxIn{
			pubkeycompressedbytes,_,err:=DecodeECDSATxInBytecode(tx.Vin[i].Bytecode)
			signature, err := btcec.ParseSignature(tx.Vin[i].Signature, btcec.S256())
			if err != nil {
				applog.Warning("%v",err)
				return false
			}
			//applog.Trace("signature[%d] %x len %d ",i,tx.Vin[i].Signature,len(tx.Vin[i].Signature))
			pubKey, err := btcec.ParsePubKey(pubkeycompressedbytes, btcec.S256())
			if err != nil {
				applog.Warning("%v",err)
				return false
			}

			verified := signature.Verify(testtxhash[:], pubKey)
			//applog.Trace("Signature %d Verified? %v ",i, verified)
			if !verified {
				return false
			}
		return true
		//}
	}
	return false
} */
