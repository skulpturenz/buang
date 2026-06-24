package services

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockService struct{}

func TestSetAndGetByType(t *testing.T) {
	svc := New()
	mock := &mockService{}

	key := reflect.TypeFor[mockService]()
	Set(svc, key, mock)

	got, ok := Get[any](svc, key)

	require.True(t, ok, "expected to get service by type")
	require.Equalf(t, got.Unwrap(), mock, "got %v, want %v", got.Unwrap(), mock)
}

func TestSetAndGetByString(t *testing.T) {
	svc := New()
	mock := &mockService{}

	Set(svc, "mock", mock)

	got, ok := Get[any](svc, "mock")

	require.True(t, ok, "expected to get service by string")
	require.Equalf(t, got.Unwrap(), mock, "got %v, want %v", got.Unwrap(), mock)
}

func TestGetNonExistentByType(t *testing.T) {
	svc := New()
	_, ok := Get[any](svc, reflect.TypeFor[mockService]())

	require.False(t, ok, "expected not found for non-existent key")
}

func TestGetNonExistentByString(t *testing.T) {
	svc := New()
	_, ok := Get[any](svc, "non-existent")

	require.False(t, ok, "expected not found for non-existent key")
}

func TestForRangeEntries(t *testing.T) {
	svc := New()
	Set(svc, "key1", "value1")
	Set(svc, "key2", "value2")

	keys := []string{}
	values := []string{}
	for k, v := range svc.Entries() {
		keys = append(keys, k.(string))
		values = append(values, v.(string))
	}

	require.Equal(t, len(keys), 2)
	require.Equal(t, len(values), 2)
}
