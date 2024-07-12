package amnesia

import "testing"

func TestFileWrite(t *testing.T) {
	fp := filePersistor{}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()

	fp.Write("abc", 123)
}
