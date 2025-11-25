package response

type Options struct {
	Label    string     `json:"label"`              // 名称
	Value    string     `json:"value"`              // 值
	Children []*Options `json:"children,omitempty"` // 子级选项
}
