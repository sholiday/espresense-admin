package eadmin

import (
	"encoding/json"
	"math"

	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type WebApp struct {
	config Config
	engine *gin.Engine

	// Synced datastructures
	mu           sync.Mutex
	deviceByName map[string]*Device
	rooms        map[string]*Room
}

type Device struct {
	Pings map[string]*Ping
}

type Room struct {
	LastPing   time.Time
	pingsByMac map[string]*Ping
}

func NewWebApp(c Config) (*WebApp, error) {
	return &WebApp{
		config:       c,
		engine:       gin.Default(),
		deviceByName: map[string]*Device{},
		rooms:        map[string]*Room{},
	}, nil
}

func (a *WebApp) addRoutes() {
	a.engine.LoadHTMLGlob("*.html")
	a.engine.StaticFile("/table.js", "./table.js")
	a.engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	a.engine.GET("/plot", func(c *gin.Context) {
		c.HTML(200, "plot.html", nil)
	})
	a.engine.GET("/table", func(c *gin.Context) {
		t := a.toTableResponse()
		c.JSON(http.StatusOK, t)
	})
	a.engine.GET("/3d", func(c *gin.Context) {
		t := a.toThreeD()
		c.JSON(http.StatusOK, t)
	})
}

func (a *WebApp) setupMqtt() error {
	opts := MQTT.NewClientOptions().AddBroker(a.config.Broker.Server)
	opts.SetUsername(a.config.Broker.Username)
	opts.SetPassword(a.config.Broker.Password)
	opts.SetClientID(a.config.Broker.ClientID)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := c.Subscribe("espresense/devices/#", 0, a.mqttHandlerDevice); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := c.Subscribe("espresense/rooms/+", 0, a.mqttHandlerRoom); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (a *WebApp) Run() error {
	err := a.setupMqtt()
	if err != nil {
		return err
	}
	a.addRoutes()
	return a.engine.Run(fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port))
}

func (a *WebApp) getOrInsertRoomLocked(roomName string) *Room {
	rEntry := a.rooms[roomName]
	if rEntry == nil {
		rEntry = &Room{
			LastPing:   time.Time{},
			pingsByMac: map[string]*Ping{},
		}
		a.rooms[roomName] = rEntry
	}
	return rEntry
}

func (a *WebApp) mqttHandlerRoom(client MQTT.Client, msg MQTT.Message) {
	t := time.Now()
	topic := msg.Topic()

	tParts := strings.Split(topic, "/")
	roomName := tParts[2]

	var ping Ping
	json.Unmarshal(msg.Payload(), &ping)
	ping.Recieved = t

	a.mu.Lock()
	defer a.mu.Unlock()
	room := a.getOrInsertRoomLocked(roomName)
	room.LastPing = t
	room.pingsByMac[ping.Mac] = &ping
}

func (a *WebApp) mqttHandlerDevice(client MQTT.Client, msg MQTT.Message) {
	t := time.Now()
	topic := msg.Topic()
	tParts := strings.Split(topic, "/")
	id := tParts[2]
	room := tParts[3]

	var ping Ping
	json.Unmarshal(msg.Payload(), &ping)
	ping.Recieved = t

	a.mu.Lock()
	defer a.mu.Unlock()
	rEntry := a.getOrInsertRoomLocked(room)
	rEntry.LastPing = t
	entry := a.deviceByName[id]
	if entry == nil {
		entry = &Device{}
	}
	if entry.Pings == nil {
		entry.Pings = make(map[string]*Ping)
	}
	entry.Pings[room] = &ping
	a.deviceByName[id] = entry
}

func (a *WebApp) gcLocked() {
	// GC deviceByName
	gcThreshold := time.Now().Add(-30 * time.Second)
	for id, device := range a.deviceByName {
		var toGC []string
		for room, ping := range device.Pings {
			if ping.Recieved.Before(gcThreshold) {
				toGC = append(toGC, room)
			}
		}
		for _, room := range toGC {
			delete(device.Pings, room)
			log.Printf("GC %s %s", id, room)
		}
	}

	// GC pingsByMac
	for rName, room := range a.rooms {
		var toGC []string
		for mac, ping := range room.pingsByMac {
			if ping.Recieved.Before(gcThreshold) {
				toGC = append(toGC, mac)
			}
		}
		for _, mac := range toGC {
			delete(room.pingsByMac, mac)
			log.Printf("GC %s %s", rName, mac)
		}
	}

	// GC rooms
	gcThreshold = time.Now().Add(time.Hour)
}

func (a *WebApp) toTableResponse() TableResponse {
	var t TableResponse
	a.mu.Lock()
	defer a.mu.Unlock()
	a.gcLocked()
	for room := range a.rooms {
		t.Rooms = append(t.Rooms, room)
	}
	for id, device := range a.deviceByName {
		var e = Entry{
			"name": id,
		}
		closest := ""
		dist := math.MaxFloat64
		for room, ping := range device.Pings {
			if ping.Distance < dist {
				closest = room
				dist = ping.Distance
			}
		}
		e["closest"] = closest
		for room := range a.rooms {
			if val, ok := device.Pings[room]; ok {
				e[room] = fmt.Sprintf("%.1f", val.Distance)
			} else {
				e[room] = ""
			}
		}
		t.Data = append(t.Data, e)
	}
	return t
}

func (a *WebApp) toThreeD() ThreeD {
	var t ThreeD
	for room := range a.rooms {
		t.Nodes = append(t.Nodes, ThreeDNode{
			ID:    room,
			Group: 1,
		})
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	mac := map[string]bool{}
	for _, device := range a.deviceByName {
		for room, ping := range device.Pings {
			mac[ping.Mac] = true
			t.Links = append(t.Links, ThreeDLink{
				Source: room,
				Target: ping.Mac,
				Value:  int(ping.Distance),
			})
		}
	}

	for m := range mac {
		t.Nodes = append(t.Nodes, ThreeDNode{
			ID:    m,
			Group: 2,
		})
	}
	return t
}