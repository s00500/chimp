package components

// type ctxKey string

// const (
// 	ctxKeySlot     ctxKey = "slot"
// 	ctxKeyTemplate ctxKey = "template"
// )

// func Root() templ.Component {
// 	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
// 		ctx = addSlotData(ctx)
// 		return templ.GetChildren(ctx).Render(ctx, w)
// 	})
// }

// var slotIndex int64

// type None struct{}

// func addSlotData(ctx context.Context) context.Context {
// 	sd := make(map[int64][]interface{})
// 	td := make(map[int64][]interface{})
// 	ctx = context.WithValue(ctx, ctxKeySlot, &sd)
// 	ctx = context.WithValue(ctx, ctxKeyTemplate, &td)
// 	return ctx
// }

// func mustGetSlotData(ctx context.Context) *map[int64][]interface{} {
// 	sd, ok := ctx.Value(ctxKeySlot).(*map[int64][]interface{})
// 	if !ok {
// 		panic("slot data not found")
// 	}
// 	return sd
// }

// func mustGetTemplateData(ctx context.Context) *map[int64][]interface{} {
// 	td, ok := ctx.Value(ctxKeyTemplate).(*map[int64][]interface{})
// 	if !ok {

// 		panic("template data not found")
// 	}
// 	return td
// }

// type Slot[T any] struct {
// 	id int64
// }

// func DefineSlot[T any]() *Slot[T] {
// 	return &Slot[T]{
// 		id: atomic.AddInt64(&slotIndex, 1),
// 	}
// }

// func (s *Slot[T]) addValue(ctx context.Context, v ...T) {
// 	sd := mustGetSlotData(ctx)
// 	if _, ok := (*sd)[s.id]; !ok {
// 		(*sd)[s.id] = make([]interface{}, 0)
// 	}
// 	if len(v) > 0 {
// 		(*sd)[s.id] = append((*sd)[s.id], v[0])
// 	} else {
// 		(*sd)[s.id] = append((*sd)[s.id], nil)
// 	}
// }

// func (s *Slot[T]) hasValue(ctx context.Context) bool {
// 	sd := mustGetSlotData(ctx)
// 	return len((*sd)[s.id]) > 0
// }

// func (s *Slot[T]) removeValue(ctx context.Context) {
// 	sd := mustGetSlotData(ctx)
// 	if !s.hasValue(ctx) {
// 		panic("slot not found")
// 	}
// 	(*sd)[s.id] = (*sd)[s.id][:len((*sd)[s.id])-1]
// }

// func (ts *Slot[T]) addTemplate(ctx context.Context, t templ.Component) {
// 	td := mustGetTemplateData(ctx)
// 	if _, ok := (*td)[ts.id]; !ok {
// 		(*td)[ts.id] = make([]interface{}, 0)
// 	}
// 	(*td)[ts.id] = append((*td)[ts.id], t)
// }

// func (ts *Slot[T]) hasTemplate(ctx context.Context) bool {
// 	td := mustGetTemplateData(ctx)
// 	return len((*td)[ts.id]) > 0
// }

// func (ts *Slot[T]) removeTemplate(ctx context.Context) {
// 	td := mustGetTemplateData(ctx)
// 	if !ts.hasTemplate(ctx) {
// 		panic("template not found")
// 	}
// 	(*td)[ts.id] = (*td)[ts.id][:len((*td)[ts.id])-1]
// }

// func (ts *Slot[T]) getTemplate(ctx context.Context) templ.Component {
// 	td := mustGetTemplateData(ctx)
// 	return (*td)[ts.id][len((*td)[ts.id])-1].(templ.Component)
// }

// func (s *Slot[T]) Value(ctx context.Context) T {
// 	sd := mustGetSlotData(ctx)
// 	if !s.hasValue(ctx) {
// 		panic("slot not found")
// 	}
// 	return (*sd)[s.id][len((*sd)[s.id])-1].(T)
// }

// func (s *Slot[T]) WithValue(v T) templ.Component {
// 	if _, ok := any(v).(None); ok {
// 		return s
// 	}
// 	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
// 		s.addValue(ctx, v)
// 		c := new(strings.Builder)
// 		templ.GetChildren(ctx).Render(ctx, c)
// 		defer s.removeValue(ctx)
// 		if s.hasTemplate(ctx) {
// 			defer s.removeTemplate(ctx)
// 			return s.getTemplate(ctx).Render(ctx, w)
// 		}
// 		_, err := io.Copy(w, strings.NewReader(c.String()))
// 		return err
// 	})
// }

// func (s *Slot[T]) Render(ctx context.Context, w io.Writer) error {
// 	s.addValue(ctx)
// 	c := new(strings.Builder)
// 	templ.GetChildren(ctx).Render(ctx, c)
// 	defer s.removeValue(ctx)
// 	if s.hasTemplate(ctx) {
// 		defer s.removeTemplate(ctx)
// 		return s.getTemplate(ctx).Render(ctx, w)
// 	}
// 	_, err := io.Copy(w, strings.NewReader(c.String()))
// 	return err
// }

// func (s *Slot[T]) Template() templ.Component {
// 	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
// 		if s.hasValue(ctx) {
// 			s.addTemplate(ctx, templ.GetChildren(ctx))
// 		}
// 		return nil
// 	})
// }
