package api

import (
	"sync"

	"bitbucket.org/digitdreamteam/mirana/core/mdl"
)

type AccountManager struct {
	Users map[string]*mdl.User
	l     *sync.Mutex
}

var _instance *AccountManager

var _once sync.Once

func Manager() *AccountManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &AccountManager{
		Users: make(map[string]*mdl.User),
		l:     &sync.Mutex{},
	}
}

func (m *AccountManager) GetUser(id string) *mdl.User {
	m.l.Lock()
	defer m.l.Unlock()
	if u, ok := m.Users[id]; ok {
		return u
	}
	return nil
}

func (m *AccountManager) AddAmount(id string, amount float32) {
	m.l.Lock()
	defer m.l.Unlock()
	u := m.getUser(id)
	if u != nil {
		u.Balance += amount
	}
}

func (m *AccountManager) CheckAmountGreaderThan(id string, amount float32) bool {
	m.l.Lock()
	defer m.l.Unlock()
	u := m.getUser(id)
	return u != nil && u.Balance >= amount
}

func (m *AccountManager) getUser(id string) *mdl.User {
	if u, ok := m.Users[id]; ok {
		return u
	}
	return nil
}
