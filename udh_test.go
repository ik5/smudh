package udh

import (
	"testing"
)

func TestValidBytesInputForstring(t *testing.T) {
	inputs := []UDHMessage{
		[]byte("0608047539040405d105d105e805db05d4002005d905d505e405d9002005e405d905e005e005e105d905dd002005d105e2002205de"),
	}

	var (
		result string
		err    error
	)
	for idx, input := range inputs {
		result, err = DecodeUDHSMSMessagemsg(input)
		if err != nil {
			t.Errorf("%d. expected %s to be valid, but error returned: %s", idx, input, err)
			continue
		}

		if result != "" {
			t.Error(result)
		}
	}
}
