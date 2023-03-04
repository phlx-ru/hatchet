package sanitize

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPhone(t *testing.T) {
	testCases := []struct {
		phone    string
		expected string
	}{
		{
			phone:    `9009009090`,
			expected: `9009009090`,
		},
		{
			phone:    `79009009090`,
			expected: `9009009090`,
		},
		{
			phone:    `89009009090`,
			expected: `9009009090`,
		},
		{
			phone:    `7 900 900 90 90`,
			expected: `9009009090`,
		},
		{
			phone:    ` 900  900  90 90  `,
			expected: `9009009090`,
		},
		{
			phone:    `[900]_900-90-90`,
			expected: `9009009090`,
		},
		{
			phone:    `+7 900 900 90 90`,
			expected: `9009009090`,
		},
		{
			phone:    `+7_[900]_900-90-90`,
			expected: `9009009090`,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			`test_phone:`+testCase.phone,
			func(t *testing.T) {
				actual := Phone(testCase.phone)
				require.Equal(t, testCase.expected, actual)
			},
		)
	}
}

func TestPhoneWithCountryCode(t *testing.T) {
	testCases := []struct {
		phone    string
		expected string
	}{
		{
			phone:    `9009009090`,
			expected: `79009009090`,
		},
		{
			phone:    `79009009090`,
			expected: `79009009090`,
		},
		{
			phone:    `89009009090`,
			expected: `79009009090`,
		},
		{
			phone:    `7 900 900 90 90`,
			expected: `79009009090`,
		},
		{
			phone:    ` 900  900  90 90  `,
			expected: `79009009090`,
		},
		{
			phone:    `[900]_900-90-90`,
			expected: `79009009090`,
		},
		{
			phone:    `+7 900 900 90 90`,
			expected: `79009009090`,
		},
		{
			phone:    `+7_[900]_900-90-90`,
			expected: `79009009090`,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			`test_phone:`+testCase.phone,
			func(t *testing.T) {
				actual := PhoneWithCountryCode(testCase.phone)
				require.Equal(t, testCase.expected, actual)
			},
		)
	}
}
