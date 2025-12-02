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

// CoreGatewayCluster holds the schema definition for the CoreGatewayCluster entity.
type CoreGatewayCluster struct {
	ent.Schema
}

// Fields of the CoreGatewayCluster.
func (CoreGatewayCluster) Fields() []ent.Field {
	return []ent.Field{
		field.String("cluster_id").Optional().Comment("集群ID").Unique(),
		field.Int64("cluster_create_time").Optional().Comment("创建时间").DefaultFunc(func() int64 { return time.Now().Unix() }),
	}
}

// Edges of the CoreGatewayCluster.
func (CoreGatewayCluster) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("cluster_to_node", CoreGatewayNode.Type),
	}
}

func (CoreGatewayCluster) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreGatewayCluster) Indexes() []ent.Index {
	// 时间戳字段的索引由 TimeMixin 提供
	return []ent.Index{
		index.Fields("cluster_id"),
		index.Fields("cluster_create_time"),
	}
}

func (CoreGatewayCluster) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("网关集群信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_gateway_cluster", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
