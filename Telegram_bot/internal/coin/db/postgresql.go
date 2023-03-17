package coin

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"go.mod/internal/coin"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ") //реплейсим все табы на пусьльу
}

func (r *repository) Create(ctx context.Context, coin *coin.Coin) error {

	q :=
		`
		INSERT INTO coin
			(name)
		VALUES 
		    ($1) 
		RETURNING id
		`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q))) // сюда можно положить все запросы

	if err := r.client.QueryRow(ctx, q, coin.Name).Scan(&coin.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr) // вывод ошибок запроса
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) (c []coin.Coin, err error) {
	q :=
		`
		SELECT id, name FROM public.coin;
		`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	coins := make([]coin.Coin, 0)

	for rows.Next() {
		var cin coin.Coin

		err = rows.Scan(&cin.ID, &cin.Name) // распихиваем поля по моделям
		if err != nil {
			return nil, err
		}

		coins = append(coins, cin)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coins, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (coin.Coin, error) {
	q :=
		`
		SELECT id, name FROM public.coin WHERE id = $1
		`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	var cin coin.Coin
	err := r.client.QueryRow(ctx, q, id).Scan(&cin.ID, &cin.Name)
	if err != nil {
		return coin.Coin{}, err
	}
	return cin, nil
}
func (r *repository) Update(ctx context.Context, user coin.Coin) error {
	//TODO implement me
	panic("implement me")
}
func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) coin.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
