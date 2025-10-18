package cache

import (
	"testing"

	"github.com/dark-vinci/nildb/constants"
)

// TestBuilder verifies the cache builder configuration
func TestBuilder(t *testing.T) {
	tests := []struct {
		name             string
		configure        func(*Builder)
		expectedMaxSize  uint
		expectedPageSize uint
		expectedPinLimit float32
		expectedK        uint
		expectedCRP      uint64
		expectPanic      bool
	}{
		{
			name:             "Default configuration",
			configure:        func(b *Builder) {},
			expectedMaxSize:  constants.DefaultMaxCacheSize,
			expectedPageSize: constants.DefaultPageSize,
			expectedPinLimit: constants.DefaultPinPercentageLimit,
			expectedK:        constants.DefaultLruK,
			expectedCRP:      constants.DefaultCRP,
		},
		{
			name: "Custom configuration",
			configure: func(b *Builder) {
				b.SetMaxSize(100).SetPageSize(8192).SetPinPercentageLimit(75.0).LruK(3).CorrelatedReferencePeriod(100)
			},
			expectedMaxSize:  100,
			expectedPageSize: 8192,
			expectedPinLimit: 75.0,
			expectedK:        3,
			expectedCRP:      100,
		},
		{
			name: "Invalid max size",
			configure: func(b *Builder) {
				b.SetMaxSize(constants.MinCacheSize - 1)
			},
			expectPanic: true,
		},
		{
			name: "Invalid pin percentage limit",
			configure: func(b *Builder) {
				b.SetPinPercentageLimit(101.0)
			},
			expectPanic: true,
		},
		{
			name: "Invalid K value",
			configure: func(b *Builder) {
				b.LruK(0)
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic, but none occurred")
					}
				}()
			}

			tt.configure(b)

			if !tt.expectPanic {
				c := b.Build()

				if c.MaxSize != tt.expectedMaxSize {
					t.Errorf("Expected MaxSize %d, got %d", tt.expectedMaxSize, c.MaxSize)
				}

				if c.PageSize != tt.expectedPageSize {
					t.Errorf("Expected PageSize %d, got %d", tt.expectedPageSize, c.PageSize)
				}

				if c.PinPercentageLimit != tt.expectedPinLimit {
					t.Errorf("Expected PinPercentageLimit %f, got %f", tt.expectedPinLimit, c.PinPercentageLimit)
				}

				if c.K != tt.expectedK {
					t.Errorf("Expected K %d, got %d", tt.expectedK, c.K)
				}

				if c.CRP != tt.expectedCRP {
					t.Errorf("Expected CRP %d, got %d", tt.expectedCRP, c.CRP)
				}
			}
		})
	}
}
