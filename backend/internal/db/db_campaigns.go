package db

import (
	"encoding/json"
	"fmt"
)

type Scope string

type Campaign struct {
	Name      string `json:"name"`
	TenantId  string `json:"tenant_id"`
	UniqueURI string `json:"uniqueURI"`
}

func (d *DB) GetCampaigns() ([]Campaign, error) {
	rows, err := d.conn.Query(`SELECT name, tenant_id, UniqueStr FROM campaigns ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := []Campaign{}
	for rows.Next() {
		var c Campaign
		if err := rows.Scan(&c.Name, &c.TenantId, &c.UniqueURI); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, rows.Err()
}

func (d *DB) InsertNewCampaign(uniqueString, name, tenantId, resource string, scopes []Scope) error {
	scopesJSON, err := json.Marshal(scopes)
	if err != nil {
		return err
	}

	q := `INSERT INTO campaigns(UniqueStr, name, tenant_id, resource, scopes) VALUES(?, ?, ?, ?, ?)`
	_, err = d.conn.Exec(q, uniqueString, name, tenantId, resource, string(scopesJSON))
	return err

}

type CampaignMeta struct {
	Id          int      `json:"id"`
	UniqueStr   string   `json:"unique_str"`
	Name        string   `json:"name"`
	Tenant_id   string   `json:"tenant_id"`
	Resource    string   `json:"resource"`
	Clicks      int      `json:"clicks"`
	Scopes      []string `json:"scopes"`
	TotalTokens int      `json:"tokens_total"`
}

// `CREATE TABLE IF NOT EXISTS campaigns(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		UniqueStr TEXT NOT NULL UNIQUE,
// 		name TEXT NOT NULL UNIQUE,
// 		tenant_id TEXT NOT NULL UNIQUE,
// 		resource TEXT NOT NULL,
// 		clicks INTEGER,
// 		scopes TEXT NOT NULL CHECK (json_valid(scopes))
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS tokens(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		tenant_id TEXT NOT NULL,
// 		unqiueStr TEXT NOT NULL,
// 		access_token TEXT NOT NULL,
// 		refresh_token TEXT NOT NULL,
// 		FOREIGN KEY (unqiueStr) REFERENCES campaigns(UniqueStr) ON DELETE CASCADE
// 	);`, // just for testing, so can verify if works

func (d *DB) GetCampaignMeta(id string) (*CampaignMeta, error) {
	q := `
		SELECT
			c.id,
			c.UniqueStr,
			c.name,
			c.tenant_id,
			c.resource,
			COALESCE(c.clicks, 0) AS clicks,
			c.scopes,
			COUNT(t.id) AS token_count
		FROM campaigns c
		LEFT JOIN tokens t
			ON t.unqiueStr = c.UniqueStr
		WHERE c.UniqueStr = ?
		GROUP BY c.id`

	var c CampaignMeta
	var jScopes string
	err := d.conn.QueryRow(q, id).Scan(&c.Id, &c.UniqueStr, &c.Name, &c.Tenant_id, &c.Resource, &c.Clicks, &jScopes, &c.TotalTokens)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal([]byte(jScopes), &c.Scopes); err != nil {
		fmt.Println(err)

		return nil, err
	}

	return &c, nil

}
func (d *DB) EnsureLandingIdExists(id string) (bool, error) {
	var count int
	err := d.conn.QueryRow("SELECT COUNT(1) FROM campaigns WHERE UniqueStr = ? LIMIT 1", id).Scan(&count)
	return count > 1, err
}
