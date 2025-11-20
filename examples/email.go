package examples

import (
	"context"
	"fmt"
	ok "gook"
)

func Email(testString string) {
	ctx := context.Background()
	emailRule := ok.All(
		ok.Not(ok.StringEmpty()),
		ok.StringLength(5, 100),
		ok.StringContains("@"),
	)
	res, valid := emailRule.Validate(ctx, testString)
	fmt.Printf("valid: %v, result: %+v", valid, res)
	fmt.Println(res.Format())
}
