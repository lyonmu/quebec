package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lyonmu/quebec/cmd/core/internal/global"
	"github.com/lyonmu/quebec/pkg/constant"
	"github.com/lyonmu/quebec/pkg/tools"
)

// CoreCert holds the schema definition for the TLS Certificate entity.
type CoreCert struct {
	ent.Schema
}

//SerialNumber string // 证书序列号
// IssuerCN     string // 颁发者CN
// IssuerOrg    string // 颁发者组织
// SubjectCN    string // 主题CN
// SubjectOrg   string // 主题组织
// NotBefore    int64  // 有效期开始时间
// NotAfter     int64  // 有效期结束时间
// Fingerprint  string // 证书指纹

// Fields of the CoreCert.
func (CoreCert) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("证书名称"),
		field.String("serial_number").Optional().Comment("证书序列号"),
		field.String("issuer_cn").Optional().Comment("颁发者CN"),
		field.String("subject_cn").Optional().Comment("主题CN"),
		field.Int64("not_before").Optional().Comment("有效期开始时间"),
		field.Int64("not_after").Optional().Comment("有效期结束时间"),
		field.Int8("secret_type").GoType(constant.CertType(1)).Optional().Comment("证书类型: 1-服务证书 2-根证书").Default(int8(constant.ServerCert)),
		field.String("certificate_hash").Optional().Comment("证书SHA256哈希"),
		field.String("certificate").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).
			Optional().Comment("证书PEM格式"),
		field.String("private_key_hash").Optional().Comment("私钥SHA256哈希"),
		field.String("private_key").SchemaType(map[string]string{dialect.MySQL: "text", dialect.SQLite: "text", dialect.Postgres: "text"}).
			Optional().Comment("私钥PEM格式"),
		field.Int8("status").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否可用 [1: 是, 2: 否]").Default(int8(constant.Yes)),
		field.Int8("is_default").Optional().GoType(constant.YesOrNo(1)).Optional().Comment("是否为默认证书 [1: 是, 2: 否]").Default(int8(constant.Yes)),
	}
}

// Edges of the CoreCert.
func (CoreCert) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CoreCert) Mixin() []ent.Mixin {
	return []ent.Mixin{
		tools.NewIDMixin(func() string {
			return fmt.Sprintf("%d", global.Id.GenID())
		}),
		tools.TimeMixin{},
	}
}

func (CoreCert) Indexes() []ent.Index {
	return []ent.Index{}
}

func (CoreCert) Annotations() []schema.Annotation {
	withCommentsEnabled := true
	return []schema.Annotation{
		schema.Comment("证书信息表"),
		entsql.Annotation{
			Table:        fmt.Sprintf("%s_core_cert", constant.ProjectName),
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_general_ci",
			WithComments: &withCommentsEnabled,
		},
		edge.Annotation{StructTag: `json:"-" gorm:"-"`},
	}
}
