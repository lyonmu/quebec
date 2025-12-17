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

// CoreRole holds the schema definition for the CoreRole entity.
type CoreRole struct {
	ent.Schema
}

// Fields of the CoreRole.
func (CoreRole) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("角色名称"),
		field.String("remark").Optional().Comment("角色备注"),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("角色状态 [1: 启用, 2: 禁用]").Default(int8(constant.Yes)),
		field.Int8("system").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否系统角色 [1: 是, 2: 否]").Default(int8(constant.No)),
	}
}

// Edges of the CoreRole.
// 角色与用户一对多关系
func (CoreRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role_to_user", CoreUser.Type),
		edge.To("role_to_data_relationship", CoreDataRelationship.Type),
	}
}

func (CoreRole) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreRole) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("name"),
		index.Fields("status"),
	}
}

func (CoreRole) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("角色信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_role", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
