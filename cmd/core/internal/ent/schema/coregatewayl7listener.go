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

// CoreGatewayL7Listener holds the schema definition for the HTTP Listener entity.
type CoreGatewayL7Listener struct {
	ent.Schema
}

// Fields of the CoreGatewayL7Listener.
func (CoreGatewayL7Listener) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("监听器名称"),
		field.String("description").Optional().Comment("监听器描述"),
		field.Uint16("port").Optional().Comment("监听端口"),
		field.String("host").Optional().Comment("监听地址").Default("0.0.0.0"),
		field.Int8("enable_tls").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否启用TLS [1: 启用, 2: 禁用]").Default(int8(constant.No)),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否启用 [1: 启用, 2: 禁用]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreGatewayL7Listener.
func (CoreGatewayL7Listener) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreGatewayL7Listener) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreGatewayL7Listener) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("host"),
		index.Fields("port"),
		index.Fields("enable_tls"),
		index.Fields("status"),
	}
}

func (CoreGatewayL7Listener) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("L7 HTTP监听器表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_gateway_l7_listener", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
