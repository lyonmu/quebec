package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreGatewayL4Listener holds the schema definition for the TCP/UDP Listener entity.
type CoreGatewayL4Listener struct {
	ent.Schema
}

// Fields of the CoreGatewayL4Listener.
func (CoreGatewayL4Listener) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("监听器名称"),
		field.String("description").Optional().Comment("监听器描述"),
		field.Uint16("port").Optional().Comment("监听端口"),
		field.String("host").Optional().Comment("监听地址").Default("0.0.0.0"),
		field.Int8("protocol").GoType(constant.ProxyProtocolType(1)).Optional().Comment("协议类型: 1-TCP 2-UDP").Default(int8(constant.ProtocolTypeTCP)),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否可用 [1: 启用, 2: 禁用]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreGatewayL4Listener.
func (CoreGatewayL4Listener) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreGatewayL4Listener) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreGatewayL4Listener) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("host"),
		index.Fields("port"),
		index.Fields("protocol"),
		index.Fields("status"),
	}
}

func (CoreGatewayL4Listener) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("L4 TCP/UDP监听器表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_gateway_l4_listener", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
