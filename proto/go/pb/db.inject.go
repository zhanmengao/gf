package pb

import (
	"fmt"
)

func NewDBConnect(UID string) *DBConnect {
	ret := &DBConnect{
		UID: UID,
	}
	return ret
}

func (p *DBConnect) GetDBKey() string {
	key := fmt.Sprintf("DBConn:%s", p.UID)
	return key
}

func (p *DBConnect) GetDBKeyFormat() string {
	return "DBConn:%s"
}

func (p *DBConnect) GetDBField() string {
	key := fmt.Sprintf("")
	return key
}

func (p *DBConnect) GetDBFieldFormat() string {
	return ""
}

func (p *DBConnect) GetDBKeyName() string {
	return "DBConnect"
}

func (p *DBConnect) GetDBName() string {
	return "Redis"
}

func NewDBSessionKey(SessionKey string) *DBSessionKey {
	ret := &DBSessionKey{
		SessionKey: SessionKey,
	}
	return ret
}

func (p *DBSessionKey) GetDBKey() string {
	key := fmt.Sprintf("DBSessionKey:%s", p.SessionKey)
	return key
}

func (p *DBSessionKey) GetDBKeyFormat() string {
	return "DBSessionKey:%s"
}

func (p *DBSessionKey) GetDBField() string {
	key := fmt.Sprintf("")
	return key
}

func (p *DBSessionKey) GetDBFieldFormat() string {
	return ""
}

func (p *DBSessionKey) GetDBKeyName() string {
	return "DBSessionKey"
}

func (p *DBSessionKey) GetDBName() string {
	return "Redis"
}
