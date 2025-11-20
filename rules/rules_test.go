package rules

import (
	"context"
	"testing"
)

func TestEmail(t *testing.T) {
	ctx := context.Background()
	rule := Email()

	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.com",
		"user_name@example.org",
	}

	invalidEmails := []string{
		"invalid",
		"@example.com",
		"user@",
		"user@domain",
		"user @example.com",
	}

	for _, email := range validEmails {
		_, ok := rule.Validate(ctx, email)
		if !ok {
			t.Errorf("Expected %q to be valid", email)
		}
	}

	for _, email := range invalidEmails {
		_, ok := rule.Validate(ctx, email)
		if ok {
			t.Errorf("Expected %q to be invalid", email)
		}
	}
}

func TestURL(t *testing.T) {
	ctx := context.Background()
	rule := URL()

	validURLs := []string{
		"http://example.com",
		"https://example.com/path",
		"ftp://files.example.com",
		"http://example.com:8080/path?query=value",
	}

	invalidURLs := []string{
		"not-a-url",
		"example.com",
		"://example.com",
		"http://",
	}

	for _, u := range validURLs {
		_, ok := rule.Validate(ctx, u)
		if !ok {
			t.Errorf("Expected %q to be valid", u)
		}
	}

	for _, u := range invalidURLs {
		_, ok := rule.Validate(ctx, u)
		if ok {
			t.Errorf("Expected %q to be invalid", u)
		}
	}
}

func TestUUID(t *testing.T) {
	ctx := context.Background()
	rule := UUID()

	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-41d1-80b4-00c04fd430c8", // v4 UUID
	}

	invalidUUIDs := []string{
		"not-a-uuid",
		"550e8400-e29b-41d4-a716",
		"550e8400-e29b-41d4-a716-446655440000-extra",
		"550e8400-e29b-31d4-a716-446655440000", // wrong version
	}

	for _, u := range validUUIDs {
		_, ok := rule.Validate(ctx, u)
		if !ok {
			t.Errorf("Expected %q to be valid", u)
		}
	}

	for _, u := range invalidUUIDs {
		_, ok := rule.Validate(ctx, u)
		if ok {
			t.Errorf("Expected %q to be invalid", u)
		}
	}
}

func TestCreditCard(t *testing.T) {
	ctx := context.Background()
	rule := CreditCard()

	// Valid test card numbers (Luhn algorithm compliant)
	validCards := []string{
		"4532015112830366", // Visa test number
		"4532-0151-1283-0366",
		"4532 0151 1283 0366",
	}

	invalidCards := []string{
		"1234567890123456", // Invalid Luhn
		"1234",
		"not-a-card",
	}

	for _, card := range validCards {
		_, ok := rule.Validate(ctx, card)
		if !ok {
			t.Errorf("Expected %q to be valid", card)
		}
	}

	for _, card := range invalidCards {
		_, ok := rule.Validate(ctx, card)
		if ok {
			t.Errorf("Expected %q to be invalid", card)
		}
	}
}

func TestIPAddress(t *testing.T) {
	ctx := context.Background()
	rule := IPAddress()

	validIPs := []string{
		"192.168.1.1",
		"8.8.8.8",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"::1",
	}

	invalidIPs := []string{
		"999.999.999.999",
		"not-an-ip",
		"256.1.1.1",
	}

	for _, ip := range validIPs {
		_, ok := rule.Validate(ctx, ip)
		if !ok {
			t.Errorf("Expected %q to be valid", ip)
		}
	}

	for _, ip := range invalidIPs {
		_, ok := rule.Validate(ctx, ip)
		if ok {
			t.Errorf("Expected %q to be invalid", ip)
		}
	}
}

func TestDomain(t *testing.T) {
	ctx := context.Background()
	rule := Domain()

	validDomains := []string{
		"example.com",
		"sub.example.com",
		"example.co.uk",
	}

	invalidDomains := []string{
		"not-a-domain",
		".example.com",
		"example.",
		"-example.com",
	}

	for _, domain := range validDomains {
		_, ok := rule.Validate(ctx, domain)
		if !ok {
			t.Errorf("Expected %q to be valid", domain)
		}
	}

	for _, domain := range invalidDomains {
		_, ok := rule.Validate(ctx, domain)
		if ok {
			t.Errorf("Expected %q to be invalid", domain)
		}
	}
}

func TestHexColor(t *testing.T) {
	ctx := context.Background()
	rule := HexColor()

	validColors := []string{
		"#FF0000",
		"#ff0000",
		"#F00",
		"#FF0000FF",
	}

	invalidColors := []string{
		"FF0000",
		"#GG0000",
		"#FF00",
		"#FF00000",
	}

	for _, color := range validColors {
		_, ok := rule.Validate(ctx, color)
		if !ok {
			t.Errorf("Expected %q to be valid", color)
		}
	}

	for _, color := range invalidColors {
		_, ok := rule.Validate(ctx, color)
		if ok {
			t.Errorf("Expected %q to be invalid", color)
		}
	}
}

func TestBase64(t *testing.T) {
	ctx := context.Background()
	rule := Base64()

	validBase64 := []string{
		"SGVsbG8gV29ybGQ=",
		"dGVzdA==",
		"",
	}

	invalidBase64 := []string{
		"SGVsbG8gV29ybGQ", // missing padding
		"SGVsbG8gV29ybGQ===", // too much padding
		"SGVsbG8gV29ybGQ!",
	}

	for _, b64 := range validBase64 {
		_, ok := rule.Validate(ctx, b64)
		if !ok {
			t.Errorf("Expected %q to be valid", b64)
		}
	}

	for _, b64 := range invalidBase64 {
		_, ok := rule.Validate(ctx, b64)
		if ok {
			t.Errorf("Expected %q to be invalid", b64)
		}
	}
}

func TestJSON(t *testing.T) {
	ctx := context.Background()
	rule := JSON()

	validJSON := []string{
		`{"key": "value"}`,
		`[1, 2, 3]`,
		`"string"`,
		`123`,
		`true`,
		`null`,
		`{}`,
		`[]`,
	}

	invalidJSON := []string{
		`{key: value}`,
		`[1, 2, 3`,
		`not json`,
		`{"key": "value"`,
	}

	for _, jsonStr := range validJSON {
		_, ok := rule.Validate(ctx, jsonStr)
		if !ok {
			t.Errorf("Expected %q to be valid", jsonStr)
		}
	}

	for _, jsonStr := range invalidJSON {
		_, ok := rule.Validate(ctx, jsonStr)
		if ok {
			t.Errorf("Expected %q to be invalid", jsonStr)
		}
	}
}

