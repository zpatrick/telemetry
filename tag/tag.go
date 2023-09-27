package tag

// A Tag is used to annotate telemetry data.
type Tag interface {
	write(f func(key string, val any))
}

// TagValue is a type that can be used as a tag value.
type TagValue interface {
	~string | ~int | ~int64 | ~float64 | ~bool
}

type tag[T TagValue] struct {
	key string
	val T
}

// New creates a new tag with the given key and value.
func New[T TagValue](key string, value T) Tag {
	return &tag[T]{key, value}
}

func (t tag[T]) write(f func(key string, val any)) {
	f(t.key, t.val)
}

type group struct {
	name string
	tags []Tag
}

// Group converts a list of tags into a single tag where each key is prefixed with the group name.
func Group(name string, tags ...Tag) Tag {
	return group{name: name, tags: tags}
}

func (g group) write(f func(key string, val any)) {
	Write(g.tags, func(key string, val any) {
		f(g.name+"."+key, val)
	})
}

// Write calls f for each key:value pair in tags.
func Write(tags []Tag, f func(key string, val any)) {
	for _, tag := range tags {
		tag.write(f)
	}
}
