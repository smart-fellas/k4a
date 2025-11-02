package context

import (
	"fmt"
)

type Manager struct {
	config *Config
}

func NewManager(cfg *Config) *Manager {
	return &Manager{config: cfg}
}

func (m *Manager) ListContexts() []string {
	var names []string
	for _, ctx := range m.config.Contexts {
		names = append(names, ctx.Name)
	}
	return names
}

func (m *Manager) SwitchContext(name string) error {
	for _, ctx := range m.config.Contexts {
		if ctx.Name == name {
			m.config.CurrentContext = name
			return m.config.Save()
		}
	}
	return fmt.Errorf("context %s not found", name)
}

func (m *Manager) GetCurrentNamespace() string {
	ctx, err := m.config.GetCurrentContext()
	if err != nil {
		return ""
	}
	return ctx.Context.Namespace
}

func (m *Manager) SetNamespace(namespace string) error {
	ctx, err := m.config.GetCurrentContext()
	if err != nil {
		return err
	}

	for i, c := range m.config.Contexts {
		if c.Name == ctx.Name {
			m.config.Contexts[i].Context.Namespace = namespace
			return m.config.Save()
		}
	}

	return fmt.Errorf("failed to update namespace")
}
