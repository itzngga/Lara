package metadata

import (
	"encoding/json"
	"math/big"
)

type StickerMetadata struct {
	Pack      string   `json:"sticker-pack-id"`
	Name      string   `json:"sticker-pack-name"`
	Publisher string   `json:"sticker-pack-publisher"`
	Emojis    []string `json:"emojis"`
}

func writeUIntLE(buffer []byte, value, offset, byteLength int64) {
	slice := make([]byte, byteLength)
	val := new(big.Int)
	val.SetUint64(uint64(value))
	valBytes := val.Bytes()

	tmp := make([]byte, len(valBytes))
	for i := range valBytes {
		tmp[i] = valBytes[len(valBytes)-1-i]
	}
	copy(slice, tmp)
	copy(buffer[offset:], slice)
}

func CreateMetadata(metadata StickerMetadata) []byte {
	bit := make([]byte, 4)
	bit = []byte{0x49, 0x49, 0x2a, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00}

	metadata.Pack = "github.com/itzngga/roxy"
	metadata.Emojis = []string{"😀"}
	jsonData, _ := json.Marshal(metadata)
	bit = append(bit, jsonData...)

	writeUIntLE(bit, int64(len(jsonData)), 14, 4)
	return bit
}
