package mq

import "errors"

var (
	ErrNilEncoder = errors.New("nil encoder")
	ErrNilDecoder = errors.New("nil decoder")
)
