package internal

import (
	"fmt"
	"github.com/thoas/go-funk"
	"regexp"
	"strconv"
	"strings"
)

type ArithmeticComparison func(valueExtractor func() (any, error), compareTo float64, otherwiseThrow func(v1 any, v2 float64) error) error

func isLesserThan(valueExtractor func() (any, error), compareTo float64, otherwiseThrow func(v1 any, v2 float64) error) error {
	return compareValue(valueExtractor, compareTo, func(a, b float64) bool { return a < b }, "is lesser than", otherwiseThrow)
}

func isGreaterThan(valueExtractor func() (any, error), compareTo float64, otherwiseThrow func(v1 any, v2 float64) error) error {
	return compareValue(valueExtractor, compareTo, func(a, b float64) bool { return a > b }, "is greater than", otherwiseThrow)
}

func isLesserOrEqualTo(valueExtractor func() (any, error), compareTo float64, otherwiseThrow func(v1 any, v2 float64) error) error {
	return compareValue(valueExtractor, compareTo, func(a, b float64) bool { return a <= b }, "is lesser or equal to", otherwiseThrow)
}

func isGreaterOrEqualTo(valueExtractor func() (any, error), compareTo float64, otherwiseThrow func(v1 any, v2 float64) error) error {
	return compareValue(valueExtractor, compareTo, func(a, b float64) bool { return a >= b }, "is greater or equal to", otherwiseThrow)
}

func isEqualTo(valueExtractor func() (any, error), compareTo any, otherwiseThrow func(v1, v2 any) error) error {
	value, err := valueExtractor()
	if err != nil {
		return err
	}
	if !funk.Equal(value, compareTo) {
		return otherwiseThrow(value, compareTo)
	}
	return nil
}

func matchesPattern(valueExtractor func() (any, error), compareTo string, otherwiseThrow func(v1 string, v2 string) error) error {
	return withStringPredicate(valueExtractor, compareTo, func(valueOf, pattern string) bool {
		re := regexp.MustCompile(pattern)
		return re.MatchString(valueOf)
	}, otherwiseThrow)
}

func startsWith(valueExtractor func() (any, error), compareTo string, otherwiseThrow func(v1, v2 string) error) error {
	return withStringPredicate(valueExtractor, compareTo, func(v1, v2 string) bool {
		return strings.HasPrefix(v1, v2)
	}, otherwiseThrow)
}

func endsWith(valueExtractor func() (any, error), compareTo string, otherwiseThrow func(v1, v2 string) error) error {
	return withStringPredicate(valueExtractor, compareTo, func(v1, v2 string) bool {
		return strings.HasSuffix(v1, v2)
	}, otherwiseThrow)
}

func withStringPredicate(valueExtractor func() (any, error), compareTo string, predicate func(v1, v2 string) bool, otherwiseThrow func(v1, v2 string) error) error {
	value, err := valueExtractor()
	if err != nil {
		return err
	}
	valueStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("value is not a string")
	}
	if !predicate(valueStr, compareTo) {
		return otherwiseThrow(valueStr, compareTo)
	}
	return nil
}

func compareValue(valueExtractor func() (any, error), compareTo float64, compare func(float64, float64) bool, errorMsg string, otherwiseThrow func(v1 any, v2 float64) error) error {
	value, err := valueExtractor()
	if err != nil {
		return err
	}
	valueOf, err := toFloat64(value)
	if err != nil {
		return fmt.Errorf("cannot determine '%v' %s %v", valueOf, errorMsg, compareTo)
	}
	if !compare(valueOf, compareTo) {
		return otherwiseThrow(value, compareTo)
	}
	return nil
}

func toFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		if fv, err := strconv.ParseFloat(v, 64); err == nil {
			return fv, nil
		} else {
			return 0, fmt.Errorf("cannot parse string '%s' to float64: %w", v, err)
		}
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}
