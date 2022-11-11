package eadmin

import (
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func newTestWebApp(t *testing.T) *WebApp {
	return &WebApp{
		clock:        clock.NewMock(),
		deviceByName: map[string]*Device{},
		rooms:        map[string]*Room{},
		manufByMac:   map[string]string{},
	}
}

func TestGetOrInsertRoom(t *testing.T) {
	wa := newTestWebApp(t)
	wa.mu.Lock()
	defer wa.mu.Unlock()
	room := wa.getOrInsertRoomLocked("new_room")
	assert.NotNil(t, room)
	assert.NotNil(t, room.pingsByMac)

	{
		val, ok := wa.rooms["new_room"]
		assert.True(t, ok)
		assert.Equal(t, val, room)
	}

	{
		val, ok := wa.rooms["invalid_room"]
		assert.False(t, ok)
		assert.Nil(t, val)
	}
}

func TestHandleRoomTelem(t *testing.T) {
	wa := newTestWebApp(t)

	m := MockMessage{
		_Payload: []byte(`{"queried": 2}`),
	}

	wa.mqttHandlerRoomTelem(nil, &m, "room_name")

	val, ok := wa.rooms["room_name"]
	assert.True(t, ok)
	assert.Equal(t, val.Telemetry.Queried, 2)
	assert.Equal(t, val.Telemetry.Recieved, wa.clock.Now())
}
