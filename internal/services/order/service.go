package order

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService struct {
	repo OrderRepository
	sqs  *sqs.SQS
}

func NewOrderService(repo OrderRepository, sqs *sqs.SQS) *OrderService {
	return &OrderService{repo: repo, sqs: sqs}
}

// TODO: requisição para o serviço de pagamento, para validar o token de pagamento e processar o pagamento
// TODO: caso o pagamento seja aprovado e os serviços de voo e hotel estejam disponíveis, criar o pedido no banco de dados com status "pending"
// TODO: publicar uma mensagem na fila do SQS para processar o pedido
func (s *OrderService) Create(ctx context.Context, order Order) error {
	// TODO: Verificar se existe duplicidade de pedido, para evitar que o mesmo pedido seja criado mais de uma vez (user_id ter o mesmo flight_id e hotel_id e com status pending)
	pgUserID := pgtype.UUID{
		Bytes: order.UserID,
		Valid: true,
	}

	orders, err := s.repo.GetOrdersByUserID(ctx, pgUserID)

	if err != nil {
		return fmt.Errorf("get orders by user id failed")
	}

	// TODO: verificação meio ruim, creio que eu consiga melhorar
	for _, o := range orders {
		if o.FlightID == order.FlightID && o.HotelID == order.HotelID && o.Status == Pending {
			return fmt.Errorf("order duplicated")
		}
	}

	// TODO: requisição para os servições de voo e hotel, para validar a disponibilidade.

	return s.repo.Create(ctx, order)
}
