package secctx

import (
	"github.com/shaj13/go-guardian/v2/auth"
	"net/http"
	"strings"
)

type SecurityContext interface {
	GetUserId() string
	GetUserToken() string
	GetApiKey() string
	IsSystem() bool
}

func Create(r *http.Request) SecurityContext {
	token := getAuthorizationToken(r)
	if token != "" {
		user := auth.User(r)
		userId := user.GetID()
		return &securityContextImpl{
			userId:   userId,
			token:    token,
			apiKey:   "",
			isSystem: false,
		}
	} else {
		apiKey := getApihubApiKey(r)
		userId := "api-key_" + apiKey
		return &securityContextImpl{
			userId:   userId,
			token:    "",
			apiKey:   apiKey,
			isSystem: false,
		}
	}
}

func CreateSystemContext() SecurityContext {
	return &securityContextImpl{isSystem: true}
}

type securityContextImpl struct {
	userId   string
	token    string
	apiKey   string
	isSystem bool
}

func getAuthorizationToken(r *http.Request) string {
	authorizationHeaderValue := r.Header.Get("authorization")
	return strings.ReplaceAll(authorizationHeaderValue, "Bearer ", "")
}

func getApihubApiKey(r *http.Request) string {
	return r.Header.Get("api-key")
}

func (ctx securityContextImpl) GetUserId() string {
	return ctx.userId
}
func (ctx securityContextImpl) GetUserToken() string {
	return ctx.token
}
func (ctx securityContextImpl) GetApiKey() string {
	return ctx.apiKey
}
func (ctx securityContextImpl) IsSystem() bool { return ctx.isSystem }
