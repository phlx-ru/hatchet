package texts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlural(t *testing.T) {
	testCases := []struct {
		name             string
		count            int
		formSingular     string
		formPluralWeak   string
		formPluralStrong string
		expected         string
	}{
		{
			name:             "default",
			count:            1,
			formSingular:     "кот",
			formPluralWeak:   "кота",
			formPluralStrong: "котов",
			expected:         "кот",
		},
		{
			name:             "default",
			count:            187632,
			formSingular:     "кот",
			formPluralWeak:   "кота",
			formPluralStrong: "котов",
			expected:         "кота",
		},
		{
			name:             "default",
			count:            63463466,
			formSingular:     "кот",
			formPluralWeak:   "кота",
			formPluralStrong: "котов",
			expected:         "котов",
		},
		{
			name:             "default",
			count:            1,
			formSingular:     "рыб",
			formPluralWeak:   "рыба",
			formPluralStrong: "рыбов",
			expected:         "рыб",
		},
		{
			name:             "default",
			count:            2,
			formSingular:     "рыб",
			formPluralWeak:   "рыба",
			formPluralStrong: "рыбов",
			expected:         "рыба",
		},
		{
			name:             "default",
			count:            5,
			formSingular:     "рыб",
			formPluralWeak:   "рыба",
			formPluralStrong: "рыбов",
			expected:         "рыбов",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Plural(testCase.count, testCase.formSingular, testCase.formPluralWeak, testCase.formPluralStrong)
			require.Equal(t, testCase.expected, actual)
		})
	}
}
