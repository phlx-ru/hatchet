package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractDatabaseName(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected string
		wantErr  bool
	}{
		{
			name:     `storage`,
			source:   `host=127.0.0.1 port=5432 user=postgres password=postgres dbname=storage sslmode=disable`,
			expected: `storage`,
		},
		{
			name:     `storage_without_sslmode`,
			source:   `host=127.0.0.1 port=5432 user=postgres password=postgres dbname=storage`,
			expected: `storage`,
		},
		{
			name:     `storage_starts_with`,
			source:   `dbname=storage host=127.0.0.1 port=5432 user=postgres password=postgres`,
			expected: `storage`,
		},
		{
			name:     `postgres`,
			source:   `host=127.0.0.1 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable`,
			expected: defaultExistedDatabase,
		},
		{
			name:    `empty`,
			source:  `host=127.0.0.1 port=5432 user=postgres password=postgres dbname= sslmode=disable`,
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := extractDatabaseNameFromSource(testCase.source)
			if testCase.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expected, actual)
			}
		})
	}
}
