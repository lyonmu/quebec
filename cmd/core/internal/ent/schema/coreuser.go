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
		field.String("password").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).Optional().Comment("密码"),
		field.String("email").Optional().Comment("邮箱"),
		field.String("nickname").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).Optional().Comment("昵称"),
		field.Int8("status").Optional().Comment("用户状态 [1: 启用, 2: 禁用]").Default(1),
		field.String("role_id").Optional().Comment("角色ID"),
		field.String("remark").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).Optional().Comment("用户备注"),
	}
}

// Edges of the CoreUser.
// 用户与角色多对一关系
func (CoreUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user_from_role", CoreRole.Type).Ref("role_to_user").Field("role_id").Unique(),
	}
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
		index.Fields("status"),
		index.Fields("role_id"),
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
