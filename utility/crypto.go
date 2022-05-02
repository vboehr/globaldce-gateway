package utility
import
(
	//"github.com/globaldce/globaldce-gateway/utility"
	//"encoding/json"
	//"encoding/binary"
	//"bytes"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/globaldce/globaldce-gateway/applog"
	"fmt"
	//"error"
	//"os"
	//"sync"
)

func NewPrivateKey() (btcec.PrivateKey,error) {
	pk,err:=btcec.NewPrivateKey()
	return *pk,err
}
func Sign(pk *btcec.PrivateKey,msgbytes []byte) ([]byte) {
	sig:=ecdsa.Sign(pk,msgbytes)
	return sig.Serialize()
}
func VerifySignature(signinghash Hash,signaturebytes []byte,pubkeycompressedbytes []byte) (error){
	//pubkeycompressedbytes,_,err:=DecodeECDSATxInBytecode(tx.Vin[i].Bytecode)
	signature, err := ecdsa.ParseSignature(signaturebytes) 
	if err != nil {
		applog.Warning("%v",err)
		return err
	}
	//applog.Trace("signature[%d] %x len %d ",i,tx.Vin[i].Signature,len(tx.Vin[i].Signature))
	pubKey, err := btcec.ParsePubKey(pubkeycompressedbytes)
	if err != nil {
		applog.Warning("%v",err)
		return err
	}
	
	verified := signature.Verify(signinghash[:], pubKey)
	//applog.Trace("Signature %d Verified? %v ",i, verified)	
	if !verified {
		return fmt.Errorf("Signature not verified")
	}
	return nil
}


func PrivKeyFromBytes(privateKeyBytes[]byte) (btcec.PrivateKey,btcec.PublicKey) {
	pk,pubkey:=btcec.PrivKeyFromBytes(privateKeyBytes)
	return *pk,*pubkey
}
