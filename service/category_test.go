package service

import (
	"testing"

	"github.com/onetooler/bistory-backend/test"
	"github.com/stretchr/testify/assert"
)

func TestFindAllCategories_Success(t *testing.T) {
	container := test.PrepareForServiceTest()

	service := NewCategoryService(container)
	result := service.FindAllCategories()

	assert.Len(t, *result, 3)
}
