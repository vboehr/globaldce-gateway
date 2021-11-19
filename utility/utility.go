package utility
import 	"fmt"

func PrintBytes (b []byte) {
	for i:=0;i<len(b);i++{
		fmt.Printf("0x%x, ",b[i])
	}
	
	fmt.Printf("")
}