package local

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tiptok/gocomm/pkg/broker/models"
)

type PgMessageReceiverRepository struct {
	Db orm.DB
	Tx orm.DB
}

func NewPgMessageReceiverRepository(db *pg.DB, tx *pg.Tx) *PgMessageReceiverRepository {
	return &PgMessageReceiverRepository{
		Db: db,
		Tx: tx,
	}
}

func (repository *PgMessageReceiverRepository) ReceiveMessage(params map[string]interface{}) error {
	type queryType struct {
		Num    int
		Status int
	}
	var query queryType
	checkSql := `select count(0) num,sum(status) status from sys_message_consume where "offset" =? and topic=? limit 1`
	_, err := repository.Db.Query(&query, checkSql, params["offset"], params["topic"])
	if err != nil {
		return err
	}
	if query.Num == 0 {
		sql := `insert into sys_message_consume(topic,partition,"offset",key,value,msg_time,create_at,status)values(?,?,?,?,?,?,?,?)`
		_, err = repository.Db.Exec(sql, params["topic"], params["partition"], params["offset"], params["key"], params["value"], params["msg_time"], params["create_at"], params["status"])
	}
	if query.Num > 0 && query.Status == int(models.Finished) {
		return fmt.Errorf("receive repeate message status:%v [%v]", query.Status, params)
	}
	return err
}

func (repository *PgMessageReceiverRepository) ConfirmReceive(params map[string]interface{}) error {
	_, err := repository.Db.Exec(`update sys_message_consume set status=? where "offset" =? and topic=?`, int(models.Finished), params["offset"], params["topic"])
	return err
}
