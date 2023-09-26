package tag_test

import (
	"context"
	"testing"

	"github.com/zpatrick/telemetry/tag"
	"github.com/zpatrick/testx/assert"
)

func TestTagsFromContext_none(t *testing.T) {
	tags := tag.TagsFromContext(context.Background())
	if tags != nil {
		t.Fatal("tags should be nil", tags)
	}
}

func TestContextWithTags(t *testing.T) {
	parent := tag.ContextWithTags(context.Background(), tag.New("parent", 1))
	child := tag.ContextWithTags(parent, tag.New("child", 2))

	// The parent context should only contain the parent tags.
	parentTags := map[string]any{}
	tag.Write(tag.TagsFromContext(parent), func(key string, val any) {
		parentTags[key] = val
	})

	assert.EqualMaps(t, parentTags, map[string]interface{}{"parent": 1})

	// The child context should contain both the parent and child tags.
	childTags := map[string]any{}
	tag.Write(tag.TagsFromContext(child), func(key string, val any) {
		childTags[key] = val
	})

	assert.EqualMaps(t, childTags, map[string]interface{}{"parent": 1, "child": 2})
}
