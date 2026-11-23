package strategyregistry

import "testing"

func TestResolveBuiltInStrategies(t *testing.T) {
	registry := NewInMemoryCalculationStrategyRegistry()

	for _, name := range []string{"add", "subtract", "multiply", "divide"} {
		if operation, ok := registry.Resolve(name); !ok || operation == nil {
			t.Fatalf("Resolve(%q) = (%v, %v), want a strategy and true", name, operation, ok)
		}
	}

	if operation, ok := registry.Resolve("unknown"); ok || operation != nil {
		t.Fatalf("Resolve(unknown) = (%v, %v), want (nil, false)", operation, ok)
	}
}
