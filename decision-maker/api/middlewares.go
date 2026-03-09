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
		})
	}
}

func getPublicEndpoints() endpointSet {
	apiSpec, err := GetSwagger()
	if err != nil {
		log.Fatal("error parsing OpenAPI spec:", err)
	}
	pe := *newEndpointSet()
	for path, pathItem := range apiSpec.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			if operation.Security == nil || len(*operation.Security) == 0 {
				endpoint := fmt.Sprintf("%s %s", strings.ToUpper(method), consts.APIBaseUrl+path)
				pe.add(endpoint)
			}
		}
	}
	pe.add("GET /api/docs")
	pe.add("GET /api/spec")
	return pe
}

type endpointSet struct {
	data map[string]any
}

func newEndpointSet() *endpointSet {
	return &endpointSet{data: make(map[string]any)}
}

func (s endpointSet) add(endpoint string) {
	s.data[endpoint] = nil
}

func (s endpointSet) contains(endpoint string) bool {
	parts := strings.Split(endpoint, " ")
	method, pathSegments := parts[0], strings.Split(parts[1], "/")

	for publicEndpoint := range s.data {
		publicEndpointParts := strings.Split(publicEndpoint, " ")
		publicEndpointPathTemplate := strings.Split(publicEndpointParts[1], "/")
		if publicEndpointParts[0] == method && pathMatchesTemplate(pathSegments, publicEndpointPathTemplate) {
			return true
		}
	}
	return false
}

func (s endpointSet) empty() bool {
	return len(s.data) == 0
}

func pathMatchesTemplate(actual, template []string) bool {
	if len(actual) != len(template) {
		return false
	}
	for i, seg := range template {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			continue
		}
		if seg != actual[i] {
			return false
		}
	}
	return true
}
