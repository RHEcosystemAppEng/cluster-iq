package inventory

import "testing"

func TestAsInstanceState(t *testing.T) {
	tests := []struct {
		input  string
		result InstanceState
	}{
		{
			input:  "running",
			result: Running,
		},
		{
			input:  "RUNNING",
			result: Running,
		},
		{
			input:  "stop",
			result: Stopped,
		},
		{
			input:  "stopped",
			result: Stopped,
		},
		{
			input:  "terminated",
			result: Terminated,
		},
		{
			input:  "unknown",
			result: Unknown,
		},
		{
			input:  "RANDOM",
			result: Unknown,
		},
	}

	for _, test := range tests {
		result := AsInstanceState(test.input)
		if test.result != result {
			t.Errorf("Instance state parsing failed. Have: %s ; Expected: %v", result, test.result)
		}
	}

}
