//go:generate mockery -name=Injector -inpkg -output=default_mock_test -testonly -case=underscore

package inject

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

var (
	inj *MockInjector
)

func TestDefaultProvide(t *testing.T) {
	provider := func() int {
		return 42
	}
	inj = &MockInjector{}
	inj.
		On("Provide", mock.AnythingOfType("func() int")).
		Return(nil).
		Once()

	defaultInjector = inj

	Provide(provider)

	inj.AssertExpectations(t)
}

func TestDefaultMap(t *testing.T) {
	var val int
	inj = &MockInjector{}
	inj.
		On("Map", mock.AnythingOfType("int")).
		Return(nil).
		Once()

	defaultInjector = inj

	Map(val)

	inj.AssertExpectations(t)
}

func TestDefaultMapTo(t *testing.T) {
	var a, b int
	inj = &MockInjector{}
	inj.
		On("MapTo", mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Return(nil).
		Once()

	defaultInjector = inj

	MapTo(a, b)

	inj.AssertExpectations(t)
}

func TestDefaultConstruct(t *testing.T) {
	var someValue int
	inj = &MockInjector{}
	inj.
		On("Construct", mock.AnythingOfType("")).
		Return(nil).
		Once()

	defaultInjector = inj

	Construct(&someValue)

	inj.AssertExpectations(t)
}

func TestDefaultConstructLater(t *testing.T) {
	var someValue int
	inj = &MockInjector{}
	inj.
		On("ConstructLater", mock.AnythingOfType("")).
		Return(nil).
		Once()

	defaultInjector = inj

	ConstructLater(&someValue)

	inj.AssertExpectations(t)
}

func TestDefaultFinishConstruct(t *testing.T) {
	inj = &MockInjector{}
	inj.
		On("FinishConstruct").
		Return(nil).
		Once()

	defaultInjector = inj

	FinishConstruct()

	inj.AssertExpectations(t)
}
