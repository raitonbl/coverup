package internal

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractJSONPathValue(t *testing.T) {
	age := float64(15)
	name := "RaitonBL"
	id := "734f37fe-fd3f-4c4c-a046-43d547a0cdf5"
	binary := []byte(fmt.Sprintf(`{
		"id":"%s",
        "name":"%s",
		"age":%v
	}`, id, name, age))
	value, err := extractJSONPathValue(binary, "id")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, id, value)
	value, err = extractJSONPathValue(binary, "name")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, name, value)
	value, err = extractJSONPathValue(binary, "age")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, age, value)
}
