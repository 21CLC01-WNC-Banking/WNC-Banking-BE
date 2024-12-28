package repositoryimplement

import (
	"context"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/jmoiron/sqlx"
)

type DebtReplyRepository struct {
	db *sqlx.DB
}

func NewDebtReplyRepository(db database.Db) repository.DebtReplyRepository {
	return &DebtReplyRepository{db: db}
}

func (repo *DebtReplyRepository) CreateCommand(ctx context.Context, reply *entity.DebtReply) error {
	insertQuery := `INSERT INTO debt_reply(debt_reminder_id, user_reply_name, content) VALUES (:debt_reminder_id,:user_reply_name, :content)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, reply)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DebtReplyRepository) GetReplyByDebtIdQuery(ctx context.Context, debtId string) (*entity.DebtReply, error) {
	var reply entity.DebtReply
	query := `SELECT * FROM debt_reply WHERE debt_reminder_id = ?`
	err := repo.db.QueryRowxContext(ctx, query, debtId).StructScan(&reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
