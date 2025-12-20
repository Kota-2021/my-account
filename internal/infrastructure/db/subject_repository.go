package db

import (
	"context"
	"my-account/internal/domain"
	"my-account/internal/infrastructure/db/dbgen"

	"github.com/jackc/pgx/v5"
)

// SaveSubjects は勘定科目マスタを一括登録・更新します
func SaveSubjects(ctx context.Context, tx pgx.Tx, subjects []domain.Subject) error {
	q := dbgen.New(tx)
	for _, s := range subjects {
		// subject_codeはsmallintなのでint16にキャスト
		err := q.SaveSubject(ctx, dbgen.SaveSubjectParams{
			SubjectCode: int16(s.Code),
			SubjectName: s.Name,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func FetchAllSubjects(ctx context.Context, conn *pgx.Conn) ([]domain.Subject, error) {
	q := dbgen.New(conn)
	dbSubjects, err := q.ListSubjects(ctx)
	if err != nil {
		return nil, err
	}

	var results []domain.Subject
	for _, s := range dbSubjects {
		results = append(results, domain.Subject{
			Code: int16(s.SubjectCode),
			Name: s.SubjectName,
		})
	}
	return results, nil
}
