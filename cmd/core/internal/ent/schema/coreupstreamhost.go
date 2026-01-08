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

// CoreUpstreamHost holds the schema definition for the Upstream Host entity.
type CoreUpstreamHost struct {
	ent.Schema
}

// Fields of the CoreUpstreamHost.
func (CoreUpstreamHost) Fields() []ent.Field {
	return []ent.Field{
		field.String("address").Optional().Comment("后端地址IP"),
		field.Int("weight").Optional().Comment("权重(相对权重)").Default(1),
		field.Int("port").Optional().Comment("后端端口"),
		field.Int8("enabled").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否可用 [1: 是, 2: 否]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreUpstreamHost.
func (CoreUpstreamHost) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreUpstreamHost) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreUpstreamHost) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("address"),
		index.Fields("port"),
		index.Fields("weight"),
		index.Fields("enabled"),
	}
}

func (CoreUpstreamHost) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("上游服务后端地址表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_upstream_host", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
