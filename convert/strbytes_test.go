package convert

import (
	"bytes"
	"testing"
)

var materials = []struct {
	str string
	bys []byte
}{
	{"a", []byte{97}},
	{"A", []byte{65}},
	{"Hello Zoe", []byte{72, 101, 108, 108, 111, 32, 90, 111, 101}},
	{"你好，周筱鲁", []byte{228, 189, 160, 229, 165, 189, 239, 188, 140, 229, 145, 168, 231, 173, 177, 233, 178, 129}},
}

func TestBytes2String(t *testing.T) {
	for _, m := range materials {
		if Bytes2String(m.bys) != m.str {
			t.Error("not equals!")
		}
	}
}

func TestString2Bytes(t *testing.T) {
	for _, m := range materials {
		if !bytes.Equal(String2Bytes(m.str), m.bys) {
			t.Error("not equals!")
		}
	}
}

func BenchmarkBytes2String(b *testing.B) {

}

func BenchmarkString2Bytes(b *testing.B) {

}
