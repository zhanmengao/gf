package tracetyp

import (
	"context"

	"go.opentelemetry.io/otel/baggage"
)

// Baggage 通过所有跟踪跨度保存数据
type Baggage struct {
	ctx context.Context
}

func NewBaggage(ctx context.Context) *Baggage {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Baggage{
		ctx: ctx,
	}
}

func (b *Baggage) Ctx() context.Context {
	return b.ctx
}

func (b *Baggage) SetValue(key string, value interface{}) context.Context {
	if s, ok := value.(string); ok {
		member, _ := baggage.NewMember(key, s)
		bag, _ := baggage.New(member)
		b.ctx = baggage.ContextWithBaggage(b.ctx, bag)
	}

	return b.ctx
}

func (b *Baggage) SetMap(data map[string]interface{}) context.Context {
	members := make([]baggage.Member, 0)
	for k, v := range data {
		if s, ok := v.(string); ok {
			member, _ := baggage.NewMember(k, s)
			members = append(members, member)
		}
	}
	bag, _ := baggage.New(members...)
	b.ctx = baggage.ContextWithBaggage(b.ctx, bag)
	return b.ctx
}

func (b *Baggage) GetMap() map[string]interface{} {
	m := make(map[string]interface{})
	members := baggage.FromContext(b.ctx).Members()
	for i := range members {
		m[members[i].Key()] = members[i].Value()
	}
	return m
}

func (b *Baggage) GetVar(key string) interface{} {
	value := baggage.FromContext(b.ctx).Member(key).Value()
	return value
}
