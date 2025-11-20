package rules

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/johan-st/gook"
)

// Email creates a rule that validates RFC 5322 compliant email addresses
// This is a simplified version that covers most common email formats
func Email() *gook.Rule[string] {
	// RFC 5322 simplified regex - covers most practical email addresses
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return gook.Test("email", func(ctx context.Context, value string) error {
		if !emailRegex.MatchString(value) {
			return errors.New("invalid email format")
		}
		return nil
	})
}

// URL creates a rule that validates URL format
func URL() *gook.Rule[string] {
	return gook.Test("url", func(ctx context.Context, value string) error {
		u, err := url.Parse(value)
		if err != nil {
			return fmt.Errorf("invalid URL format: %v", err)
		}
		if u.Scheme == "" {
			return errors.New("URL must have a scheme (e.g., http, https)")
		}
		if u.Host == "" {
			return errors.New("URL must have a host")
		}
		return nil
	})
}

// PhoneUS creates a rule that validates US phone number formats
// Accepts formats: (123) 456-7890, 123-456-7890, 123.456.7890, 1234567890, +1 123 456 7890
func PhoneUS() *gook.Rule[string] {
	// US phone: 10 digits, optional country code +1, various separators
	phoneRegex := regexp.MustCompile(`^(\+?1[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})$`)
	return gook.Test("phone-us", func(ctx context.Context, value string) error {
		cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(value, "")
		if len(cleaned) < 10 || len(cleaned) > 11 {
			return errors.New("US phone number must have 10 or 11 digits")
		}
		if len(cleaned) == 11 && cleaned[0] != '1' {
			return errors.New("US phone number with country code must start with 1")
		}
		if !phoneRegex.MatchString(value) {
			return errors.New("invalid US phone number format")
		}
		return nil
	})
}

// PhoneInternational creates a rule that validates international phone number formats
// Accepts formats with country codes (e.g., +44 20 7946 0958, +1-555-123-4567)
func PhoneInternational() *gook.Rule[string] {
	// International: starts with +, followed by 1-15 digits with optional separators
	phoneRegex := regexp.MustCompile(`^\+\d{1,3}[-.\s]?\d{1,14}[-.\s]?\d{1,14}$`)
	return gook.Test("phone-international", func(ctx context.Context, value string) error {
		if !phoneRegex.MatchString(value) {
			return errors.New("invalid international phone number format (must start with +)")
		}
		cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(value, "")
		if len(cleaned) < 7 || len(cleaned) > 15 {
			return errors.New("international phone number must have 7-15 digits")
		}
		return nil
	})
}

// UUID creates a rule that validates UUID v4 format
func UUID() *gook.Rule[string] {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return gook.Test("uuid", func(ctx context.Context, value string) error {
		if !uuidRegex.MatchString(strings.ToLower(value)) {
			return errors.New("invalid UUID v4 format (expected: xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx)")
		}
		return nil
	})
}

// luhnCheck validates a number using the Luhn algorithm
func luhnCheck(number string) bool {
	sum := 0
	alternate := false
	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

// CreditCard creates a rule that validates credit card numbers using Luhn algorithm
// Accepts formats with or without spaces/dashes
func CreditCard() *gook.Rule[string] {
	// Credit card: 13-19 digits, may contain spaces or dashes
	cardRegex := regexp.MustCompile(`^[\d\s\-]{13,19}$`)
	return gook.Test("credit-card", func(ctx context.Context, value string) error {
		cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(value, "")
		if len(cleaned) < 13 || len(cleaned) > 19 {
			return errors.New("credit card number must have 13-19 digits")
		}
		if !cardRegex.MatchString(value) {
			return errors.New("invalid credit card format")
		}
		if !luhnCheck(cleaned) {
			return errors.New("invalid credit card number (Luhn check failed)")
		}
		return nil
	})
}

// IPAddress creates a rule that validates IPv4 or IPv6 addresses
func IPAddress() *gook.Rule[string] {
	return gook.Test("ip-address", func(ctx context.Context, value string) error {
		ip := net.ParseIP(value)
		if ip == nil {
			return errors.New("invalid IP address format (must be IPv4 or IPv6)")
		}
		return nil
	})
}

// IPv4 creates a rule that validates IPv4 addresses only
func IPv4() *gook.Rule[string] {
	return gook.Test("ipv4", func(ctx context.Context, value string) error {
		ip := net.ParseIP(value)
		if ip == nil {
			return errors.New("invalid IPv4 address format")
		}
		if ip.To4() == nil {
			return errors.New("not an IPv4 address (use IPAddress() for IPv6 support)")
		}
		return nil
	})
}

// IPv6 creates a rule that validates IPv6 addresses only
func IPv6() *gook.Rule[string] {
	return gook.Test("ipv6", func(ctx context.Context, value string) error {
		ip := net.ParseIP(value)
		if ip == nil {
			return errors.New("invalid IPv6 address format")
		}
		if ip.To4() != nil {
			return errors.New("not an IPv6 address (use IPv4() for IPv4 support)")
		}
		return nil
	})
}

// Domain creates a rule that validates domain name format
func Domain() *gook.Rule[string] {
	// Domain: alphanumeric, hyphens, dots; must start/end with alphanumeric; TLD at least 2 chars
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return gook.Test("domain", func(ctx context.Context, value string) error {
		if !domainRegex.MatchString(value) {
			return errors.New("invalid domain name format")
		}
		if len(value) > 253 {
			return errors.New("domain name too long (max 253 characters)")
		}
		return nil
	})
}

// HexColor creates a rule that validates CSS hex color codes (#RRGGBB or #RRGGBBAA)
func HexColor() *gook.Rule[string] {
	hexColorRegex := regexp.MustCompile(`^#([0-9a-fA-F]{6}|[0-9a-fA-F]{3}|[0-9a-fA-F]{8})$`)
	return gook.Test("hex-color", func(ctx context.Context, value string) error {
		if !hexColorRegex.MatchString(value) {
			return errors.New("invalid hex color format (expected #RRGGBB, #RGB, or #RRGGBBAA)")
		}
		return nil
	})
}

// Base64 creates a rule that validates Base64 encoded strings
func Base64() *gook.Rule[string] {
	return gook.Test("base64", func(ctx context.Context, value string) error {
		_, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return fmt.Errorf("invalid Base64 encoding: %v", err)
		}
		return nil
	})
}

// JSON creates a rule that validates JSON string format
func JSON() *gook.Rule[string] {
	return gook.Test("json", func(ctx context.Context, value string) error {
		var js interface{}
		if err := json.Unmarshal([]byte(value), &js); err != nil {
			return fmt.Errorf("invalid JSON format: %v", err)
		}
		return nil
	})
}

