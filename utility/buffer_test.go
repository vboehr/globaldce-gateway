package utility
import (
	//"github.com/globaldce/globaldce-gateway/utility"
	"testing"
)

func TestEndOfBytes(t *testing.T){
	tmpbw:=NewBufferWriter()
	//tmpbw.PutUint32(1)
	/*
	//x1 := binary.LittleEndian.Uint32(tmpbr.GetContent()[0:4])
	
	//x2 := binary.LittleEndian.Uint16(b[2:])
	
	applog.Notice("% x ", tmpbw.GetContent())
	tmpbr:=utility.NewBufferReader(tmpbw.GetContent())
	applog.Notice("%d", tmpbr.GetUint32())

	applog.Notice("%d", tmpbr.GetUint32())
	if tmpbr.GetError()!=nil{
		applog.Notice("error: %v ",tmpbr.GetError())
	}
	*/
	tmpbw.PutVarUint(160)
	//applog.Trace("%v %d",tmpbw.GetContent(),len(tmpbw.GetContent()))
	tmpbr:=NewBufferReader(tmpbw.GetContent())
	_=tmpbr.GetVarUint()

	
	
	if !tmpbr.EndOfBytes(){
		t.Errorf("Buffer reader error - end of bytes not reached")
	}


}