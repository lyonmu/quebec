package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	corecommon "github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreMenu holds the schema definition for the CoreMenu entity.
type CoreMenu struct {
	ent.Schema
}

// Fields of the CoreMenu.
func (CoreMenu) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("菜单名称"),
		field.Int8("menu_type").GoType(corecommon.MenuType(1)).Optional().Comment("菜单类型 [1: 目录, 2: 菜单, 3: 按钮]").Default(int8(corecommon.MenuTypeDirectory)),
		field.String("api_path").Optional().Comment("菜单API路径"),
		field.String("api_path_method").Optional().Comment("菜单API方法"),
		field.Int8("order").Optional().Comment("菜单排序").Default(1),
		field.String("parent_id").Optional().Comment("父菜单ID"),
		field.Int8("status").GoType(constant.YesOrNo(1)).Optional().Comment("菜单状态 [1: 启用, 2: 禁用]").Default(int8(constant.Yes)),
		field.String("component").Optional().Comment("菜单组件"),
		field.String("remark").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).Optional().Comment("菜单备注"),
	}
}

// Edges of the CoreMenu.
func (CoreMenu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu_from_parent", CoreMenu.Type).Ref("menu_to_children").Field("parent_id").Unique(),
		edge.To("menu_to_data_relationship", CoreDataRelationship.Type),
		edge.To("menu_to_children", CoreMenu.Type),
	}
}

func (CoreMenu) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreMenu) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("name"),
		index.Fields("menu_type"),
		index.Fields("api_path"),
		index.Fields("status"),
		index.Fields("parent_id"),
	}
}

func (CoreMenu) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("菜单信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_menu", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
