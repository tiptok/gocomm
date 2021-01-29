package local

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/broker/models"
)

type PgMessageProduceRepository struct {
	Db orm.DB
	Tx orm.DB
}

func (repository *PgMessageProduceRepository) SaveMessage(message *models.Message) error {
	sql := `insert into sys_message_produce (id,topic,value,msg_time,status)values(?,?,?,?,?)`
	_, err := repository.DB().Exec(sql, message.Id, message.Topic, common.JsonAssertString(message), message.MsgTime, int64(models.UnFinished))
	return err
}
func (repository *PgMessageProduceRepository) FindNoPublishedStoredMessages() ([]*models.Message, error) {

	sql := `select value from sys_message_produce where status=?`
	var values []string
	_, e := repository.Db.Query(&values, sql, int64(models.UnFinished))
	var messages = make([]*models.Message, 0)
	if e != nil {
		return messages, nil
	}
	for _, v := range values {
		item := &models.Message{}
		common.JsonUnmarshal(v, item)
		if item.Id != 0 {
			messages = append(messages, item)
		}
	}
	return messages, nil
}
func (repository *PgMessageProduceRepository) FinishMessagesStatus(messageIds []int64, finishStatus int) error {
	_, err := repository.DB().Exec("update sys_message_produce set status=? where id in (?)", finishStatus, pg.In(messageIds))
	return err
}

func (repository *PgMessageProduceRepository) DB() orm.DB {
	if repository.Tx != nil {
		return repository.Tx
	}
	return repository.Db
}
func NewPgMessageProduceRepository(db orm.DB, tx orm.DB) *PgMessageProduceRepository {
	return &PgMessageProduceRepository{
		Db: db,
		Tx: tx,
	}
}
