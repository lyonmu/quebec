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
	"github.com/lyonmu/quebec/pkg/common"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreUser holds the schema definition for the CoreUser entity.
type CoreUser struct {
	ent.Schema
}

// Fields of the CoreUser.
func (CoreUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Unique().Optional().Comment("用户名"),
		field.String("password").Optional().Comment("密码"),
		field.String("email").Optional().Comment("邮箱"),
		field.String("nickname").Optional().Comment("昵称"),
	}
}

// Edges of the CoreUser.
func (CoreUser) Edges() []ent.Edge {
	return nil
}

func (CoreUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(global.Id),
		tools.TimeMixin{},
	}
}

func (CoreUser) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("username"),
		index.Fields("email"),
		index.Fields("nickname"),
	}
}

func (CoreUser) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("用户信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_user", common.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
