package api

import (
	"context"
	"decisionMaker/consts"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(validateToken func(token string) (string, error)) func(http.Handler) http.Handler {
	pe := getPublicEndpoints()
	isPublic := func(e string) bool {
		return pe.contains(e)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint := r.Method + " " + r.URL.Path
			if isPublic(endpoint) { // if endpoint is public, don't authenticate
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"message":"missing or invalid authorization header"}`, http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			uid, err := validateToken(token)
			if err != nil {
				http.Error(w, `{"message":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", uid)
			next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
		})
	}
}

func getPublicEndpoints() set {
	apiSpec, err := GetSwagger()
	if err != nil {
		log.Fatal("error parsing OpenAPI spec:", err)
	}
	pe := *newSet()
	for path, pathItem := range apiSpec.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			if operation.Security == nil || len(*operation.Security) == 0 {
				key := fmt.Sprintf("%s %s", strings.ToUpper(method), consts.APIBaseUrl+path)
				pe.add(key)
			}
		}
	}
	pe.add("GET /api/docs")
	pe.add("GET /api/spec")
	return pe
}

type set struct {
	data map[string]any
}

func newSet() *set {
	return &set{data: make(map[string]any)}
}

func (s set) add(item string) {
	s.data[item] = nil
}

func (s set) contains(item string) bool {
	if _, ok := s.data[item]; ok {
		return true
	}
	return false
}

func (s set) empty() bool {
	return len(s.data) == 0
}
