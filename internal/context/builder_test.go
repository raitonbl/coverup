package context

import (
	"errors"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type DefaultComponent struct {
	Attributes map[string]any
	Component  *DefaultComponent
}

func (d *DefaultComponent) GetPathValue(x string) (any, error) {
	if strings.HasPrefix(x, "Attributes.") {
		return d.Attributes[x[11:]], nil
	} else if strings.HasPrefix(x, "Component.Attributes") {
		return d.Component.GetPathValue(x[10:])
	} else {
		return nil, errors.New("unexpected")
	}
}

func TestBuilder_ResolveOrGetValue(t *testing.T) {
	r := map[string]any{
		"name": "RaitonBL",
		"id":   "619a4863-f902-49cf-b47d-1dd4baf86ae8",
	}
	s := map[string]any{
		"name": "coverup",
		"id":   "c4ad02db-2d1b-4f27-be13-1c13dc0bf9ab",
	}
	c := &DefaultComponent{
		Attributes: r,
		Component: &DefaultComponent{
			Attributes: s,
			Component:  nil,
		},
	}
	b := Builder{
		context: nil,
		references: map[string]Component{
			"DefaultComponent": c,
		},
		aliases: map[string]map[string]Component{
			"DefaultComponent": {
				"R": c,
			},
		},
	}
	valueOf, err := b.GetValue("{{DefaultComponent.R.Attributes.name}}")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, c.Attributes["name"], valueOf)
	valueOf, err = b.GetValue("{{DefaultComponent.R.Component.Attributes.id}}")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, c.Component.Attributes["id"], valueOf)
	valueOf, err = b.GetValue("{{ DefaultComponent.R.Component.Attributes.id }}")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, c.Component.Attributes["id"], valueOf)

	valueOf, err = b.GetValue("{{ DefaultComponent.R.Component.Attributes.id }} > {{ DefaultComponent.R.Component.Attributes.name }}")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, c.Component.Attributes["id"].(string)+" > "+c.Component.Attributes["name"].(string), valueOf)
}
