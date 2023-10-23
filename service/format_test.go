package service

import (
	"testing"

	"github.com/onetooler/bistory-backend/test"
	"github.com/stretchr/testify/assert"
)

func TestFindAllFormats_Success(t *testing.T) {
	container := test.PrepareForServiceTest()

	service := NewFormatService(container)
	result := service.FindAllFormats()

	assert.Len(t, *result, 2)
}
