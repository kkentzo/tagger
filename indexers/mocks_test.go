package indexers

import "github.com/stretchr/testify/mock"

type MockRvmHandler struct {
	mock.Mock
}

func (rvm *MockRvmHandler) IsRuby(path string) bool {
	args := rvm.Called(path)
	return args.Get(0).(bool)
}

func (rvm *MockRvmHandler) GemsetPath(path string) (string, error) {
	args := rvm.Called(path)
	gpath := args.Get(0).(string)
	err := args.Get(1)
	if err != nil {
		return gpath, err.(error)
	} else {
		return gpath, nil
	}
}
