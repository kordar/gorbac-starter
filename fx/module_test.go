package gorbacstarterfx

import (
	"testing"

	gocfgmodulefx "github.com/kordar/gocfg-load-module/fx/v2"
	"github.com/kordar/gorbac"
)

func TestModule_NormalizeConfig(t *testing.T) {
	t.Parallel()

	cfg := normalizeConfig(ModuleConfig{})
	if cfg.DB != "gorbac" {
		t.Fatalf("expected default DB 'gorbac', got %q", cfg.DB)
	}
	if cfg.Driver != "mysql" {
		t.Fatalf("expected default Driver 'mysql', got %q", cfg.Driver)
	}
	if cfg.GuestRole != "guest" {
		t.Fatalf("expected default GuestRole 'guest', got %q", cfg.GuestRole)
	}
	if cfg.TableRule != "sys_auth_rule" {
		t.Fatalf("expected default TableRule 'sys_auth_rule', got %q", cfg.TableRule)
	}
	if cfg.CacheTTL == 0 {
		t.Fatalf("expected non-zero CacheTTL")
	}
}

func TestModuleConfig(t *testing.T) {
	t.Parallel()

	cfg := ModuleConfig{
		DB:         "sys",
		Driver:     "mysql",
		GuestRole:  "anonymous",
		CacheStore: "redis",
	}
	if cfg.Driver != "mysql" {
		t.Fatalf("expected Driver 'mysql', got %q", cfg.Driver)
	}
}

func TestBuildModuleConfig(t *testing.T) {
	data := map[string]any{
		"driver":       "mysql",
		"db":           "sys",
		"guest_role":   "anonymous",
		"cache_enabled": "true",
		"cache_store":  "redis",
		"cache_ttl":    "5m",
		"table_rule":   "my_auth_rule",
	}
	cfg := buildModuleConfig(data)

	if cfg.Driver != "mysql" {
		t.Fatalf("expected Driver 'mysql', got %q", cfg.Driver)
	}
	if cfg.DB != "sys" {
		t.Fatalf("expected DB 'sys', got %q", cfg.DB)
	}
	if cfg.GuestRole != "anonymous" {
		t.Fatalf("expected GuestRole 'anonymous', got %q", cfg.GuestRole)
	}
	if !cfg.CacheEnabled {
		t.Fatal("expected CacheEnabled to be true")
	}
	if cfg.CacheStore != "redis" {
		t.Fatalf("expected CacheStore 'redis', got %q", cfg.CacheStore)
	}
	if cfg.CacheTTL.String() != "5m0s" {
		t.Fatalf("expected CacheTTL '5m0s', got %q", cfg.CacheTTL.String())
	}
	if cfg.TableRule != "my_auth_rule" {
		t.Fatalf("expected TableRule 'my_auth_rule', got %q", cfg.TableRule)
	}
}

func TestBuildModuleConfig_Empty(t *testing.T) {
	cfg := buildModuleConfig(nil)
	if cfg.Driver != "" {
		t.Fatalf("expected empty Driver, got %q", cfg.Driver)
	}
}

func TestStarterModule_ImplementsGoCfgModule(t *testing.T) {
	m := StarterModule("gorbac")
	if m.Name() != "gorbac" {
		t.Fatalf("expected name 'gorbac', got %q", m.Name())
	}

	var _ gocfgmodulefx.GoCfgModule = m
}

func TestStarterModule_Load(t *testing.T) {
	data := map[string]any{
		"driver":       "mysql",
		"db":           "sys",
	}
	m := StarterModule("gorbac")
	opts := m.Load(data)

	if len(opts) != 1 {
		t.Fatalf("expected 1 fx.Option, got %d", len(opts))
	}
}

var _ *gorbac.RbacService = nil // ensure gorbac import is used
