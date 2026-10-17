package db

import "encoding/json"

type Token struct {
	Id             int    `json:"id"`
	AccessToken    string `json:"access_token"`
	AccessTokenLen int    `json:"access_len"`

	RefreshToken   string `json:"refresh_token"`
	RfreshTokenLen int    `json:"refresh_len"`
	Created_at     string `json:"created_at"`
	TenantId       string `json:"tenant_id"`
}

type TokenRefreshInfo struct {
	Refresh  string
	TenantId string
	Scopes   []string
	Resource string
}

func (d *DB) GetToken(unqiueStr, token_id string) (TokenRefreshInfo, error) {
	q := `SELECT refresh_token, tenant_id FROM tokens WHERE token_id = ? and unqiueStr = ?`

	var t TokenRefreshInfo
	err := d.conn.QueryRow(q, token_id, unqiueStr).Scan(&t.Refresh, &t.TenantId)
	if err != nil {
		return TokenRefreshInfo{}, err
	}

	q2 := `SELECT resource, scopes FROM campaigns WHERE UniqueStr = ?`
	var scopes string
	err = d.conn.QueryRow(q2, unqiueStr).Scan(&t.Resource, &scopes)
	if err != nil {
		return TokenRefreshInfo{}, err
	}
	if err := json.Unmarshal([]byte(scopes), &t.Scopes); err != nil {
		return TokenRefreshInfo{}, err
	}

	return t, nil
}

func (d *DB) ListTokens(id string) ([]Token, error) {
	q := `SELECT access_token, refresh_token, tenant_id, created_at FROM tokens WHERE unqiueStr = ?`
	rows, err := d.conn.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []Token

	for rows.Next() {
		var t Token
		if err := rows.Scan(&t.AccessToken, &t.RefreshToken, &t.TenantId, &t.Created_at); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (d *DB) DeleteAllTokens(id string) error {
	q := `DELETE FROM tokens WHERE unqiueStr = ?`
	_, err := d.conn.Exec(q, id)
	return err
}

func (d *DB) DeleteToken(unqiueId, tokenId string) error {
	q := `DELETE FROM tokens WHERE unqiueStr = ? AND token_id = ?`
	_, err := d.conn.Exec(q, unqiueId, tokenId)
	return err
}
