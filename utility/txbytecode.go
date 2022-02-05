package utility

import (
	"fmt"
)


const RegistredNameMaxSize=20
func DecodeBytecodeId(bytecode []byte) (uint32){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return 0
	}
	return primitivemoduleid
}

func DecodeECDSATxOutBytecode(bytecode []byte) (*Hash,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSATxOut{
		return nil,nil,fmt.Errorf("Not an ECDSA Vout Module")
	}

	txoutaddr:=tmpbr.GetHash()
	extradata:=tmpbr.GetExtradata()
	/*
	extradatalen:=tmpbr.GetVarUint()
	if extradatalen>ExtradataMaxSize{
		return nil,nil,fmt.Errorf("ExtradataMaxSize exceeded")
	}
	var extradata Extradata
	if extradatalen!=0{
		extradata.Size=extradatalen
		extradata.Hash=tmpbr.GetHash()
	}*/
	//if tmpbr.GetError()!=nil{
	//	return false
	//}
	//return true
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,fmt.Errorf("End of bytes not reached")
	}
	return &txoutaddr,extradata,nil

}

func DecodeECDSATxInBytecode(bytecode []byte) ([]byte,*Extradata,error ){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSATxIn{
		return nil,nil,fmt.Errorf("Not an ECDSA Vin Module")
	}
	
	pubkeycompressedlen:=tmpbr.GetVarUint()
	pubkeycompressed:=tmpbr.GetBytes(uint(pubkeycompressedlen))
	extradata:=tmpbr.GetExtradata()
	/*
	extradatalen:=tmpbr.GetVarUint()
	if extradatalen>ExtradataMaxSize{
	return nil,nil,fmt.Errorf("ExtradataMaxSize exceeded")
	}

	var extradata Extradata
	if extradatalen!=0{
		extradata.Size=extradatalen
		extradata.Hash=tmpbr.GetHash()
	}
	*/
	//signaturelen:=tmpbr.GetVarUint()
	//signature:=tmpbr.GetBytes(uint(signaturelen))
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return pubkeycompressed,extradata,nil
}
//
func DecodeECDSANamePublicPost(bytecode []byte) ([]byte,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSANamePublicPost{
		return nil,nil,fmt.Errorf("Not an ECDSA Name Unregistration")
	}
	pubkeycompressedlen:=tmpbr.GetVarUint()
	pubkeycompressed:=tmpbr.GetBytes(uint(pubkeycompressedlen))

	ed:=tmpbr.GetExtradata()
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return pubkeycompressed,ed,nil
}

func DecodeECDSANameUnregistration(bytecode []byte) ([]byte,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSANameUnregistration{
		return nil,nil,fmt.Errorf("Not an ECDSA Name Unregistration")
	}

	pubkeycompressedlen:=tmpbr.GetVarUint()
	pubkeycompressed:=tmpbr.GetBytes(uint(pubkeycompressedlen))
	/*
	extradatalen:=tmpbr.GetVarUint()
	if extradatalen>ExtradataMaxSize{
		return nil,nil,fmt.Errorf("ExtradataMaxSize exceeded")
	}
	var extradata Extradata
	if extradatalen!=0{
		extradata.Size=extradatalen
		extradata.Hash=tmpbr.GetHash()
	}
	*/
	extradata:=tmpbr.GetExtradata()
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return pubkeycompressed,extradata,nil
}
func DecodeECDSANameRegistration(bytecode []byte) (*Hash,[]byte,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSANameRegistration{
		return nil,nil,nil,fmt.Errorf("Not an ECDSA Name Unregistration")
	}

	namelen:=tmpbr.GetVarUint()
	if namelen>RegistredNameMaxSize{
		return nil,nil,nil,fmt.Errorf("Name in ECDSANameRegistration is too long - %d",namelen)
	}
	name:=tmpbr.GetBytes(uint(namelen)) 
	pubkeyhash:=tmpbr.GetHash()
	/*
	extradatalen:=tmpbr.GetVarUint()
	if extradatalen>ExtradataMaxSize{
		return nil,nil,nil,fmt.Errorf("ExtradataMaxSize exceeded")
	}
	var extradata Extradata
	if extradatalen!=0{
		extradata.Size=extradatalen
		extradata.Hash=tmpbr.GetHash()
	}
	*/
	extradata:=tmpbr.GetExtradata()
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return &pubkeyhash,name,extradata,nil

}
//
func DecodeECDSAEngagement(bytecode []byte) (uint32,*Hash,uint32,*Hash,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSAEngagementPublicPost{
		return 0,nil,0,nil,nil,fmt.Errorf("Not an Engagement bytecode")
	}
	eid:=tmpbr.GetUint32()
	//namelen:=tmpbr.GetVarUint()
	//if namelen>RegistredNameMaxSize{
	//	return 0,nil,0,nil,nil,fmt.Errorf("Name in Engagement bytecode is too long - %d",namelen)
	//}
	//name:=tmpbr.GetBytes(uint(namelen)) 
	
	/*
	extradatalen:=tmpbr.GetVarUint()
	if extradatalen>ExtradataMaxSize{
		return nil,nil,nil,fmt.Errorf("ExtradataMaxSize exceeded")
	}
	var extradata Extradata
	if extradatalen!=0{
		extradata.Size=extradatalen
		extradata.Hash=tmpbr.GetHash()
	}
	*/
	hash:=tmpbr.GetHash()
	index:=tmpbr.GetUint32()
	claimaddress:=tmpbr.GetHash()
	extradata:=tmpbr.GetExtradata()
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return 0,nil,0,nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return 0,nil,0,nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return eid,&hash,index,&claimaddress,extradata,nil

}
func DecodeECDSAEngagementRewardClaim (bytecode []byte) ([]byte,*Extradata,error){
	tmpbr:=NewBufferReader(bytecode)
	primitivemoduleid:=tmpbr.GetUint32()
	if primitivemoduleid != ModuleIdentifierECDSAEngagementPublicPostRewardClaim {
		return nil,nil,fmt.Errorf("Not an ECDSA EngagementPublicPostRewardClaim")
	}
	pubkeycompressedlen:=tmpbr.GetVarUint()
	pubkeycompressed:=tmpbr.GetBytes(uint(pubkeycompressedlen))

	ed:=tmpbr.GetExtradata()
	tmpbrerr:=tmpbr.GetError()
	if tmpbrerr!=nil{
		return nil,nil,tmpbrerr
	}
	if !tmpbr.EndOfBytes(){
		return nil,nil,fmt.Errorf("End of bytes not reached")
	}
	
	return pubkeycompressed,ed,nil
}
