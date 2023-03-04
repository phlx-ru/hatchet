package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVar(t *testing.T) {
	testCases := []struct {
		name     string
		variable any
		tag      string
		ok       bool
	}{
		{
			name:     "min_ok",
			variable: 5,
			tag:      "min=3",
			ok:       true,
		},
		{
			name:     "min_fail",
			variable: 5,
			tag:      "min=8",
			ok:       false,
		},
		{
			name:     "rfc3339_ok",
			variable: "2022-09-17T22:17:17.958458Z",
			tag:      "datetime",
			ok:       false,
		},
		{
			name:     "rfc3339_ok",
			variable: "2022-09-17T22:17:17.958458Z",
			tag:      "datetime=2006-01-02T15:04:05Z07:00",
			ok:       true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := Var(testCase.variable, testCase.tag)
			require.NoError(t, err)
			if testCase.ok {
				if len(actual) != 0 {
					require.NoError(t, actual[0])
				}
			} else {
				if len(actual) != 0 {
					require.Error(t, actual[0])
				} else {
					require.NotEmpty(t, actual)
				}
			}
		})
	}
}
