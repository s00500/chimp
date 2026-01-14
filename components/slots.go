package components

import (
	"context"
	"io"
	"strings"
	"sync/atomic"

	"github.com/a-h/templ"
)

type ctxKey string

const (
	ctxKeySlot     ctxKey = "slot"
	ctxKeyTemplate ctxKey = "template"
)

var slotIndex int64

type None struct{}

func addSlotData(ctx context.Context) context.Context {
	sd := make(map[int64][]interface{})
	td := make(map[int64][]interface{})
	ctx = context.WithValue(ctx, ctxKeySlot, &sd)
	ctx = context.WithValue(ctx, ctxKeyTemplate, &td)
	return ctx
}

func ensureSlotData(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKeySlot).(*map[int64][]interface{}); ok {
		return ctx
	}
	return addSlotData(ctx)
}

func getSlotData(ctx context.Context) *map[int64][]interface{} {
	sd, ok := ctx.Value(ctxKeySlot).(*map[int64][]interface{})
	if !ok {
		return nil
	}
	return sd
}

func getTemplateData(ctx context.Context) *map[int64][]interface{} {
	td, ok := ctx.Value(ctxKeyTemplate).(*map[int64][]interface{})
	if !ok {
		return nil
	}
	return td
}

type Slot[T any] struct {
	id int64
}

func DefineSlot[T any]() *Slot[T] {
	return &Slot[T]{
		id: atomic.AddInt64(&slotIndex, 1),
	}
}

func (s *Slot[T]) addValue(ctx context.Context, v ...T) {
	sd := getSlotData(ctx)
	if sd == nil {
		return
	}
	if _, ok := (*sd)[s.id]; !ok {
		(*sd)[s.id] = make([]interface{}, 0)
	}
	if len(v) > 0 {
		(*sd)[s.id] = append((*sd)[s.id], v[0])
	} else {
		(*sd)[s.id] = append((*sd)[s.id], nil)
	}
}

func (s *Slot[T]) hasValue(ctx context.Context) bool {
	sd := getSlotData(ctx)
	if sd == nil {
		return false
	}
	return len((*sd)[s.id]) > 0
}

func (s *Slot[T]) removeValue(ctx context.Context) {
	sd := getSlotData(ctx)
	if sd == nil || !s.hasValue(ctx) {
		return
	}
	(*sd)[s.id] = (*sd)[s.id][:len((*sd)[s.id])-1]
}

func (ts *Slot[T]) addTemplate(ctx context.Context, t templ.Component) {
	td := getTemplateData(ctx)
	if td == nil {
		return
	}
	if _, ok := (*td)[ts.id]; !ok {
		(*td)[ts.id] = make([]interface{}, 0)
	}
	(*td)[ts.id] = append((*td)[ts.id], t)
}

func (ts *Slot[T]) hasTemplate(ctx context.Context) bool {
	td := getTemplateData(ctx)
	if td == nil {
		return false
	}
	return len((*td)[ts.id]) > 0
}

func (ts *Slot[T]) removeTemplate(ctx context.Context) {
	td := getTemplateData(ctx)
	if td == nil || !ts.hasTemplate(ctx) {
		return
	}
	(*td)[ts.id] = (*td)[ts.id][:len((*td)[ts.id])-1]
}

func (ts *Slot[T]) getTemplate(ctx context.Context) templ.Component {
	td := getTemplateData(ctx)
	if td == nil {
		return nil
	}
	return (*td)[ts.id][len((*td)[ts.id])-1].(templ.Component)
}

func (s *Slot[T]) Value(ctx context.Context) T {
	sd := getSlotData(ctx)
	if sd == nil || !s.hasValue(ctx) {
		var zero T
		return zero
	}
	return (*sd)[s.id][len((*sd)[s.id])-1].(T)
}

func (s *Slot[T]) WithValue(v T) templ.Component {
	if _, ok := any(v).(None); ok {
		return s
	}
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = ensureSlotData(ctx)
		s.addValue(ctx, v)
		c := new(strings.Builder)
		templ.GetChildren(ctx).Render(ctx, c)
		defer s.removeValue(ctx)
		if s.hasTemplate(ctx) {
			defer s.removeTemplate(ctx)
			return s.getTemplate(ctx).Render(ctx, w)
		}
		_, err := io.Copy(w, strings.NewReader(c.String()))
		return err
	})
}

func (s *Slot[T]) Render(ctx context.Context, w io.Writer) error {
	ctx = ensureSlotData(ctx)
	s.addValue(ctx)
	c := new(strings.Builder)
	templ.GetChildren(ctx).Render(ctx, c)
	defer s.removeValue(ctx)
	if s.hasTemplate(ctx) {
		defer s.removeTemplate(ctx)
		return s.getTemplate(ctx).Render(ctx, w)
	}
	_, err := io.Copy(w, strings.NewReader(c.String()))
	return err
}

func (s *Slot[T]) Template() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = ensureSlotData(ctx)
		if s.hasValue(ctx) {
			s.addTemplate(ctx, templ.GetChildren(ctx))
		}
		return nil
	})
}
