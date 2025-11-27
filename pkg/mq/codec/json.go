package codec

import (
	"bytes"
	"io"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/lyonmu/quebec/pkg/mq"
)

var (
	DefaultSize = 1 * 1024 * 1024
	MaxSize     = 10 * 1024 * 1024
)

var poolJson = &sync.Pool{
	New: func() any {
		return NewJson()
	},
}

func AcquireJson() *Json {
	return poolJson.Get().(*Json)
}

func ReleaseJson(j *Json) {
	j.Reset()
	poolJson.Put(j)
}

type Json struct {
	defaultSize int
	maxSize     int
	len         int
	data        []byte
	reader      *bytes.Reader
	Encoder     sonic.Encoder
}

func NewJson() *Json {
	res := &Json{
		data: make([]byte, DefaultSize),
	}
	res.reader = bytes.NewReader(res.data)
	res.Encoder = sonic.ConfigDefault.NewEncoder(res)
	return res
}

func (j *Json) Marshal(v interface{}) error {
	err := j.Encoder.Encode(v)
	if err != nil {
		j.Encoder = sonic.ConfigDefault.NewEncoder(j)
		err = j.Encoder.Encode(v)
	}
	if err == nil {
		j.ResetReader()
	}
	return err
}
func (j *Json) Unmarshal(v interface{}) error {
	return sonic.ConfigDefault.NewDecoder(j.reader).Decode(v)
}

func (j *Json) Write(p []byte) (n int, err error) {
	needed := j.len + len(p)
	if needed > len(j.data) {
		// 指数扩容策略
		newCap := len(j.data) * 2
		if newCap < needed {
			newCap = needed
		}
		j.data = append(j.data[:j.len], make([]byte, newCap-j.len)...)
	}
	copy(j.data[j.len:], p)
	j.len += len(p)
	return len(p), nil
}

func (j *Json) Read(p []byte) (n int, err error) {
	return j.reader.Read(p)
}

func (j *Json) Len() int {
	return j.len
}

// 整个结构体清空
func (j *Json) Reset() {
	if j.len >= j.maxSize {
		j.data = make([]byte, j.defaultSize)
	}
	j.len = 0
	clear(j.data)
	j.ResetReader()
}

// 只重置偏移量，不重置数据
func (j *Json) ResetOffSet() {
	j.reader.Seek(0, io.SeekStart)
}

func (j *Json) ResetReader() {
	j.reader.Reset(j.data)
}

// WriteString 写入一个字符串到 Json.data
func (j *Json) WriteString(s string) (n int, err error) {
	// 将字符串转为 []byte 并调用已有的 Write 方法
	n, err = j.Write([]byte(s))
	j.ResetReader()
	return n, err
}

func (j *Json) WriteBytes(b []byte) (n int, err error) {
	n, err = j.Write(b)
	j.ResetReader()
	return n, err
}

/*----------------------- alert_global IStructCache interface ------------------------*/
func InitJsonSize(defaultSize, maxSize int) {
	DefaultSize = defaultSize
	MaxSize = maxSize
}

type JsonEncoder[K, V any] struct{}

func NewJsonEncoder[K, V any]() mq.Encoder[K, V] {
	return &JsonEncoder[K, V]{}
}

func (e *JsonEncoder[K, V]) EncodeKey(key *K) ([]byte, error) {
	return sonic.Marshal(key)
}

func (e *JsonEncoder[K, V]) EncodeValue(value *V) ([]byte, error) {
	return sonic.Marshal(value)
}

type JsonDecoder[K, V any] struct{}

func NewJsonDecoder[K, V any]() mq.Decoder[K, V] {
	return &JsonDecoder[K, V]{}
}

func (d *JsonDecoder[K, V]) DecodeKey(raw []byte, key *K) error {
	return nil
}

func (d *JsonDecoder[K, V]) DecodeValue(raw []byte, value *V) error {
	decoder := AcquireJson()
	defer ReleaseJson(decoder)
	_, _ = decoder.WriteBytes(raw)
	return decoder.Unmarshal(value)
}

type JsonDecoderWithKey[K, V any] struct{}

func NewJsonDecoderWithKey[K, V any]() mq.Decoder[K, V] {
	return &JsonDecoderWithKey[K, V]{}
}

func (d *JsonDecoderWithKey[K, V]) DecodeKey(raw []byte, key *K) error {
	if len(raw) > 0 && key != nil {
		decoder := AcquireJson()
		defer ReleaseJson(decoder)
		_, _ = decoder.WriteBytes(raw)
		return decoder.Unmarshal(key)
	}
	return nil
}

func (d *JsonDecoderWithKey[K, V]) DecodeValue(raw []byte, value *V) error {
	decoder := AcquireJson()
	defer ReleaseJson(decoder)
	_, _ = decoder.WriteBytes(raw)
	return decoder.Unmarshal(value)
}
