package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPrint(t *testing.T) {
	ss := []Service{
		{
			Name: "foo",
		},
		{
			Name: "bar",
		},
	}
	expected := "[foo bar]"
	actual := fmt.Sprintf("%v", ss)

	assert.Equal(t, expected, actual)
}
