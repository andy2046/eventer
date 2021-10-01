# eventer

[![Documentation](https://godoc.org/github.com/andy2046/eventer?status.svg)](http://godoc.org/github.com/andy2046/eventer)
[![GitHub issues](https://img.shields.io/github/issues/andy2046/eventer.svg)](https://github.com/andy2046/eventer/issues)
[![license](https://img.shields.io/github/license/andy2046/eventer.svg)](https://github.com/andy2046/eventer/LICENSE)
[![Release](https://img.shields.io/github/release/andy2046/eventer.svg?label=Release)](https://github.com/andy2046/eventer/releases)

----

## event emitter made easy

## Install

```
go get github.com/andy2046/eventer
```

## Usage

```go
package main

import (
	"github.com/andy2046/eventer"
)

type MockEventListener struct{}

func (m *MockEventListener) HandleEvent(e eventer.Event) {
	println("HandleEvent", e)
}

type testEvent struct{}

func main() {
	l := &MockEventListener{}
	emitter := &eventer.SyncEventEmitter{}
	// emitter := &eventer.AsyncEventEmitter{}

	emitter.AddListener(l)
	defer emitter.RemoveListener(l)

	emitter.EmitEvent(testEvent{})
}
```
