package inventory

import "testing"

func TestAsInstanceStatus(t *testing.T) {
	tests := []struct {
		input  string
		result InstanceStatus
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
			input:  "RANDOM",
			result: Running,
		},
	}

	for _, test := range tests {
		result := AsInstanceStatus(test.input)
		if test.result != result {
			t.Errorf("Instance status parsing failed. Have: %s ; Expected: %v", result, test.result)
		}
	}

}
