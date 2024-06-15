package http

import (
	"fmt"
)

type FactoryOpts[S any] struct {
	Settings                    S
	AssertTrue                  bool
	AssertAlias                 bool
	ResolveValueBeforeAssertion bool
}

func createResponseBodyJsonPathRegexp(expr string) []string {
	return createResponseBodyRegexp(fmt.Sprintf(`\$(\S+) %s`, expr))
}

func createResponseBodyRegexp(expr string) []string {
	return []string{
		fmt.Sprintf(`^(?i)response body %s$`, expr),
		fmt.Sprintf(`^(?i)the response body %s$`, expr),
		fmt.Sprintf(`^(?i)the Response body %s$`, expr),
	}
}

func createAliasedRequestResponseBodyJsonPathRegexp(expr string) []string {
	return createAliasedRequestResponseBodyRegexp(fmt.Sprintf(`\$(\S+) %s`, expr))
}

func createAliasedRequestResponseBodyRegexp(expr string) []string {
	return []string{
		fmt.Sprintf(`%s response body %s$`, httpRequestRegex, expr),
	}
}

func createRequestLinePart(expr string) []string {
	return []string{
		fmt.Sprintf(`^%s`, expr),
		fmt.Sprintf(`^(?i)the %s`, expr),
		//fmt.Sprintf(`^(?i)%s the `, strings.ToUpper(string(expr[0]))+expr[1:]),
	}
}
