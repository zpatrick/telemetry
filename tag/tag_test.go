package tag_test

import (
	"fmt"

	"github.com/zpatrick/telemetry/tag"
)

func ExampleWrite() {
	tags := []tag.Tag{
		tag.New("host", "localhost"),
		tag.New("port", 8080),
	}

	tag.Write(tags, func(key string, val interface{}) {
		fmt.Println(key, val)
	})

	// Output:
	// host localhost
	// port 8080
}

func ExampleGroup() {
	g := tag.Group("server", tag.New("host", "localhost"), tag.New("port", 8080))

	tag.Write([]tag.Tag{g}, func(key string, val interface{}) {
		fmt.Println(key, val)
	})

	// Output:
	// server.host localhost
	// server.port 8080
}
