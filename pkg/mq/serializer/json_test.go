package serializer

import (
	"strings"
	"testing"
)

// --- 测试数据结构 ---

type TestStruct struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type NestedStruct struct {
	Outer string     `json:"outer"`
	Inner TestStruct `json:"inner"`
	Items []string   `json:"items"`
}

// =============================================================================
// 1. 底层引擎测试 (Low-Level Json Struct)
// =============================================================================

func TestJson_Marshal_BasicTypes(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{"string", "hello world"},
		{"int", 42},
		{"float", 3.14159},
		{"bool_true", true},
		{"slice", []int{1, 2, 3}},
		{"map", map[string]int{"a": 1, "b": 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			j := AcquireJson()
			defer ReleaseJson(j)

			// Marshal
			if err := j.Marshal(tc.value); err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// Unmarshal
			// 注意：Sonic 默认行为，解码需要匹配类型
			switch v := tc.value.(type) {
			case string:
				var res string
				if err := j.Unmarshal(&res); err != nil {
					t.Fatal(err)
				}
				if res != v {
					t.Errorf("expected %v, got %v", v, res)
				}
			case int:
				var res int
				if err := j.Unmarshal(&res); err != nil {
					t.Fatal(err)
				}
				if res != v {
					t.Errorf("expected %v, got %v", v, res)
				}
			}
		})
	}
}

func TestJson_Marshal_Struct_Newline(t *testing.T) {
	j := AcquireJson()
	defer ReleaseJson(j)

	input := TestStruct{ID: 123, Name: "test", Enabled: true}
	if err := j.Marshal(input); err != nil {
		t.Fatal(err)
	}

	// 验证底层 Sonic Stream Encoder 默认添加了换行符
	got := string(j.Bytes())
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("expected newline at end of raw stream marshal, got: %q", got)
	}

	// 验证内容正确性
	expected := `{"id":123,"name":"test","enabled":true}`
	gotTrimmed := strings.TrimSpace(got)
	if gotTrimmed != expected {
		t.Errorf("expected %s, got %s", expected, gotTrimmed)
	}
}

func TestJson_Pool_Reuse_And_Reset(t *testing.T) {
	// 模拟高并发下的复用
	obj1 := AcquireJson()
	obj1.Marshal("payload_1")
	
	// 验证 obj1 数据
	if !strings.Contains(string(obj1.Bytes()), "payload_1") {
		t.Error("obj1 data mismatch")
	}

	// 归还 obj1
	ReleaseJson(obj1)

	// 重新获取（极大可能获取到同一个对象）
	obj2 := AcquireJson()
	// 验证 Reset 是否生效：Len 应为 0
	if obj2.Len() != 0 {
		t.Errorf("expected len 0 after release/acquire, got %d", obj2.Len())
	}

	// 写入新数据
	obj2.Marshal("payload_2")
	if !strings.Contains(string(obj2.Bytes()), "payload_2") {
		t.Error("obj2 data mismatch")
	}
	// 确保没有残留旧数据
	if strings.Contains(string(obj2.Bytes()), "payload_1") {
		t.Error("obj2 contains dirty data from previous usage")
	}

	ReleaseJson(obj2)
}

func TestJson_ZeroCopy_Write_Read(t *testing.T) {
	j := AcquireJson()
	defer ReleaseJson(j)

	// 测试 WriteString (零拷贝转换)
	largeStr := strings.Repeat("a", 1000)
	n, err := j.WriteString(largeStr)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1000 {
		t.Errorf("expected 1000 bytes written, got %d", n)
	}

	// 测试 ResetReader 后的 Read
	buf := make([]byte, 500)
	readN, err := j.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if readN != 500 {
		t.Errorf("expected 500 bytes read, got %d", readN)
	}
}

// =============================================================================
// 2. 泛型 Codec 测试 (High-Level Codec)
// =============================================================================

