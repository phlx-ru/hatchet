package maps

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

func TestExtractor(t *testing.T) {
	testCases := []struct {
		name         string
		value        map[string]any
		manipulation func(*Extractor) any
		expected     any
	}{
		{
			name: `int`,
			value: map[string]any{
				"count": 8,
			},
			manipulation: func(getter *Extractor) any {
				return getter.Get(`count`).ToInt()
			},
			expected: 8,
		},
		{
			name: `int-pointer`,
			value: map[string]any{
				"count": 8,
			},
			manipulation: func(getter *Extractor) any {
				return getter.Get(`count`).ToPointerInt()
			},
			expected: pointer.ToInt(8),
		},
		{
			name: `get-by-cargo-direct`,
			value: map[string]any{
				"name":      "John",
				"name_uid":  "Jane",
				"name_type": "jxs:string",
			},
			manipulation: func(getter *Extractor) any {
				return getter.GetByCargoType(`name`).ToString()
			},
			expected: "John",
		},
		{
			name: `get-by-cargo-fallback`,
			value: map[string]any{
				"name_uid":  "Jane",
				"name_type": "jxs:string",
			},
			manipulation: func(getter *Extractor) any {
				return getter.GetByCargoType(`name`).ToString()
			},
			expected: "Jane",
		},
		{
			name: `get-by-cargo-none`,
			value: map[string]any{
				"name_uid":  "Jane",
				"name_type": "nonexistent:string",
			},
			manipulation: func(getter *Extractor) any {
				return getter.GetByCargoType(`name`).ToString()
			},
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			getter, err := MakeExtractor(testCase.value)
			require.NoError(t, err)
			actual := testCase.manipulation(getter)
			require.Equal(t, testCase.expected, actual)
		})
	}
}
