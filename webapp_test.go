package eadmin

import (
	"log"
	"os"
	"testing"

	"github.com/GPORTALcloud/ouidb/pkg/ouidb"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func newTestWebApp(t *testing.T) *WebApp {
	db, _ := ouidb.New()
	return &WebApp{
		clock:        clock.NewMock(),
		ouidb:        db,
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
		_Topic:   "espresense/rooms/room_name/telemetry",
		_Payload: []byte(`{"queried": 2}`),
	}
	wa.mqttHandlerRoom(nil, &m)

	val, ok := wa.rooms["room_name"]
	assert.True(t, ok)
	assert.Equal(t, val.Telemetry.Queried, 2)
	assert.Equal(t, val.Telemetry.Recieved, wa.clock.Now())
}

const testMac = "000181dead01"

func TestHandleRoomDevice(t *testing.T) {

	wa := newTestWebApp(t)
	m := MockMessage{
		_Topic:   "espresense/rooms/room_name",
		_Payload: []byte(`{"mac": "000181dead01"}`),
	}
	wa.mqttHandlerRoom(nil, &m)

	room, ok := wa.rooms["room_name"]
	assert.True(t, ok)
	assert.Equal(t, room.LastPing, wa.clock.Now())
	assert.Equal(t, room.pingsByMac[testMac].Recieved, wa.clock.Now())

	assert.Equal(t, wa.manufByMac[testMac], "NortelNe")
}

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Exit(m.Run())
}
