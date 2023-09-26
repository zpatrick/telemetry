package tag

import "context"

type tagContextKey int

// ContextWithTags returns a new context with the given tags.
// If the context already has tags, the new tags will be appended.
func ContextWithTags(ctx context.Context, tags ...Tag) context.Context {
	return context.WithValue(ctx, tagContextKey(0), append(TagsFromContext(ctx), tags...))
}

// TagsFromContext returns the tags from the given context.
// If the context does not have any tags, nil is returned.
func TagsFromContext(ctx context.Context) []Tag {
	if tags, ok := ctx.Value(tagContextKey(0)).([]Tag); ok {
		return tags
	}

	return nil
}
