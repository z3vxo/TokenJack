package codes

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/z3vxo/TokenJack/internal/db"
	"github.com/z3vxo/TokenJack/internal/events"
)

type ActivePolls struct {
	expires_at time.Time
	DeviceCode string
	Tenat      string
	UnqiueStr  string
	Status     int // 0 = polling, 1 = failed, 2 = success
}

type Codes struct {
	mu     sync.Mutex
	Active []ActivePolls
	DB     *db.DB
	Hub    *events.Hub
}

func New(d *db.DB, hub *events.Hub) *Codes {
	c := &Codes{DB: d, Hub: hub}
	go c.Purger()
	return c
}

type TokenEvent struct {
	Type      string `json:"type"`
	UniqueStr string `json:"unique_str"`
	Message   string `json:"message"`
}

func (c *Codes) publish(uniqueStr, message string) {
	if c.Hub == nil {
		return
	}
	payload, err := json.Marshal(TokenEvent{
		Type:      "token",
		UniqueStr: uniqueStr,
		Message:   message,
	})
	if err != nil {
		return
	}
	_ = c.Hub.Send(uniqueStr, payload)
}

func isExpired(code ActivePolls) bool {

	return !time.Now().Before(code.expires_at)
}

func (c *Codes) Insert(dcr EntraResp) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	new := ActivePolls{
		expires_at: time.Now().Add(time.Duration(dcr.Expires) * time.Second),
		DeviceCode: dcr.DeviceCode,
		Tenat:      dcr.Tenant,
		UnqiueStr:  dcr.UnqiueStr,
		Status:     0,
	}

	c.Active = append(c.Active, new)
	return nil
}

func (c *Codes) Remove(dcr EntraResp) {
	c.mu.Lock()
	defer c.mu.Unlock()
	active := c.Active[:0]
	for _, poll := range c.Active {
		if poll.DeviceCode != dcr.DeviceCode {
			active = append(active, poll)
		}
	}
	c.Active = active
}

func (c *Codes) Purger() {
	for {
		time.Sleep(2 * time.Minute)
		c.mu.Lock()

		active := c.Active[:0]
		for _, poll := range c.Active {
			if !isExpired(poll) {
				active = append(active, poll)
			}
		}
		c.Active = active
		c.mu.Unlock()
	}
}

func (c *Codes) InsertAndBeginPoll(CodeInfo db.CodeInfo) (string, string, error) {

	EntraResponse, err := c.PostEntraDeviceCode(CodeInfo)
	if err != nil {
		return "", "", err
	}

	fmt.Println(EntraResponse.UserCode, EntraResponse.VerifyUrl)

	c.Insert(EntraResponse)
	go c.PollDeviceCode(EntraResponse)

	return EntraResponse.UserCode, EntraResponse.VerifyUrl, nil
}
