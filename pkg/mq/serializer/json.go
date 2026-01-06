package serializer

import (
	"bytes"
	"errors"
	"sync"
	"unsafe"

	"github.com/bytedance/sonic"
)

const (
	DefaultSize = 4 * 1024
	MaxSize     = 10 * 1024 * 1024
)

var sonicAPI = sonic.ConfigFastest

var poolJson = &sync.Pool{
	New: func() any {
		return newJson()
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
	len    int
	data   []byte
	reader *bytes.Reader // 保持指针，避免拷贝结构体
}

func newJson() *Json {
	res := &Json{
		data: make([]byte, DefaultSize),
	}
	// 初始化时创建一个 reader 实例，后续只 Reset 它
	res.reader = bytes.NewReader(res.data[:0])
	return res
}

// Write 实现 io.Writer
func (j *Json) Write(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		return 0, nil
	}
	if j.len+n > cap(j.data) {
		if err := j.grow(n); err != nil {
			return 0, err // ✅ 返回错误而不是 panic
		}
	}
	copy(j.data[j.len:], p)
	j.len += n
	return n, nil
}

// Read 实现 io.Reader 接口
// 必须存在，否则无法被 sonic.NewDecoder 或其他 IO 工具调用
func (j *Json) Read(p []byte) (n int, err error) {
	return j.reader.Read(p)
}

// WriteString 零拷贝写入字符串
// 包含 ResetReader 以确保写入后立即可读（保持与原有逻辑兼容）
func (j *Json) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}
	// 优化：使用 unsafe 避免 string -> []byte 的内存分配与拷贝
	src := unsafe.Slice(unsafe.StringData(s), len(s))
	n, err = j.Write(src)
	
	// 关键：写入后重置 Reader 范围，使新数据对 Read 可见
	// 如果不加这一行，Write 后直接 Read 会读不到刚写入的数据
	j.ResetReader()
	return n, err
}

// WriteBytes 写入字节切片
// 包含 ResetReader 以确保写入后立即可读
func (j *Json) WriteBytes(b []byte) (n int, err error) {
	n, err = j.Write(b)
	j.ResetReader()
	return n, err
}

// grow 指数扩容 + 错误处理
func (j *Json) grow(n int) error {
	needed := j.len + n
	currentCap := cap(j.data)
	if needed <= currentCap {
		return nil
	}
	newCap := currentCap * 2
	if newCap == 0 {
		newCap = DefaultSize
	}
	for newCap < needed {
		newCap *= 2
	}
	if newCap > MaxSize {
		return errors.New("json buffer overflow: exceeds max size limit")
	}
	newData := make([]byte, newCap)
	copy(newData, j.data[:j.len])
	j.data = newData
	return nil
}

// Reset 逻辑
func (j *Json) Reset() {
	if cap(j.data) > MaxSize {
		j.data = make([]byte, DefaultSize)
	}
	j.len = 0
	j.ResetReader()
}

func (j *Json) ResetReader() {
	// ✅ 修正：复用对象，不产生 GC
	j.reader.Reset(j.data[:j.len])
}

// Marshal 封装
func (j *Json) Marshal(v any) error {
	j.len = 0 // 直接重置
	// Sonic Encoder 默认会加 \n，这是 Stream 协议决定的
	err := sonicAPI.NewEncoder(j).Encode(v)
	if err != nil {
		return err
	}
	j.ResetReader()
	return nil
}

// Unmarshal 封装
func (j *Json) Unmarshal(v any) error {
	j.ResetReader()
	return sonicAPI.NewDecoder(j.reader).Decode(v)
}

// Bytes 获取引用
func (j *Json) Bytes() []byte {
	return j.data[:j.len]
}

// Len 获取当前有效数据长度
// 复杂度: O(1)
func (j *Json) Len() int {
	return j.len
}


// --- MQ 接口实现 ---

// JsonCodec 泛型 JSON 编解码器
type JsonCodec[K, V any] struct{}

func NewJsonCodec[K, V any]() *JsonCodec[K, V] {
	return &JsonCodec[K, V]{}
}

// MarshalKey 实现 Serializer 接口
func (c *JsonCodec[K, V]) MarshalKey(key *K) ([]byte, error) {
	return c.encode(key)
}

// MarshalValue 实现 Serializer 接口
func (c *JsonCodec[K, V]) MarshalValue(value *V) ([]byte, error) {
	return c.encode(value)
}

func (c *JsonCodec[K, V]) encode(v any) ([]byte, error) {
	j := AcquireJson()
	defer ReleaseJson(j)

	if err := j.Marshal(v); err != nil {
		return nil, err
	}

	b := j.Bytes()
	// 处理 sonic Encoder 带来的尾部换行
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}

	// ✅ 必须拷贝，Go 1.20+ 推荐 bytes.Clone
	return bytes.Clone(b), nil
}

// UnmarshalKey 实现 Deserializer 接口
func (c *JsonCodec[K, V]) UnmarshalKey(raw []byte, key *K) error {
	// ✅ 修正：直接解码，零拷贝
	if len(raw) == 0 {
		return nil
	}
	return sonicAPI.Unmarshal(raw, key)
}

// UnmarshalValue 实现 Deserializer 接口
func (c *JsonCodec[K, V]) UnmarshalValue(raw []byte, value *V) error {
	// ✅ 修正：直接解码，零拷贝
	if len(raw) == 0 {
		return nil
	}
	return sonicAPI.Unmarshal(raw, value)
}