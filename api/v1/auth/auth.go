package auth

import (
	"context"

	"github.com/nairobi-gophers/fupisha/internal/config"
	"github.com/nairobi-gophers/fupisha/internal/store"
)

//Resource defines dependencies for auth handlers.
type Resource struct {
	Store  store.Store
	Config *config.Config
}

//NewResource returns a configured authentication resource.
func NewResource(store store.Store, cfg *config.Config) *Resource {
	return &Resource{
		Store:  store,
		Config: cfg,
	}
}

//FromContext extracts user id from a Context
func FromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
