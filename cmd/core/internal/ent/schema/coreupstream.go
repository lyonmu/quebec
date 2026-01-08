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

// CoreUpstream holds the schema definition for the Upstream entity.
type CoreUpstream struct {
	ent.Schema
}

// Fields of the CoreUpstream.
func (CoreUpstream) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("上游服务名称"),
		field.String("description").Optional().Comment("上游服务描述"),
		field.Int8("lb_policy").GoType(constant.ProxyLbPolicy(1)).Optional().Comment("负载均衡策略: 1-ROUND_ROBIN 2-LEAST_REQUEST 3-RANDOM 4-RING_HASH 5-MAGLEV").Default(int8(constant.LbPolicyMaglev)),
		field.Int("connect_timeout_ms").Optional().Comment("连接超时(毫秒)").Default(constant.DefaultConnectTimeoutMs),
		field.Int("max_connections").Optional().Comment("最大连接数").Default(constant.DefaultMaxConnections),
		field.Int("max_pending_requests").Optional().Comment("最大等待请求数").Default(constant.DefaultMaxPendingRequests),
		field.Int("max_requests").Optional().Comment("最大请求数").Default(constant.DefaultMaxRequests),
		field.Int("max_retries").Optional().Comment("最大重试次数").Default(constant.DefaultMaxRetries),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("状态 [1-启用 2-禁用]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreUpstream.
func (CoreUpstream) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreUpstream) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreUpstream) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("lb_policy"),
		index.Fields("connect_timeout_ms"),
		index.Fields("max_connections"),
		index.Fields("max_pending_requests"),
		index.Fields("max_requests"),
		index.Fields("max_retries"),
		index.Fields("status"),
	}
}

func (CoreUpstream) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("上游服务信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_upstream", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
