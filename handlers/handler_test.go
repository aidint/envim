package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockHandler struct {
	state         HandlerState
	errors        []error
	postExecState HandlerState
	shouldProceed bool
	dependencies  []HandlerType
}

func (mh *mockHandler) GetType() HandlerType {
	return BaseType
}

func (mh *mockHandler) GetState() HandlerState {
	return mh.state
}

func (mh *mockHandler) GetErrors() []error {
	return mh.errors
}

func (mh *mockHandler) Execute(state map[HandlerType]Handler) {
	mh.state = mh.postExecState
}

func (mh *mockHandler) ShouldProceed() bool {
	return mh.shouldProceed
}

func (mh *mockHandler) DependsOn() []HandlerType {
	return mh.dependencies
}

type endState struct {
	state      HandlerState
	errorCount int
	proceed    bool
}

func TestChainHandler(t *testing.T) {
	tests := []struct {
		mocks []mockHandler
		end   endState
	}{
		// Test 1
		{
			mocks: []mockHandler{
				{
					postExecState: HandlerSuccessState,
					shouldProceed: true,
				},
				{
					postExecState: HandlerSuccessState,
					shouldProceed: true,
				},
			},
			end: endState{
				state:      HandlerSuccessState,
				errorCount: 0,
				proceed:    true,
			},
		},
		// Test 2
		{
			mocks: []mockHandler{
				{
					postExecState: HandlerErrorState,
					shouldProceed: true,
				},
				{
					postExecState: HandlerSuccessState,
					shouldProceed: false,
				},
			},
			end: endState{
				state:      HandlerErrorState,
				errorCount: 0,
        proceed:    false,
			},
		},
	}

	for idx, test := range tests {
		ch := &ChainHandler{}
		for _, mh := range test.mocks {
			ch.AddHandler(&mh)
		}
		t.Run(fmt.Sprintf("Test %d", idx+1), func(t *testing.T) {
			t.Parallel()
			ch.Execute(nil)
			require.Equal(t,
				test.end.state,
				ch.GetState(),
				"Expected state: %s, returned state: %s",
				test.end.state,
				ch.GetState())

			require.Equal(t,
				test.end.errorCount,
				len(ch.GetErrors()),
				"Expected number of errors: %d, returned number of errors: %d",
				test.end.errorCount,
				len(ch.GetErrors()))
		})
	}

}
