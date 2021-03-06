// Automatically generated by MockGen. DO NOT EDIT!
// Source: notifier.go

package engine

import (
  gomock "code.google.com/p/gomock/gomock"
)

// Mock of Notifier interface
type MockNotifier struct {
  ctrl     *gomock.Controller
  recorder *_MockNotifierRecorder
}

// Recorder for MockNotifier (not exported)
type _MockNotifierRecorder struct {
  mock *MockNotifier
}

func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
  mock := &MockNotifier{ctrl: ctrl}
  mock.recorder = &_MockNotifierRecorder{mock}
  return mock
}

func (_m *MockNotifier) EXPECT() *_MockNotifierRecorder {
  return _m.recorder
}

func (_m *MockNotifier) Notify(who string, rule_name string, msg string) error {
  ret := _m.ctrl.Call(_m, "Notify", who, rule_name, msg)
  ret0, _ := ret[0].(error)
  return ret0
}

func (_mr *_MockNotifierRecorder) Notify(arg0, arg1, arg2 interface{}) *gomock.Call {
  return _mr.mock.ctrl.RecordCall(_mr.mock, "Notify", arg0, arg1, arg2)
}
