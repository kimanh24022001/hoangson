package user

import (
	"errors"
	"strings"
)

var (
	Error_PhoneInvalidLen    = errors.New("The number of characters in the phone number is invalid")
	Error_PhoneInvalidChar   = errors.New("Invalid characters exist in the phone number")
	Error_EmailInvalid       = errors.New("Invalid email")
	Error_EmailInvalidPrefix = errors.New("Invalid email prefix")
	Error_EmailInvalidDomain = errors.New("Invalid email domain")
)

// TODO(duong): make sure it is correct.
func ValidatePhone(phone string) error {
	phoneLen := len(phone)

	if phoneLen != 9 &&
		phoneLen != 10 &&
		phoneLen != 11 {
		return Error_PhoneInvalidLen
	}

	ok := true
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			continue
		}

		ok = false
		break
	}

	if !ok {
		return Error_PhoneInvalidLen
	}

	return nil
}

func ValidateEmail(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return Error_EmailInvalid
	}

	prefix := strings.ToLower(parts[0])
	domain := parts[1]

	// NOTE:
	// Acceptable email prefix formats:
	// - Allowed characters: letters (a-z), numbers, underscores, periods,
	// and dashes.
	// - An underscore, period, or dash must be followed by one or more
	// letter or number.
	{
		ok := true
		for _, r := range prefix {
			if (r >= 'a' && r <= 'z') ||
				(r >= '0' && r <= '9') {
				ok = true
				continue
			}

			if r == '_' || r == '.' || r == '-' {
				ok = false
				continue
			}

			ok = false
			break
		}

		if !ok {
			return Error_EmailInvalidPrefix
		}
	}

	// NOTE: Check the domain
	// Acceptable email domain formats
	// - Allowed characters: letters, numbers, dashes.
	// - The last portion of the domain must be at least two characters, for
	// example: .com, .org, .cc
	{
		{
			count := 0
			runes := []rune(domain)

			for i := len(runes) - 1; i >= 0; i-- {
				if runes[i] == '.' {
					if i == 0 {
						return Error_EmailInvalidDomain
					}
					break
				}
				count++
			}

			if count < 2 {

				return Error_EmailInvalidDomain
			}
		}

		var ok = true
		for _, r := range domain {
			if (r >= 'a' && r <= 'z') ||
				(r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') ||
				(r == '.') ||
				(r == '-') {
				continue
			}

			ok = false
			break
		}

		if !ok {
			return Error_EmailInvalidDomain
		}
	}

	return nil
}
