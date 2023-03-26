package excuse

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestNewJohn(t *testing.T) {
	john, err := NewJohn()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	env := Env{prng: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var sb strings.Builder
	for i := 0; i < 10; i++ {
		sb.Reset()
		if err = john.Expand(ctx, &sb, &env); err != nil {
			t.Fatal(err)
		}
		t.Log(">", sb.String())
	}
}
