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
	handlerType   HandlerType
}

func (mh *mockHandler) GetType() HandlerType {
	return mh.handlerType
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
		mocks      []mockHandler
		end        endState
		initExeMap map[HandlerType]Handler
	}{
		// Case 1
		{
			mocks: []mockHandler{
				{
					postExecState: HandlerSuccessState,
					shouldProceed: true,
					handlerType:   -1,
				},
				{
					postExecState: HandlerSuccessState,
					shouldProceed: true,
					handlerType:   -2,
				},
			},
			end: endState{
				state:      HandlerSuccessState,
				errorCount: 0,
				proceed:    true,
			},
		},
		// Case 2
		{
			mocks: []mockHandler{
				{
					postExecState: HandlerErrorState,
					shouldProceed: true,
					errors:        []error{fmt.Errorf("Error 1")},
					handlerType:   -1,
				},
				{
					postExecState: HandlerSuccessState,
					errors:        []error{fmt.Errorf("Error 2")},
					shouldProceed: true,
					handlerType:   -2,
					dependencies:  []HandlerType{-1},
				},
			},
			end: endState{
				state:      HandlerSuccessState,
				errorCount: 2,
				proceed:    true,
			},
		},
		{
			mocks: []mockHandler{
				{
					postExecState: HandlerSuccessState,
					shouldProceed: true,
					handlerType:   -1,
				},
				{
					postExecState: HandlerSuccessState,
					errors:        []error{fmt.Errorf("Error 2")},
					shouldProceed: false,
					handlerType:   -2,
					dependencies:  []HandlerType{-1, -3},
				},
			},
			end: endState{
				state:      HandlerErrorState,
				errorCount: 1,
				proceed:    false,
			},
			initExeMap: map[HandlerType]Handler{
				-3: nil,
			},
		},
	}

	for idx, test := range tests {
		ch := &ChainHandler{ExeMap: test.initExeMap}
		for _, mh := range test.mocks {
			ch.AddHandler(&mh)
		}
		t.Run(fmt.Sprintf("Test %d", idx+1), func(t *testing.T) {
			t.Parallel()

			require.Equal(t,
				ch.GetType(),
				ChainType,
				"Expected type: %s, returned type: %s",
				CheckEnvironmentType,
				ch.GetType())

			require.Equal(t,
				ch.DependsOn(),
				[]HandlerType{},
				"Expected dependencies: %s, returned dependencies: %s",
				[]HandlerType{},
				ch.DependsOn())

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

			require.Equal(t,
				ch.ShouldProceed(),
				test.end.proceed,
				"Expected proceed: %t, returned proceed: %t",
				test.end.proceed,
				ch.ShouldProceed())

			ekeys := make([]HandlerType, 0, len(ch.ExeMap))
			for ht := range ch.ExeMap {
				ekeys = append(ekeys, ht)
			}

			mockKeys := make([]HandlerType, 0, len(test.mocks))
			for _, mh := range test.mocks {
				mockKeys = append(mockKeys, mh.handlerType)
			}

			require.ElementsMatch(t,
				ekeys,
				mockKeys,
				"Expected execMap keys: %v, returned execMap keys: %v",
				mockKeys,
				ekeys,
			)

		})
	}

}
