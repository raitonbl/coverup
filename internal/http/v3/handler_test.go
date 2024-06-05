package v3

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages/go/v21"
	"github.com/raitonbl/coverup/internal/context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestHttpContext_WithRequest(t *testing.T) {
	c := &HttpContext{
		ctx: &V3Context{
			references: make(map[string]context.Component),
			aliases:    make(map[string]map[string]context.Component),
		},
	}
	err := c.WithRequest()
	require.Nil(t, err)
	require.Equal(t, 0, len(c.ctx.(*V3Context).aliases))
	require.Equal(t, 1, len(c.ctx.(*V3Context).references))
}

func TestHttpContext_WithRequestWhenAlias(t *testing.T) {
	c := &HttpContext{
		ctx: &V3Context{
			references: make(map[string]context.Component),
			aliases:    make(map[string]map[string]context.Component),
		},
	}
	err := c.WithRequestWhenAlias("SendSmsRequest")
	require.Nil(t, err)
	require.Equal(t, 1, len(c.ctx.(*V3Context).aliases))
	require.Equal(t, 1, len(c.ctx.(*V3Context).references))
}

func TestHttpContext_WithHeaders(t *testing.T) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"Cookie":       "developer=RaitonBL",
	}
	doOnRequest(t, func(h *HttpContext) {
		rows := make([]*messages.PickleTableRow, 0)
		for k, v := range headers {
			rows = append(rows, &messages.PickleTableRow{
				Cells: []*messages.PickleTableCell{
					{
						Value: k,
					},
					{
						Value: v,
					},
				},
			})
		}
		err := h.WithHeaders(&godog.Table{
			Rows: rows,
		})
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.NotNil(t, req.headers)
		require.Len(t, req.headers, len(headers))
		for k, v := range headers {
			require.Equal(t, v, req.headers[k])
		}
	})
}

func TestHttpContext_WithHeader(m *testing.T) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"Cookie":       "developer=RaitonBL",
	}
	for name, value := range headers {
		m.Run(name, func(t *testing.T) {
			doOnRequest(t, func(h *HttpContext) {
				err := h.WithHeader(name, value)
				if err != nil {
					t.Fatal(err)
				}
			}, func(h *HttpContext) {
				req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
				require.NotNil(t, req.headers)
				require.Len(t, req.headers, 1)
				require.Equal(t, value, req.headers[name])
			})
		})
	}

}

func TestHttpContext_WithMethod(m *testing.T) {
	methods := []string{"OPTIONS", "HEAD", "GET", "PUT", "POST", "PATCH", "DELETE"}
	for _, method := range methods {
		m.Run(method, func(t *testing.T) {
			doOnRequest(t, func(h *HttpContext) {
				err := h.WithMethod(method)
				if err != nil {
					t.Fatal(err)
				}
			}, func(h *HttpContext) {
				req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
				require.Equal(t, method, req.method)
			})
		})
	}
}

func TestHttpContext_WithPath(m *testing.T) {
	arr := []string{"", "home", "vouchers", "me", "userinfo"}
	for _, p := range arr {
		m.Run(p, func(t *testing.T) {
			path := "/" + p
			doOnRequest(t, func(h *HttpContext) {
				err := h.WithPath(path)
				if err != nil {
					t.Fatal(err)
				}
			}, func(h *HttpContext) {
				req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
				require.Equal(t, path, req.path)
			})
		})
	}
}

func TestHttpContext_WithHttpPath(t *testing.T) {
	value := "localhost:8080"
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithHttpPath(value)
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Empty(t, req.path)
		require.Equal(t, "http://"+value, req.serverURL)
	})
}

func TestHttpContext_WithHttpsPath(t *testing.T) {
	value := "localhost:8080"
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithHttpsPath(value)
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Empty(t, req.path)
		require.Equal(t, "https://"+value, req.serverURL)
	})
}

func TestHttpContext_WithServerURL(t *testing.T) {
	value := "http://localhost:8080"
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithServerURL(value)
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Equal(t, value, req.serverURL)
	})
}

func TestHttpContext_WithBody(t *testing.T) {
	value := []byte(`{ "id":"83c4bdd9-58d5-4c32-998f-dd869334de73" }`)
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithBody(&godog.DocString{
			Content: string(value),
		})
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Nil(t, req.form)
		require.Equal(t, string(value), string(req.body))
	})
}

func TestHttpContext_WithBodyFileURI(t *testing.T) {
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithBodyFileURI("testdata/features/design-api/picture.base64")
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Nil(t, req.form)
		require.NotNil(t, req.body)
		require.True(t, len(req.body) > 0)
	})
}

func TestHttpContext_WithAcceptHeader(t *testing.T) {
	contentType := "application/json"
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithAcceptHeader(contentType)
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.NotNil(t, req.headers)
		require.Len(t, req.headers, 1)
		require.Equal(t, contentType, req.headers["Accept"])
	})
}

func TestHttpContext_WithContentTypeHeader(t *testing.T) {
	contentType := "application/json"
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithContentTypeHeader(contentType)
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.NotNil(t, req.headers)
		require.Len(t, req.headers, 1)
		require.Equal(t, contentType, req.headers["Content-Type"])
	})
}

func TestHttpContext_WithFormEncType(m *testing.T) {
	arr := []string{"multipart/form-data", "application/x-www-form-urlencoded"}
	for _, value := range arr {
		m.Run(value, func(t *testing.T) {
			doOnRequest(t, func(h *HttpContext) {
				err := h.WithFormEncType(value)
				if err != nil {
					t.Fatal(err)
				}
			}, func(h *HttpContext) {
				req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
				require.Nil(t, req.body)
				require.NotNil(t, req.form)
				require.Equal(t, value, req.form.encType)
			})
		})
	}
}

func TestHttpContext_WithFormAttribute(t *testing.T) {
	form := map[string]string{
		"born_at":   "2000",
		"full_name": "RaitonBL",
	}
	doOnRequest(t, func(h *HttpContext) {
		for k, v := range form {
			err := h.WithFormAttribute(k, v)
			if err != nil {
				t.Fatal(err)
			}
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Nil(t, req.body)
		require.NotNil(t, req.form)
		require.Len(t, req.form.attributes, len(form))
		for k, v := range form {
			require.Equal(t, v, req.form.attributes[k])
		}
	})
}

func TestHttpContext_WithFormFile(t *testing.T) {
	doOnRequest(t, func(h *HttpContext) {
		err := h.WithFormFile("file", "testdata/features/design-api/picture.base64")
		if err != nil {
			t.Fatal(err)
		}
	}, func(h *HttpContext) {
		req := h.ctx.(*V3Context).references[ComponentType].(*HttpRequest)
		require.Nil(t, req.body)
		require.NotNil(t, req.form)
		require.Len(t, req.form.attributes, 1)
		require.True(t, len(req.form.attributes["file"].([]byte)) > 0)
	})
}

func doOnRequest(t *testing.T, action func(*HttpContext), doAssert func(*HttpContext)) {
	home, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	c := &HttpContext{
		ctx: &V3Context{
			workDirectory: home,
			references:    make(map[string]context.Component),
			aliases:       make(map[string]map[string]context.Component),
		},
	}
	err = c.WithRequest()
	if err != nil {
		t.Fatal(err)
	}
	action(c)
	doAssert(c)
}
