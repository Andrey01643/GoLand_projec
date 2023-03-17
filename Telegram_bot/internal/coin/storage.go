package coin

import "context"

type Repository interface {
	Create(ctx context.Context, coin *Coin) error
	FindAll(ctx context.Context) (c []Coin, err error)
	FindOne(ctx context.Context, id string) (Coin, error)
	Update(ctx context.Context, user Coin) error
	Delete(ctx context.Context, id string) error
}
