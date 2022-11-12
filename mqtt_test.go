package eadmin

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MockMessage struct {
	_Duplicate bool
	_Qos       byte
	_Retained  bool
	_Topic     string
	_MessageID uint16
	_Payload   []byte
	_Acked     bool
}

func (m *MockMessage) Duplicate() bool {
	return m._Duplicate
}

func (m *MockMessage) Qos() byte {
	return m._Qos
}

func (m *MockMessage) Retained() bool {
	return m._Retained
}

func (m *MockMessage) Topic() string {
	return m._Topic
}

func (m *MockMessage) MessageID() uint16 {
	return m._MessageID
}

func (m *MockMessage) Payload() []byte {
	return m._Payload
}

func (m *MockMessage) Ack() {
	m._Acked = true
}

var _ MQTT.Message = (*MockMessage)(nil)
