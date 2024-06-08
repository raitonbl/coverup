package v3

import "testing"

func TestHttpContext_when_person(t *testing.T) {
	doApply(t, "features/design-api/person.feature", nil)
}

func TestHttpContext_when_product(t *testing.T) {
	doApply(t, "features/design-api/product.feature", nil)
}
