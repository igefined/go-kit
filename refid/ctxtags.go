package refid

import "context"

type (
	ctxTags       map[string]any
	ctxTagsMarker struct{}
)

func PutTag(ctx context.Context, key string, value any) context.Context {
	return PutTags(ctx, ctxTags{key: value})
}

func PutTags(ctx context.Context, tags ctxTags) context.Context {
	for k, v := range tags {
		m, ok := markers[k]
		if ok && !m.IsValid(v) {
			tags[k] = m.Generate()
		}
	}

	cv := ctx.Value(ctxTagsMarker{})
	existingTags, ok := cv.(ctxTags)
	if ok {
		for k, v := range tags {
			existingTags[k] = v
		}
		return ctx
	}

	return context.WithValue(ctx, ctxTagsMarker{}, tags)
}

func GetTag(ctx context.Context, key string) (any, bool) {
	cv := ctx.Value(ctxTagsMarker{})
	if cv == nil {
		return nil, false
	}

	tags, ok := cv.(ctxTags)
	if !ok {
		return nil, false
	}

	v, ok := tags[key]

	return v, ok
}

func GetStringField(ctx context.Context, key string) (string, bool) {
	v, ok := GetTag(ctx, key)
	if !ok {
		return "", false
	}

	if vv, found := v.(string); found {
		return vv, true
	}

	return "", false
}

func GetFields(ctx context.Context) map[string]any {
	tags, ok := ctx.Value(ctxTagsMarker{}).(ctxTags)
	if !ok {
		return nil
	}

	return tags
}
