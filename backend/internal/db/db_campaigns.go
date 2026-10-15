package db

import "encoding/json"

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
