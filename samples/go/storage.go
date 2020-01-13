package main

import (
	"math/rand"
	"strconv"

	"github.com/coreos/go-oidc"
)

type stateStorage map[string]struct{}

func (storage stateStorage) new() string {
	state := randomString()
	storage[state] = struct{}{}
	return state
}

func (storage stateStorage) isValid(state string) bool {
	_, ok := storage[state]
	return ok
}

func (storage stateStorage) pop(state string) {
	delete(storage, state)
}

type userInfoStorage map[string]*oidc.UserInfo

func (storage userInfoStorage) add(userInfo *oidc.UserInfo) string {
	id := randomString()
	storage[id] = userInfo
	return id
}

func (storage userInfoStorage) get(id string) *oidc.UserInfo {
	return storage[id]
}

func randomString() string {
	// Do not do this in production!
	return strconv.Itoa(rand.Intn(999999))
}
