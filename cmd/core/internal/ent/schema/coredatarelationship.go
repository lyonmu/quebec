package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreDataRelationship holds the schema definition for the CoreDataRelationship entity.
type CoreDataRelationship struct {
	ent.Schema
}

// Fields of the CoreDataRelationship.
func (CoreDataRelationship) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("data_relationship_type").GoType(corecommon.DataRelationshipType(1)).Optional().Comment("数据关系类型"),
		field.String("menu_id").Optional().Comment("菜单ID"),
		field.String("role_id").Optional().Comment("角色ID"),
	}
}

// Edges of the CoreDataRelationship.
func (CoreDataRelationship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("menu", CoreMenu.Type),
		edge.To("role", CoreRole.Type),
	}
}

func (CoreDataRelationship) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreDataRelationship) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{}
}

func (CoreDataRelationship) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("数据关系信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_data_relationship", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
