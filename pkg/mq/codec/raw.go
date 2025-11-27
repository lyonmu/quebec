package codec

import (
	"bytes"
	"github.com/lyonmu/quebec/pkg/mq"
)

type rawEncoder struct{}

func NewRawEncoder() mq.Encoder[[]byte, []byte] {
	return &rawEncoder{}
}

func (e *rawEncoder) EncodeKey(key *[]byte) ([]byte, error) {
	return *key, nil
}

func (e *rawEncoder) EncodeValue(value *[]byte) ([]byte, error) {
	return *value, nil
}

type rawDecoder struct {
}

func NewRawDecoder() mq.Decoder[bytes.Buffer, bytes.Buffer] {
	return &rawDecoder{}
}

func (d *rawDecoder) DecodeKey(key []byte, dst *bytes.Buffer) error {
	if len(key) > 0 {
		dst.Write(key)
	}
	return nil
}

func (d *rawDecoder) DecodeValue(value []byte, dst *bytes.Buffer) error {
	if len(value) > 0 {
		dst.Write(value)
	}

	return nil
}
