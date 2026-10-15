package codes

import (
	"fmt"
	"sync"
	"time"

	"github.com/z3vxo/TokenJack/internal/db"
)

type ActivePolls struct {
	expires_at time.Time
	DeviceCode string
	Tenat      string
	Status     int // 0 = polling, 1 = failed, 2 = success
}

type Codes struct {
	mu     sync.Mutex
	Active []ActivePolls
	DB     *db.DB
}

func New(d *db.DB) *Codes {
	c := &Codes{DB: d}
	go c.Purger()
	return c
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

func (c *Codes) UpdateStatus(dcr EntraResp, status int) {

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

func (c *Codes) InsertAndBeginPoll(CodeInfo db.CodeInfo) error {

	EntraResponse, err := c.PostEntraDeviceCode(CodeInfo)
	if err != nil {
		return err
	}

	fmt.Println(EntraResponse.UserCode, EntraResponse.VerifyUrl)

	c.Insert(EntraResponse)
	go c.PollDeviceCode(EntraResponse)

	return nil
}
