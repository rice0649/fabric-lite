package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/rice0649/fabric-lite/internal/providers"
)

var (
	defaultProviderManager *ProviderManager
	pmOnce                 sync.Once
)

// SetDefaultProviderManager sets the global ProviderManager instance.
// This should only be called once during application initialization.
func SetDefaultProviderManager(pm *ProviderManager) {
	pmOnce.Do(func() {
		defaultProviderManager = pm
	})
}

// GetDefaultProviderManager returns the global ProviderManager instance.
// Panics if the manager has not been set.
func GetDefaultProviderManager() *ProviderManager {
	if defaultProviderManager == nil {
		panic("Default ProviderManager has not been initialized. Call SetDefaultProviderManager first.")
	}
	return defaultProviderManager
}

// ProviderManager manages AI provider instances and their lifecycle
type ProviderManager struct {
	providers   map[string]providers.Provider
	config      *providers.Config
	mutex       sync.RWMutex
	initialized bool
}


// NewProviderManager creates a new provider manager
func NewProviderManager(config *providers.Config) *ProviderManager {
	return &ProviderManager{
		providers: make(map[string]providers.Provider),
		config:    config,
	}
}

// InitializeAll creates and registers all configured providers
func (pm *ProviderManager) InitializeAll() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.initialized {
		return nil
	}

	for _, providerConfig := range pm.config.Providers {
		provider, err := providers.NewProvider(providerConfig)
		if err != nil {
			// Log warning but continue with other providers
			fmt.Printf("Warning: failed to create provider %s: %v\n", providerConfig.Name, err)
			continue
		}
		pm.providers[providerConfig.Name] = provider
	}

	pm.initialized = true
	return nil
}

// Get returns a specific provider by name
func (pm *ProviderManager) Get(name string) (providers.Provider, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if !pm.initialized {
		return nil, fmt.Errorf("provider manager not initialized")
	}

	provider, ok := pm.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}

	return provider, nil
}

// GetDefault returns the configured default provider
func (pm *ProviderManager) GetDefault() (providers.Provider, error) {
	if pm.config == nil {
		return nil, fmt.Errorf("no configuration available")
	}

	defaultName := pm.config.DefaultProvider
	if defaultName == "" {
		defaultName = "openai" // fallback
	}

	return pm.Get(defaultName)
}

// ListAvailable returns a list of all available providers
func (pm *ProviderManager) ListAvailable() []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	available := make([]string, 0, len(pm.providers))
	for name := range pm.providers {
		available = append(available, name)
	}
	return available
}

// ListReady returns a list of providers that are ready for use (IsAvailable() = true)
func (pm *ProviderManager) ListReady() []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	ready := make([]string, 0)
	for name, provider := range pm.providers {
		if provider.IsAvailable() {
			ready = append(ready, name)
		}
	}
	return ready
}

// CheckAvailability checks if a specific provider is ready for use
func (pm *ProviderManager) CheckAvailability(name string) bool {
	provider, err := pm.Get(name)
	if err != nil {
		return false
	}
	return provider.IsAvailable()
}

// GetModels returns available models for a specific provider
func (pm *ProviderManager) GetModels(name string) ([]string, error) {
	provider, err := pm.Get(name)
	if err != nil {
		return nil, err
	}

	return provider.GetModels(), nil
}

// AddProvider dynamically adds a new provider
func (pm *ProviderManager) AddProvider(name string, provider providers.Provider) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.providers[name] = provider
}

// RemoveProvider removes a provider from the manager
func (pm *ProviderManager) RemoveProvider(name string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	delete(pm.providers, name)
}

// ReloadProvider reloads a specific provider from configuration
func (pm *ProviderManager) ReloadProvider(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Find provider config
	var providerConfig *providers.ProviderConfig
	for _, pc := range pm.config.Providers {
		if pc.Name == name {
			providerConfig = &pc
			break
		}
	}

	if providerConfig == nil {
		return fmt.Errorf("provider configuration not found: %s", name)
	}

	// Create new provider instance
	provider, err := providers.NewProvider(*providerConfig)
	if err != nil {
		return fmt.Errorf("failed to recreate provider %s: %w", name, err)
	}

	pm.providers[name] = provider
	return nil
}

// GetConfigForProvider returns the configuration for a specific provider
func (pm *ProviderManager) GetConfigForProvider(name string) (*providers.ProviderConfig, error) {
	if pm.config == nil {
		return nil, fmt.Errorf("no configuration available")
	}

	for _, provider := range pm.config.Providers {
		if provider.Name == name {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("provider configuration not found: %s", name)
}

// Execute executes a completion request using the specified provider
func (pm *ProviderManager) Execute(ctx context.Context, providerName string, request providers.CompletionRequest) (*providers.CompletionResponse, error) {
	provider, err := pm.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", providerName, err)
	}

	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider not available: %s", providerName)
	}

	return provider.Execute(ctx, request)
}

// ExecuteStream executes a streaming completion request
func (pm *ProviderManager) ExecuteStream(ctx context.Context, providerName string, request providers.CompletionRequest) (<-chan providers.StreamChunk, error) {
	provider, err := pm.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", providerName, err)
	}

	if !provider.IsAvailable() {
		return nil, fmt.Errorf("provider not available: %s", providerName)
	}

	return provider.ExecuteStream(ctx, request)
}
