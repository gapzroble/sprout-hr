package sprout_test

import (
	"fmt"
	"testing"

	"github.com/gapzroble/sprout-hr/pkg/sprout"
)

func TestGetToken(t *testing.T) {
	token, _ := sprout.GetRequestVerificationToken()
	fmt.Println(token)
}
