package tools

import (
	"context"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// IDMixin 提供字符串类型的 id 字段
type IDMixin struct {
	mixin.Schema
	defaultFunc func() string
}

// NewIDMixin 创建 IDMixin 实例
func NewIDMixin(defaultFunc func() string) *IDMixin {
	return &IDMixin{defaultFunc: defaultFunc}
}

// Fields 返回 id 字段
func (i *IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(64).
			Unique().
			Comment("主键ID").
			DefaultFunc(i.defaultFunc).
			NotEmpty().
			Immutable(),
	}
}

// TimeMixin 提供 created_at, updated_at, deleted_at 字段
type TimeMixin struct {
	mixin.Schema
}

// Fields 返回时间戳字段
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("删除时间"),
	}
}

// Indexes 返回时间戳字段的索引
func (TimeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("updated_at"),
		index.Fields("deleted_at"),
		index.Fields("id"),
	}
}

// Hooks 返回自动更新时间戳的 hooks
func (TimeMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		// 创建时设置 created_at 和 updated_at
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpCreate) {
					now := time.Now()
					_ = setTimeField(m, "created_at", now)
					_ = setTimeField(m, "updated_at", now)
				}
				return next.Mutate(ctx, m)
			})
		},
		// 更新时设置 updated_at
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpUpdate) || m.Op().Is(ent.OpUpdateOne) {
					_ = setTimeField(m, "updated_at", time.Now())
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}

// setTimeField 安全地设置时间字段，如果字段不存在则忽略
func setTimeField(m ent.Mutation, fieldName string, value time.Time) error {
	if _, exists := m.Field(fieldName); !exists {
		return nil
	}
	if err := m.SetField(fieldName, value); err != nil {
		// 忽略字段不存在的错误（可能是字段名错误或其他原因）
		if strings.Contains(err.Error(), "field") && strings.Contains(err.Error(), "not found") {
			return nil
		}
		return err
	}
	return nil
}
