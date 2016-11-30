package mux

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type contextKey int

const (
	queriesKey contextKey = iota
	routeKey
)

// getQueries returns the route variables for the current request, if any.
func getQueries(r *http.Request) queries {
	if rv := contextGet(r, queriesKey); rv != nil {
		return rv.(queries)
	}
	return nil
}

// CurrentRoute returns the matched route for the current request, if any.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns, unless the KeepContext option is set on the
// Router.
func CurrentRoute(r *http.Request) *Route {
	if rv := contextGet(r, routeKey); rv != nil {
		return rv.(*Route)
	}
	return nil
}

func setQueries(r *http.Request) *http.Request {
	queries, _ := extractQueries(r)
	return contextSet(r, queriesKey, queries)
}

func setCurrentRoute(r *http.Request, val interface{}) *http.Request {
	return contextSet(r, routeKey, val)
}

func contextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func contextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}

type queries map[string][]string

// Get return the key value, of the current *http.Request quer
func (q queries) Get(key string) []string {
	if value, found := q[key]; found {
		return value
	}
	return []string{}
}

func extractQueries(req *http.Request) (queries, error) {

	queriesRaw, err := url.ParseQuery(req.URL.RawQuery)

	if err != nil {
		return nil, err
	}

	queries := queries(map[string][]string{})
	for k, v := range queriesRaw {
		for _, item := range v {
			values := strings.Split(item, ",")
			queries[k] = append(queries[k], values...)
		}
	}

	return queries, nil
}
