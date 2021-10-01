package eventer_test

import (
	"sync"
	"testing"

	. "github.com/andy2046/eventer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	MockEventListener struct {
		mock.Mock
	}

	testEvent struct{}

	EventerTestSuite struct {
		suite.Suite
		e EventEmitter
	}
)

func (m *MockEventListener) HandleEvent(e Event) {
	m.Called(e)
}

const methodName = "HandleEvent"

func (s *EventerTestSuite) TestRegisterEventListener() {
	var (
		wg1 sync.WaitGroup
		wg2 sync.WaitGroup
	)

	// expect listener to be invoked once
	wg1.Add(1)

	l1 := &MockEventListener{}
	l1.On(methodName, testEvent{}).Run(func(args mock.Arguments) {
		wg1.Done()
	}).Return()

	s.e.AddListener(l1)
	defer s.e.RemoveListener(l1)

	// expect listener to be invoked twice
	wg2.Add(2)

	l2 := &MockEventListener{}
	l2.On(methodName, testEvent{}).Run(func(args mock.Arguments) {
		wg2.Done()
	}).Return()

	s.e.AddListener(l2)
	defer s.e.RemoveListener(l2)

	s.e.EmitEvent(testEvent{})

	wg1.Wait()

	l1.AssertCalled(s.T(), methodName, testEvent{})
	l1.AssertNumberOfCalls(s.T(), methodName, 1)

	s.e.RemoveListener(l1)

	s.e.EmitEvent(testEvent{})

	wg2.Wait()

	l1.AssertCalled(s.T(), methodName, testEvent{})
	l1.AssertNumberOfCalls(s.T(), methodName, 1)

	l2.AssertNumberOfCalls(s.T(), methodName, 2)
}

func (s *EventerTestSuite) TestRegisterNilEventListener() {
	added := s.e.AddListener(nil)
	s.Assert().False(added, "expected to not add a nil listener")

	// to test that there is no nil pointer dereference when emitting an
	// event after a nil listener has been added
	s.e.EmitEvent(testEvent{})
}

func (s *EventerTestSuite) TestEmitEventWithoutEventListener() {
	// to test that emitting with empty EventListener list works.
	s.e.EmitEvent(testEvent{})
}

func (s *EventerTestSuite) TestDeregisterUnregisteredEventListener() {
	unregistered := &MockEventListener{}
	removed := s.e.RemoveListener(unregistered)
	s.Assert().False(removed, "expected to not remove a listener since it has never been added")
}

func (s *EventerTestSuite) TestDeregisterNilEventListener() {
	removed := s.e.RemoveListener(nil)
	s.Assert().False(removed, "expect to not remove a nil listener since it has not been added")

	s.e.AddListener(nil)

	removed = s.e.RemoveListener(nil)
	s.Assert().False(removed, "expect to not remove a nil listener since it can not be added")
}

func (s *EventerTestSuite) TestRemoveEventListenerDuringHandleEvent() {
	var wg1 sync.WaitGroup

	// expect listener to be invoked once
	wg1.Add(1)

	l1 := &MockEventListener{}
	l1.On(methodName, testEvent{}).Run(func(args mock.Arguments) {
		s.e.RemoveListener(l1)
		wg1.Done()
	}).Return()

	s.e.AddListener(l1)
	defer s.e.RemoveListener(l1)

	s.e.EmitEvent(testEvent{})

	wg1.Wait()

	l1.AssertCalled(s.T(), methodName, testEvent{})
	l1.AssertNumberOfCalls(s.T(), methodName, 1)

	s.e.EmitEvent(testEvent{})

	l1.AssertCalled(s.T(), methodName, testEvent{})
	l1.AssertNumberOfCalls(s.T(), methodName, 1)
}

func (s *EventerTestSuite) TestRegisterEventListenerTwice() {
	var wg1 sync.WaitGroup

	// expect listener to be invoked once
	wg1.Add(1)

	l1 := &MockEventListener{}
	l1.On(methodName, testEvent{}).Run(func(args mock.Arguments) {
		wg1.Done()
	}).Return()
	defer s.e.RemoveListener(l1)

	added := s.e.AddListener(l1)
	s.Assert().True(added, "expected to add a listener on the first time")

	added = s.e.AddListener(l1)
	s.Assert().False(added, "expected to not add a listener on the second time")

	s.e.EmitEvent(testEvent{})

	wg1.Wait()

	l1.AssertCalled(s.T(), methodName, testEvent{})
	l1.AssertNumberOfCalls(s.T(), methodName, 1)
}

func TestSyncEventEmitter(t *testing.T) {
	suite.Run(t, &EventerTestSuite{
		e: &SyncEventEmitter{},
	})
}

func TestAsyncEventEmitter(t *testing.T) {
	suite.Run(t, &EventerTestSuite{
		e: &AsyncEventEmitter{},
	})
}
