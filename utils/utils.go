package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

// when err isn't nil
// then print err
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// inteface -> byte slice
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

// set interface data (byte) -> decoded data
func FromBytes(i interface{}, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}
