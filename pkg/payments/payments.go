package payments

import (
	"errors"
	"github.com/DaniilOr/goMongo/cmd/service/app/dtos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNoToken  = errors.New("no token")
)

type Service struct {
	db *mongo.Database
}

func NewService(db *mongo.Database) *Service {
	return &Service{db: db}
}

func (s *Service) GetPayments(r *http.Request, id int64) ([]dtos.Payment, error) {
	cursor, err := s.db.Collection("suggestions").Find(r.Context(),
		bson.D{{"user_id", id}})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(r.Context()); cerr != nil {
			log.Print(cerr)
		}
	}()
	var result dtos.Response
	for cursor.Next(r.Context()) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 3)
	return result.SuggestedPayments, nil
}

func (s *Service) AddPredictedPayment(r *http.Request, id int64, payment dtos.Payment) error {
	result, err := s.db.Collection("suggestions").UpdateOne(r.Context(), bson.M{"user_id": id}, bson.D{
		{"$push", bson.D{{"suggested_payments", bson.D{{"icon", payment.Icon},
			{"name", payment.Name}, {"link", payment.Link}}}}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
