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
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreGatewayNode holds the schema definition for the CoreGatewayNode entity.
type CoreGatewayNode struct {
	ent.Schema
}

// Fields of the CoreGatewayNode.
func (CoreGatewayNode) Fields() []ent.Field {
	return []ent.Field{
		field.String("node_id").Optional().Comment("节点ID").Unique(),
		field.String("cluster_id").Optional().Comment("集群ID"),
		field.Int64("node_register_time").Optional().Comment("注册时间").DefaultFunc(func() int64 { return time.Now().Unix() }),
		field.Int64("node_last_request_time").Optional().Comment("最新请求时间").DefaultFunc(func() int64 { return time.Now().Unix() }),
	}
}

// Edges of the CoreGatewayNode.
func (CoreGatewayNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("node_from_cluster", CoreGatewayCluster.Type).Ref("cluster_to_node").Field("cluster_id").Unique(),
	}
}

func (CoreGatewayNode) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreGatewayNode) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("cluster_id"),
		index.Fields("node_id"),
		index.Fields("node_register_time"),
		index.Fields("node_last_request_time"),
	}
}

func (CoreGatewayNode) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("网关节点信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_gateway_node", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
