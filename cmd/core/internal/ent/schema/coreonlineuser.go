package schema

import (
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/lyonmu/quebec/cmd/core/internal/common"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreOnLineUser holds the schema definition for the CoreOnLineUser entity.
type CoreOnLineUser struct {
	ent.Schema
}

// Fields of the CoreOnLineUser.
func (CoreOnLineUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").Unique().Optional().Comment("用户ID").Unique(),
		field.String("access_ip").Optional().Comment("访问IP"),
		field.Int64("last_operation_time").Optional().Comment("最后操作时间").DefaultFunc(func() int64 { return time.Now().Unix() }),
		field.Int("operation_type").Optional().GoType(common.OperationType(1)).Optional().Comment("操作类型 [1: 登陆]"),
		field.String("os").Optional().Comment("操作系统"),
		field.String("platform").Optional().Comment("操作平台"),
		field.String("browser_name").Optional().Comment("浏览器名称"),
		field.String("browser_version").Optional().Comment("浏览器版本"),
		field.String("browser_engine_name").Optional().Comment("浏览器引擎名称"),
		field.String("browser_engine_version").Optional().Comment("浏览器引擎版本"),
	}

}

// Edges of the CoreOnLineUser.
func (CoreOnLineUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("on_line_from_user", CoreUser.Type).Ref("on_line_to_user").Field("user_id").Unique(),
	}
}

func (CoreOnLineUser) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreOnLineUser) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("access_ip"),
		index.Fields("last_operation_time"),
		index.Fields("operation_type"),
	}
}

func (CoreOnLineUser) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("在线用户信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_on_line_user", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
