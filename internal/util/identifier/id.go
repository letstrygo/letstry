package identifier

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
)

const (
	idLength   = 6
	charset    = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetLen = int64(len(charset))
)

type ID string

func NewID() ID {
	b := make([]byte, idLength)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(charsetLen))
		b[i] = charset[num.Int64()]
	}
	return ID(b)
}

func (id ID) String() string {
	return string(id)
}

func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*id = ID(s)
	return nil
}
