package cast

import (
	"strings"
	"unicode"
	"unsafe"
)

func Kilobyte(v int) int {
	return v * 1024
}

func Megabyte(v int) int64 {
	return int64(Kilobyte(v)) * 1024
}

func Gigabyte(v int) int64 {
	return Megabyte(v) * 1024
}

func BytesToString(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Convert camel case string to all lower snake case.
//
// Example: ThisIsTheText => this_is_the_text
func StringLowerSnakeCase(s string) string {
	builder := strings.Builder{}
	builder.Grow(len(s) * 2)

	runes := []rune(s)
	for i, r := range runes {
		lowerR := unicode.ToLower(r)

		if r != lowerR && i != 0 {
			if !unicode.IsUpper(runes[i-1]) {
				builder.WriteRune('_')
			} else if i != len(runes) - 1 && !unicode.IsUpper(runes[i+1]) {
				builder.WriteRune('_')
			}
		}

		builder.WriteRune(lowerR)
	}

	return builder.String()
}

// NOTE: Not exactly correct for every English plural word.
func StringPlural(s string) string {
	builder := strings.Builder{}
	builder.Grow(len(s) + 3)

	runes := []rune(s)

	if runes[len(runes)-2] == 's' &&  runes[len(runes)-1] == 'h' { // sh => shes
		builder.WriteString(s)
		builder.WriteString("es")
	} else if runes[len(runes)-2] == 'c' &&  runes[len(runes)-1] == 'h' { // ch => ches
		builder.WriteString(s)
		builder.WriteString("es")
	} else if runes[len(runes)-2] == 'i' &&  runes[len(runes)-1] == 's' { // is => es
		builder.WriteString(s[:len(runes)-2])
		builder.WriteString("es")
	} else {
		switch runes[len(runes)-1] {
		case 's':
			builder.WriteString(s)
			builder.WriteString("es")
		case 'y':
			builder.WriteString(s[:len(s)-1])
			builder.WriteString("ies")
		default:
			builder.WriteString(s)
			builder.WriteString("s")		
		}
	}

	return builder.String()
}

// NOTE: might be a waste.
func CopyToSliceAny[T any](slice []T) []any {
	result := make([]any, 0, len(slice))
	for _, v := range slice {
		result = append(result, v)
	}
	return result
}
