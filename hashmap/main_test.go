package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMap(t *testing.T) {
	cm := NewCacheMap(5)

	// Добавление значений
	cm.Set("one", "1")       // новый элемент
	cm.Set("two", "2")       // коллизия
	cm.Set("three", "3")     // новый элемент
	cm.Set("four", "4")      // новый элемент
	cm.Set("five", "5")      // коллизия
	cm.Set("six", "6")       // коллизия
	cm.Set("seven", "seven") // коллизия

	v, ok := cm.Get("three")
	assert.True(t, ok)
	assert.Equal(t, "3", v)

	v, ok = cm.Get("one")
	assert.True(t, ok)
	assert.Equal(t, "1", v)

	v, ok = cm.Get("two")
	assert.True(t, ok)
	assert.Equal(t, "2", v)

	v, ok = cm.Get("seven")
	assert.True(t, ok)
	assert.Equal(t, "seven", v)

	// Удаление
	cm.Delete("three")

	// Проверка после удаления
	v, ok = cm.Get("three")
	assert.False(t, ok)
	assert.Empty(t, v)

	cm.Set("seven", "7") // обновление

	// Проверка после редактирования
	v, ok = cm.Get("seven")
	assert.True(t, ok)
	assert.Equal(t, "7", v)
}
