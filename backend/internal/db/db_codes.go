package db

import "encoding/json"

type CodeInfo struct {
	Resource string
	Scopes   []Scope
	TenantId string
}

func (d *DB) GetCodeInfo(id string) (CodeInfo, error) {

	var code CodeInfo
	var scopes string
	err := d.conn.QueryRow(`SELECT resource, tenant_id, scopes FROM campaigns WHERE UniqueStr = ?`, id).Scan(&code.Resource, &code.TenantId, &scopes)
	if err != nil {
		return CodeInfo{}, err
	}

	if err := json.Unmarshal([]byte(scopes), &code.Scopes); err != nil {
		return CodeInfo{}, err
	}

	return code, nil

}

func (d *DB) InsertTokens(tenant, access, refresh string) error {
	q := `INSERT INTO tokens(tenant_id, access_token, refresh_token) VALUES(?,?,?)`

	_, err := d.conn.Exec(q, tenant, access, refresh)

	return err
}
