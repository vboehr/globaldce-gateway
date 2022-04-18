package utility

import (
	//"github.com/globaldce/go-globaldce/applog"
	"fmt"
	"encoding/binary"
	"bytes"
	"log"
	"os"
	"io"
)

const ChunkFileMaxSize=20*1024*1024//100*1024*1024

// ChunkStorage is
type ChunkStorage struct {
	Path string
	file [] *os.File
	//
	Chunkposition []int64
	Chunksize []int64
	Chunkfileid []int
}
func (cs *ChunkStorage) NbChunks() int{
	return int(len(cs.Chunkposition))
}
// OpenChunkStorage is
func OpenChunkStorage(storagepath string) *ChunkStorage {
	cs := new(ChunkStorage)
	cs.Path = storagepath

	var activechunkfileid int=0
for {
	
	filepath:=fmt.Sprintf("%s%03d",storagepath,activechunkfileid)
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	cs.file =append(cs.file,f)
	
	//--------------------------

	var position int64 =-4
	var chunksize uint32 =uint32 (0)
	var i int =0

	for  {
		position+=int64(chunksize+4)
		_, seekerr := cs.file[activechunkfileid].Seek(position, 0)
		if seekerr != nil {
			cs.file[activechunkfileid].Close() // ignore error; Write error takes precedence
			log.Fatal(seekerr)
		}
		//applog.Trace("position %d", position)

		bufferchunksize := make([]byte, 4)
		_, readerr := cs.file[activechunkfileid].Read(bufferchunksize)

		if readerr == io.EOF {
			break
		}

		if readerr != nil {
			cs.file[activechunkfileid].Close() // ignore error; Write error takes precedence
			log.Fatal(readerr)
		}

		
		readerchunksize := bytes.NewReader(bufferchunksize)


		binary.Read(readerchunksize, binary.LittleEndian, &chunksize)

		cs.Chunksize=append(cs.Chunksize ,int64 (chunksize))
		cs.Chunkposition=append(cs.Chunkposition , int64 (position))
		cs.Chunkfileid=append(cs.Chunkfileid,activechunkfileid)

		i++
	}
	//---------------------------
	if _, err := os.Stat(fmt.Sprintf("%s%03d",storagepath,activechunkfileid+1)); os.IsNotExist(err) {
		break
	}else{
		activechunkfileid++
	}
	//---------------------------
}

	//--------------------------
	return cs
}

func (cs *ChunkStorage) AddChunk(data []byte) error {
	var activechunkfileid int
	activechunkfileid=len(cs.file)-1
	var newchunkfile bool =false


	//------------------------------
	fileinfo, staterr := cs.file[activechunkfileid].Stat()
	if staterr != nil {
		log.Fatal(staterr)
	}
	//applog.Trace("size %d", fileinfo.Size())
	if (fileinfo.Size()>ChunkFileMaxSize){
		activechunkfileid++
		filepath:=fmt.Sprintf("%s%03d",cs.Path,activechunkfileid)
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		cs.file =append(cs.file,f)
		newchunkfile=true
	}

	//------------------------------


	bufferchunkfilesize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bufferchunkfilesize, uint32(len(data)))
	_,serrs:=cs.file[activechunkfileid].Seek(0,os.SEEK_END)
	if serrs != nil {
		cs.file[activechunkfileid].Close() // error
		log.Fatal(serrs)
	}
	_, lerr := cs.file[activechunkfileid].Write(bufferchunkfilesize)
	if lerr != nil {
		cs.file[activechunkfileid].Close() // error
		log.Fatal(lerr)
	}
	_,serrd:=cs.file[activechunkfileid].Seek(0,os.SEEK_END)
	if serrd != nil {
		cs.file[activechunkfileid].Close() // error
		log.Fatal(serrd)
	}

	_, err := cs.file[activechunkfileid].Write(data)
	if err != nil {
		cs.file[activechunkfileid].Close() // error
		log.Fatal(err)
	}

	if (!newchunkfile) && (cs.NbChunks()>=1){
			newposition:=int64 (  int (cs.Chunkposition[cs.NbChunks()-1]+cs.Chunksize[cs.NbChunks()-1])+4 )
			cs.Chunkposition=append(cs.Chunkposition , newposition)	
		} else {
			cs.Chunkposition=append(cs.Chunkposition, int64(0))
		}






	cs.Chunksize=append(cs.Chunksize ,int64 (len(data)))

	cs.Chunkfileid=append(cs.Chunkfileid ,activechunkfileid)

	return err
}


func (cs *ChunkStorage) GetChunkById(chunkid int) []byte {


	//applog.Trace("position %d size %d file %d", cs.Chunkposition[chunkid]+4, cs.Chunksize[chunkid],cs.Chunkfileid[chunkid])
	return cs.GetChunk(cs.Chunkposition[chunkid]+4, cs.Chunksize[chunkid],cs.Chunkfileid[chunkid])
}


func (cs *ChunkStorage) GetChunk(position int64, length int64,fileid int) []byte {
	_, seekerr := cs.file[fileid].Seek(position, 0)
	//check(err)
	if seekerr != nil {
		cs.file[fileid].Close() // ignore error; Write error takes precedence
		log.Fatal(seekerr)
	}
	chunk := make([]byte, length)
	_, readerr := cs.file[fileid].Read(chunk)

	if readerr != nil {
		cs.file[fileid].Close() // ignore error; Write error takes precedence
		log.Fatal(readerr)
	}
	return chunk
}

