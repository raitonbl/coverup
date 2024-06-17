package checks

import (
	"regexp"
)

var IsoDateRegexp, _ = regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
var IsoTimeRegexp, _ = regexp.Compile(`^([01]\d|2[0-3]):([0-5]\d)(:[0-5]\d)?(\.\d{3})?$`)
var IsoDateTimeRegexp, _ = regexp.Compile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])T([01]\d|2[0-3]):([0-5]\d):([0-5]\d)(\.\d{3})?(Z|([+-](0[0-9]|1[0-3]):[0-5]\d))?$`)

func IsBoolean(value any) bool {
	_, canBeCasted := value.(bool)
	return canBeCasted
}

func IsFloat64(value any) bool {
	_, canBeCasted := value.(float64)
	return canBeCasted
}

func IsString(value any) bool {
	_, canBeCasted := value.(string)
	return canBeCasted
}

func IsTime(value any) bool {
	return MatchesRegexp(value, IsoTimeRegexp)
}

func IsDate(value any) bool {
	return MatchesRegexp(value, IsoDateRegexp)
}

func ISDateTime(value any) bool {
	return MatchesRegexp(value, IsoDateTimeRegexp)
}

func MatchesRegexp(v any, regex *regexp.Regexp) bool {
	value, isString := v.(string)
	if !isString {
		return false
	}
	return regex.Match([]byte(value))
}

func IsLesserThan(v1 float64, v2 float64) bool {
	return v1 < v2
}

func IsLesserOrEqualTo(v1 float64, v2 float64) bool {
	return v1 <= v2
}

func IsGreaterThan(v1 float64, v2 float64) bool {
	return v1 > v2
}

func IsGreaterOrEqualTo(v1 float64, v2 float64) bool {
	return v1 >= v2
}