func TestJsonCodec_RoundTrip(t *testing.T) {
	// 实例化泛型 Codec: Key=string, Value=TestStruct
	codec := NewJsonCodec[string, TestStruct]()

	key := "user:1001"
	val := TestStruct{ID: 1, Name: "Admin", Enabled: true}

	// 1. Encode Key
	// 注意：必须传入指针
	keyBytes, err := codec.MarshalKey(&key)
	if err != nil {
		t.Fatal(err)
	}
	// Codec 层应该去掉了换行符
	if string(keyBytes) != `"user:1001"` {
		t.Errorf("EncoderKey output incorrect: %q", string(keyBytes))
	}

	// 2. Encode Value
	valBytes, err := codec.MarshalValue(&val)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(valBytes), "\n") {
		t.Error("MarshalValue should strip trailing newline")
	}

	// 3. Decode Key
	var decodedKey string
	if err := codec.UnmarshalKey(keyBytes, &decodedKey); err != nil {
		t.Fatal(err)
	}
	if decodedKey != key {
		t.Errorf("Key mismatch: %v != %v", decodedKey, key)
	}

	// 4. Decode Value
	var decodedVal TestStruct
	if err := codec.UnmarshalValue(valBytes, &decodedVal); err != nil {
		t.Fatal(err)
	}
	if decodedVal != val {
		t.Errorf("Value mismatch: %+v", decodedVal)
	}
}

func TestJsonCodec_MemorySafety_DeepCopy(t *testing.T) {
	codec := NewJsonCodec[string, TestStruct]()
	val := TestStruct{ID: 999, Name: "Critical Data"}

	// Encode
	bytes1, err := codec.MarshalValue(&val)
	if err != nil {
		t.Fatal(err)
	}

	// 模拟：修改返回的切片内容
	// 如果是引用 Pool 内存，这里会破坏后续操作或被后续操作覆盖
	bytes1[0] = '{' // 已经是 '{'，这里只是占位
	
	// 强制多次进行 Pool 操作，触发内存复用
	for i := 0; i < 10; i++ {
		temp := TestStruct{ID: i}
		_, _ = codec.MarshalValue(&temp)
	}

	// 验证 bytes1 依然有效且数据未被 Pool 污染
	// 只有实现了 Deep Copy (bytes.Clone)，这里才安全
	var checkVal TestStruct
	// 注意：DecoderValue 内部使用了 sonic 直接解码，不依赖 Pool 写操作
	if err := sonicAPI.Unmarshal(bytes1, &checkVal); err != nil {
		t.Fatalf("Data corruption detected! The returned slice was likely a pool reference. Error: %v", err)
	}
	if checkVal.Name != "Critical Data" {
		t.Error("Data content changed, deep copy failed")
	}
}

func TestJsonCodec_Decoder_ZeroCopy_Path(t *testing.T) {
	// 验证 DecoderValue 是否能处理非 Pool 生成的普通 []byte
	codec := NewJsonCodec[string, map[string]int]()

	rawJSON := []byte(`{"a": 100, "b": 200}`)
	
	var result map[string]int
	err := codec.UnmarshalValue(rawJSON, &result)
	if err != nil {
		t.Fatal(err)
	}
	
	if result["a"] != 100 {
		t.Error("Decoder failed on raw bytes")
	}
}

// =============================================================================
// 3. 基准测试 (Benchmarks)
// =============================================================================

// 基准测试：底层 Json 结构体的复用性能
func BenchmarkJson_Marshal_Struct_Pooled(b *testing.B) {
	data := TestStruct{
		ID:      12345,
		Name:    "benchmark_test_name_long_string_to_fill_buffer",
		Enabled: true,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := AcquireJson()
		_ = j.Marshal(data)
		ReleaseJson(j)
	}
}

// 基准测试：泛型 Codec 的完整 Encode 流程 (含 Deep Copy)
func BenchmarkCodec_Encode_Value(b *testing.B) {
	codec := NewJsonCodec[string, TestStruct]()
	data := TestStruct{ID: 1, Name: "CodecBench", Enabled: true}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 这里包含了: Acquire -> Marshal -> DeepCopy -> Release
		_, _ = codec.MarshalValue(&data)
	}
}

// 基准测试：泛型 Codec 的 Decode 流程 (Zero Copy)
func BenchmarkCodec_Decode_Value(b *testing.B) {
	codec := NewJsonCodec[string, TestStruct]()
	data := []byte(`{"id":1,"name":"CodecBench","enabled":true}`)
	var v TestStruct

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 这里应该是 0 allocs
		_ = codec.UnmarshalValue(data, &v)
	}
}