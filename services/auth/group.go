// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package auth

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
)

// Ensure the struct implements the interface.
var (
	_ Method        = &Group{}
	_ Initializable = &Group{}
	_ Freeable      = &Group{}
)

// Group implements the Auth interface with serval Auth.
type Group struct {
	methods []Method
}

// NewGroup creates a new auth group
func NewGroup(methods ...Method) *Group {
	return &Group{
		methods: methods,
	}
}

// Add adds a new method to group
func (b *Group) Add(method Method) {
	b.methods = append(b.methods, method)
}

// Name returns group's methods name
func (b *Group) Name() string {
	names := make([]string, 0, len(b.methods))
	for _, m := range b.methods {
		if n, ok := m.(Named); ok {
			names = append(names, n.Name())
		} else {
			names = append(names, reflect.TypeOf(m).Elem().Name())
		}
	}
	return strings.Join(names, ",")
}

// Init does nothing as the Basic implementation does not need to allocate any resources
func (b *Group) Init(ctx context.Context) error {
	for _, method := range b.methods {
		initializable, ok := method.(Initializable)
		if !ok {
			continue
		}

		if err := initializable.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Free does nothing as the Basic implementation does not have to release any resources
func (b *Group) Free() error {
	for _, method := range b.methods {
		freeable, ok := method.(Freeable)
		if !ok {
			continue
		}
		if err := freeable.Free(); err != nil {
			return err
		}
	}
	return nil
}

// Verify extracts and validates
func (b *Group) Verify(req *http.Request, w http.ResponseWriter, store DataStore, sess SessionStore) *user_model.User {
	if !db.HasEngine {
		return nil
	}

	// Try to sign in with each of the enabled plugins
	for _, ssoMethod := range b.methods {
		user := ssoMethod.Verify(req, w, store, sess)
		if user != nil {
			if store.GetData()["AuthedMethod"] == nil {
				if named, ok := ssoMethod.(Named); ok {
					store.GetData()["AuthedMethod"] = named.Name()
				}
			}
			return user
		}
	}

	return nil
}
