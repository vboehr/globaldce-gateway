package wallet
import
(
	//"fmt"
	"math/rand"
	"time"
)
func GenerateRandomSeedString() string{

	rand.Seed(time.Now().UnixNano())
	//fmt.Println("***",rand.Intn(5000))
	n:=uint(50+rand.Intn(100))
	//fmt.Println("n:",n)
	var letterRunes = []rune("adefghijklqrstvxyzABCDEFGHIJKLMNPQRSTUVWXYZ0123456789*+!?%#")

    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
} 
