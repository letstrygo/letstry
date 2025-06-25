package identifier

import (
	"crypto/rand"
	"encoding/json"
	"math/big"

	"github.com/fatih/color"
)

const (
	idLength   = 6
	charsetLen = 36 // 26 letters + 10 digits
)

func indexToChar(i int64) byte {
	if i < 26 {
		return 'a' + byte(i)
	}
	return '0' + byte(i-26)
}

type ID string

func NewID() ID {
	b := make([]byte, idLength)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(charsetLen))
		b[i] = indexToChar(num.Int64())
	}
	return ID(b)
}

func ParseIDPtr(v string) *ID {
	val := ID(v)
	return &val
}

func (id ID) String() string {
	return string(id)
}

func (id ID) FormattedString() string {
	return color.HiGreenString(id.String())
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
