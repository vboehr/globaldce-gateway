package wallet

//import
//(
//
//"encoding/json"
//"encoding/binary"
//"bytes"

//"github.com/globaldce/globaldce-gateway/applog"
//"fmt"
//"os"
//"sync"
//)
import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	//"encoding/base64"
	"fmt"
	//"crypto/rand"
	//"crypto/rsa"
	"crypto/x509"
	//"encoding/pem"
	//"fmt"
	//"io/ioutil"
	"github.com/globaldce/globaldce-gateway/utility"
)

type Commcredential struct {
	Name             string
	RSAPrivateKey    rsa.PrivateKey
	RSAPublicKeyHash utility.Hash
}

func (wlt *Wallet) GenerateCommKey(tmpname []byte) []byte {
	// The GenerateKey method takes in a reader that returns random bits, and
	// the number of bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Error:%v", err)
		return nil
	}

	// The public key is a part of the *rsa.PrivateKey struct
	//publicKey := privateKey.PublicKey
	tmpPubkeyBytes, serr := SerializeRSAPublicKey(&privateKey.PublicKey)
	if serr != nil {
		fmt.Printf("Error:%v", serr)
		return nil
	}
	tmpPubkeyhash := utility.ComputeHash(tmpPubkeyBytes)

	tmpCommcredential := Commcredential{
		Name:             string(tmpname),
		RSAPrivateKey:    *privateKey,
		RSAPublicKeyHash: tmpPubkeyhash,
	}
	wlt.Commcredentialarray = append(wlt.Commcredentialarray, tmpCommcredential)
	return tmpPubkeyBytes
}

func (wlt *Wallet) EncryptCommText(tmptext []byte, recipientPubkeyBytes []byte, sendername []byte) ([]byte, []byte, error) {

	var senderRSAPrivKeyIndex int

	for i := (len(wlt.Commcredentialarray) - 1); i <= 0; i++ { //i,tmpCommcredential:=range wlt.Commcredentialarray {
		if wlt.Commcredentialarray[i].Name == string(sendername) {
			senderRSAPrivKeyIndex = i
		}
	}

	recipientPubkey, userr := UnserializeRSAPublicKey(recipientPubkeyBytes)
	if userr != nil {
		return nil, nil, userr
	}
	encryptedBytes, eerr := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&recipientPubkey, //wlt.Commcredentialarray[selectedPrivKeyIndex],
		tmptext,
		nil)
	if eerr != nil {
		return nil, nil, eerr
	}

	//fmt.Println("encrypted bytes: ", encryptedBytes)
	//
	//txt := []byte("verifiable message")

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	encryptedBytesHash := sha256.New()
	_, eherr := encryptedBytesHash.Write(encryptedBytes) //(txt)
	if eherr != nil {
		return nil, nil, eherr
	}
	encryptedBytesHashBytes := encryptedBytesHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	senderprivateKey := wlt.Commcredentialarray[senderRSAPrivKeyIndex].RSAPrivateKey
	signature, serr := rsa.SignPSS(rand.Reader, &senderprivateKey, crypto.SHA256, encryptedBytesHashBytes, nil)
	if serr != nil {
		return nil, nil, serr
	}
	/////////////////////
	//
	return encryptedBytes, signature, nil
}
func VerifyRSASignature(encryptedBytes []byte, sendersignature []byte, senderpublickeyBytes []byte) error {
	//
	senderpublickey, userr := UnserializeRSAPublicKey(senderpublickeyBytes)
	_ = userr

	encryptedBytesHash := sha256.New()
	_, eherr := encryptedBytesHash.Write(encryptedBytes) //(txt)
	if eherr != nil {
		return eherr
	}
	encryptedBytesHashBytes := encryptedBytesHash.Sum(nil)
	verr := rsa.VerifyPSS(&senderpublickey, crypto.SHA256, encryptedBytesHashBytes, sendersignature, nil)
	if verr != nil {
		fmt.Println("could not verify signature: ", verr)
		return verr
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
	////////////////////////////////
	return nil
}
func (wlt *Wallet) DecryptCommText(encryptedBytes []byte, recepientpublickeyhash utility.Hash) ([]byte, error) {
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now

	/*
		pubkeyBytes := x509.MarshalPKCS1PublicKey(&publicKey)
		privkeyBytes :=x509.MarshalPKCS1PrivateKey(privateKey)
		//	fmt.Printf("privateKey %s\n", exportPrivKeyAsPEMStr(privateKey))
		//	fmt.Printf("publicKey %s\n", exportPubKeyAsPEMStr(&publicKey))

			fmt.Printf("privateKey len %d\n", len(privkeyBytes))
			fmt.Printf("publicKey len %d\n", len(pubkeyBytes))

		pubkeyParsed,perr:= x509.ParsePKCS1PublicKey(pubkeyBytes)// (*rsa.PublicKey, error)
		if perr != nil {
			panic(perr)
		}
		fmt.Printf("orignal%v == parsed %v\n",publicKey,*pubkeyParsed)
		fmt.Printf("**********\n")
		//func ParsePKCS1PrivateKey(der []byte) (*rsa.PrivateKey, error)
		privkeyParsed,pperr:= x509.ParsePKCS1PrivateKey(privkeyBytes)// (*rsa.PublicKey, error)
		if pperr != nil {
			panic(pperr)
		}
		fmt.Printf("orignal%v == parsed %v\n",*privateKey,*privkeyParsed)
		//
	*/
	//
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OEAPOptions in the end signify that we encrypted the data using OEAP, and that we used
	// SHA256 to hash the input.
	var recepientRSAPrivKeyIndex int
	for i := (len(wlt.Commcredentialarray) - 1); i <= 0; i++ { //i,tmpCommcredential:=range wlt.Commcredentialarray {
		if wlt.Commcredentialarray[i].RSAPublicKeyHash == recepientpublickeyhash {
			recepientRSAPrivKeyIndex = i
		}
	}
	recepientprivateKey := wlt.Commcredentialarray[recepientRSAPrivKeyIndex].RSAPrivateKey
	decryptedBytes, derr := recepientprivateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if derr != nil {
		return nil, derr
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	fmt.Println("decrypted message: ", string(decryptedBytes))
	return decryptedBytes, nil
}
func SerializeRSAPrivateKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	return privkeyBytes, nil
}
func SerializeRSAPublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	pubkeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	return pubkeyBytes, nil
}
func UnserializeRSAPrivateKey(privkeyBytes []byte) (rsa.PrivateKey, error) {
	privkeyParsed, err := x509.ParsePKCS1PrivateKey(privkeyBytes) // (*rsa.PublicKey, error)
	return *privkeyParsed, err
}
func UnserializeRSAPublicKey(pubkeyBytes []byte) (rsa.PublicKey, error) {
	pubkeyParsed, err := x509.ParsePKCS1PublicKey(pubkeyBytes) // (*rsa.PublicKey, error)
	return *pubkeyParsed, err
}
