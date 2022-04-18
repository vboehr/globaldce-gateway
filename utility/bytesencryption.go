package utility

import (
	//"github.com/globaldce/go-globaldce/applog"
	//"encoding/binary"
	//"bytes"
	//"os"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    //"log"
)
/*
func LoadJSONFile(path string) *[]byte{
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("error:", err)
	}
	bufferWalletfileseize := make([]byte, 4)
	_, rserr := f.Read(bufferWalletfileseize)
	if rserr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rserr)
	}
	var walletfileseize uint32
	readerWalletfileseize := bytes.NewReader(bufferWalletfileseize)

	binary.Read(readerWalletfileseize, binary.LittleEndian, &walletfileseize)
	walletfilestring := make([]byte, walletfileseize)
	_, rerr := f.Read(walletfilestring)
	if rerr != nil {
		//log.Fatal(err)
		fmt.Println("error:", rerr)
	}
	return &walletfilestring
	
}
*/

// The following code was created with the help of Stack Overflow question
// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64
// As answered by the following user:
// https://stackoverflow.com/users/1567738/intermernet
// The answer was also modified by the following user: 
// https://stackoverflow.com/users/2588732/roundsparrow-hilltx

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	//applog.Trace("iv:%x", iv)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	//applog.Trace("iv:%x", iv)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}