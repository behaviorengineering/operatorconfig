package operatorconfig

import (
	"errors"
	"sync"

	"github.com/zalando/go-keyring"
)

// ErrNotFound means the secret is not stored in the keyring.
var ErrNotFound = errors.New("operatorconfig: secret not found in keyring")

// Keyring stores runtime secrets outside config files.
type Keyring interface {
	Get(service, account string) (string, error)
	Set(service, account, secret string) error
	Delete(service, account string) error
}

// OSKeyring uses the platform credential store via zalando/go-keyring.
type OSKeyring struct{}

func (OSKeyring) Get(service, account string) (string, error) {
	v, err := keyring.Get(service, account)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return SanitizeSecret(v), nil
}

func (OSKeyring) Set(service, account, secret string) error {
	return keyring.Set(service, account, SanitizeSecret(secret))
}

func (OSKeyring) Delete(service, account string) error {
	return keyring.Delete(service, account)
}

var (
	keyringMu      sync.RWMutex
	keyringBackend Keyring = OSKeyring{}
)

// DefaultKeyring returns the process-wide keyring backend (OS store in production).
func DefaultKeyring() Keyring {
	keyringMu.RLock()
	defer keyringMu.RUnlock()
	return keyringBackend
}

// SetTestKeyring swaps the default backend for tests. Reset in t.Cleanup when possible.
func SetTestKeyring(kr Keyring) {
	keyringMu.Lock()
	keyringBackend = kr
	keyringMu.Unlock()
}

// MemKeyring is an in-memory Keyring for unit tests.
type MemKeyring struct {
	data map[string]string
}

// NewMemKeyring returns an empty in-memory keyring.
func NewMemKeyring() *MemKeyring {
	return &MemKeyring{data: make(map[string]string)}
}

func (m *MemKeyring) key(service, account string) string {
	return service + "\x00" + account
}

func (m *MemKeyring) Get(service, account string) (string, error) {
	if m == nil || m.data == nil {
		return "", ErrNotFound
	}
	v, ok := m.data[m.key(service, account)]
	if !ok {
		return "", ErrNotFound
	}
	v = SanitizeSecret(v)
	if v == "" {
		return "", ErrNotFound
	}
	return v, nil
}

func (m *MemKeyring) Set(service, account, secret string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[m.key(service, account)] = SanitizeSecret(secret)
	return nil
}

func (m *MemKeyring) Delete(service, account string) error {
	if m.data == nil {
		return nil
	}
	delete(m.data, m.key(service, account))
	return nil
}
