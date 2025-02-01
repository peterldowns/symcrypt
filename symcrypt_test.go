package symcrypt_test

import (
	"testing"

	"github.com/peterldowns/symcrypt"
)

func TestHello(t *testing.T) {
	t.Parallel()
	if symcrypt.Hello() != "Hello world" {
		t.Fail()
	}
}
