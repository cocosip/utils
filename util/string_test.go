package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnonymous(t *testing.T) {
	v1 := Anonymous("张三", 3, '*')
	assert.Equalf(t, "张**", v1, "anonymous value not equal %s", v1)

	v2 := Anonymous("张三丰", 3, '*')
	assert.Equalf(t, "张*丰", v2, "anonymous value not equal %s", v2)

	v3 := Anonymous("李世民", 4, 'x')
	assert.Equalf(t, "李xx民", v3, "anonymous value not equal %s", v3)

}
