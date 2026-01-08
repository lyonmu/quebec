package schema

import (
	"fmt"
	"net/http"

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

// CoreGatewayHttpRoute holds the schema definition for the HTTP Route entity.
type CoreGatewayHttpRoute struct {
	ent.Schema
}

// Fields of the CoreGatewayHttpRoute.
func (CoreGatewayHttpRoute) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("路由名称"),
		field.String("description").Optional().Comment("路由描述"),
		field.Int8("match_type").GoType(constant.ProxyHttpRouteMatchType(1)).Optional().Comment("匹配类型: 1-前缀 2-精确 3-正则").Default(int8(constant.HttpRouteMatchTypePrefix)),
		field.String("match_pattern").Optional().Comment("匹配规则，如 /api/v1/*"),
		field.Int("timeout_ms").Optional().Comment("路由超时(毫秒，默认15000=15秒)，包括所有重试").Default(constant.DefaultRouteTimeoutMs),
		field.Int8("enable_path_rewrite").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否启用路径重写 [1-启用 2-禁用]").Default(int8(constant.No)),
		field.String("path_rewrite").Optional().Comment("路径重写规则，如 /api/v1/* /v1/*"),
		field.Int8("enable_redirect").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否启用重定向 [1-启用 2-禁用]").Default(int8(constant.No)),
		field.String("redirect_url").Optional().Comment("重定向URL"),
		field.Int("redirect_code").Optional().Comment("重定向状态码").Default(http.StatusMovedPermanently),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("状态  [1-启用 2-禁用]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreGatewayHttpRoute.
func (CoreGatewayHttpRoute) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreGatewayHttpRoute) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreGatewayHttpRoute) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("match_type"),
		index.Fields("timeout_ms"),
		index.Fields("enable_path_rewrite"),
		index.Fields("enable_redirect"),
		index.Fields("status"),
	}
}

func (CoreGatewayHttpRoute) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("L7 HTTP路由规则表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_gateway_http_route", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
