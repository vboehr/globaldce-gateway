package utility

import (
	"fmt"
	"golang.org/x/term"
	"syscall"
)

func PrintBytes(b []byte) {
	for i := 0; i < len(b); i++ {
		fmt.Printf("0x%x, ", b[i])
	}

	fmt.Printf("")
}

/*
func main() {
	fmt.Print("Enter Password: ")
	password, _ := GetUserSecretEntry()
	fmt.Printf(" Password: %s\n", password)
}
*/
// go get golang.org/x/term
func GetUserSecretEntry() (string, error) {
	byteUserSecretEntry, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	UserSecretEntry := string(byteUserSecretEntry)
	return UserSecretEntry, nil
}
