package examples

import (
	"context"
	"errors"
	"fmt"
	ok "gook"
)

// RunBasicExamples demonstrates basic validation features
func RunBasicExamples() {
	ctx := context.Background()

	fmt.Println("=== Unified Rule Tree Validation Framework Examples ===")

	// Example 1: Basic String Validation
	fmt.Println("1. Basic String Validation")
	fmt.Println("-------------------------")
	
	stringRule := ok.All(
		ok.Not(ok.StringEmpty()),
		ok.StringLength(5, 100),
		ok.StringContains("@"),
	)

	fmt.Println("Testing valid email: test@example.com")
	result, valid := stringRule.Validate(ctx, "test@example.com")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	fmt.Println("\nTesting invalid email: ''")
	result, valid = stringRule.Validate(ctx, "")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	// Example 2: Short-Circuit Behavior
	fmt.Println("\n\n2. Short-Circuit Behavior")
	fmt.Println("-------------------------")
	
	failRule1 := ok.Test("fail1", func(ctx context.Context, n int) error {
		return errors.New("first failure")
	})
	failRule2 := ok.Test("fail2", func(ctx context.Context, n int) error {
		return errors.New("second failure")
	})
	passRule := ok.Test("pass", func(ctx context.Context, n int) error {
		return nil
	})

	allRule := ok.All(failRule1, failRule2, passRule)
	fmt.Println("Testing All rule (should stop at first failure):")
	result, valid = allRule.Validate(ctx, 42)
	fmt.Printf("Result: %v\n", valid)
	fmt.Println(result.Format())

	anyRule := ok.Any(failRule1, passRule, failRule2)
	fmt.Println("\nTesting Any rule (should stop at first success):")
	result, valid = anyRule.Validate(ctx, 42)
	fmt.Printf("Result: %v\n", valid)
	fmt.Println(result.Format())

	// Example 3: Numeric Validation
	fmt.Println("\n\n3. Numeric Validation")
	fmt.Println("---------------------")
	
	intRule := ok.NumericRange(10, 100)
	fmt.Println("Testing numeric range (10-100) with value 50:")
	result, valid = intRule.Validate(ctx, 50)
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	fmt.Println("\nTesting numeric range (10-100) with value 5:")
	result, valid = intRule.Validate(ctx, 5)
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	// Example 4: Required/Optional
	fmt.Println("\n\n4. Required/Optional")
	fmt.Println("--------------------")
	
	requiredRule := ok.Required[string]("name")
	fmt.Println("Testing Required with empty string:")
	result, valid = requiredRule.Validate(ctx, "")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	fmt.Println("\nTesting Required with non-empty string:")
	result, valid = requiredRule.Validate(ctx, "John")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	optionalRule := ok.Optional(ok.StringEmpty())
	fmt.Println("\nTesting Optional with empty string:")
	result, valid = optionalRule.Validate(ctx, "")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	fmt.Println("\nTesting Optional with valid email:")
	result, valid = optionalRule.Validate(ctx, "test@example.com")
	fmt.Printf("Result: %v\n", valid)
	if !valid {
		fmt.Println(result.Format())
	}

	fmt.Println("\n=== Examples Complete ===")
}
