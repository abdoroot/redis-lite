package redis

import (
	"testing"
)

func TestUMarshal(t *testing.T) {
	resp, err := Unmarshal("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nmohamed\r\n")
	if err != nil {
		t.Error(err)
	}
	if len(resp.([]any)) != 3 {
		t.Error("expected slice of any with length = 3")
	}
}

func TestResp(t *testing.T) {
	testCases := []struct {
		Data, Expected string
	}{
		{"$-1\r\n", "bstring"},
		{"*1\r\n$4\r\nping\r\n", "array"},
		{"*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n", "array"}, //mean ["echo","hello world"]
		{"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", "array"},
		{"+OK\r\n", "sstring"},
		{"-Error message\r\n", "error"},
		{"$0\r\n\r\n", "bstring"},
		{"+hello world\r\n", "sstring"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Expected, func(t *testing.T) {
			got := dataType(testCase.Data)
			if testCase.Expected != got {
				t.Errorf("expected %v got %v", testCase.Expected, got)
			}
		})
	}
}
