// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package config

import "fmt"

// ConfigurationChangeCallback is called for updates of configuration data
type ConfigurationChangeCallback func(key string, provider string, configdata interface{})

// A ConfigurationProvider provides a unified interface to accessing
// configuration systems.
type ConfigurationProvider interface {
	// the Name of the provider (YAML, Env, etc)
	Name() string
	// GetValue pulls a config value
	GetValue(key string) ConfigurationValue
	Scope(prefix string) ConfigurationProvider

	// A RegisterChangeCallback provides callback registration for config providers.
	// These callbacks are noop if a dynamic provider is not configured for the service.
	RegisterChangeCallback(key string, callback ConfigurationChangeCallback) string
	UnregisterChangeCallback(token string) bool
}

func keyNotFound(key string) error {
	return fmt.Errorf("couldn't find key %q", key)
}

// ScopedProvider defines recursive interface of providers based on the prefix
type ScopedProvider struct {
	ConfigurationProvider

	prefix string
}

// NewScopedProvider creates a child provider given a prefix
func NewScopedProvider(prefix string, provider ConfigurationProvider) ConfigurationProvider {
	return &ScopedProvider{provider, prefix}
}

// GetValue returns configuration value
func (sp ScopedProvider) GetValue(key string) ConfigurationValue {
	if sp.prefix != "" {
		key = fmt.Sprintf("%s.%s", sp.prefix, key)
	}
	return sp.ConfigurationProvider.GetValue(key)
}

// Scope returns new scoped provider, given a prefix
func (sp ScopedProvider) Scope(prefix string) ConfigurationProvider {
	return NewScopedProvider(prefix, sp)
}