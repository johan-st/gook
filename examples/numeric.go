package examples

import (
	"context"
	"errors"
	"fmt"

	ok "github.com/johan-st/gook"
)

func Numeric(testString any) {
	ctx := context.Background()
	pipeline := ok.Then(
		ok.All(ok.Test("not-nil", func(ctx context.Context, v any) error {
			if v == nil { return errors.New("required") }
			return nil
		})),
		func(v any) (int, error) {
			s, ok := v.(string)
			if !ok { return 0, errors.New("must be string") }
			var n int
			_, err := fmt.Sscanf(s, "%d", &n)
			if err != nil { return 0, errors.New("must be numeric") }
			return n, nil
		},
		ok.All(
			ok.Test("range", func(ctx context.Context, n int) error {
				if n < 10 || n > 100 {
					return fmt.Errorf("value must be between 10 and 100")
				}
				return nil
			}),
			ok.Not(ok.Test("not-13", func(ctx context.Context, n int) error {
				if n == 13 {
					return nil
				}
				return fmt.Errorf("value is not 13")
			})),
		))
	result, ok := pipeline.Validate(ctx, testString)
	fmt.Printf("valid: %v\n", ok)
	fmt.Println(result.Format())
}