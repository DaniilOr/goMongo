package payments

import (
	"errors"
	"github.com/DaniilOr/goMongo/cmd/service/app/dtos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var (
	ErrNotFound = errors.New("not found")
	ErrNoToken  = errors.New("no token")
)

type Service struct {
	db *mongo.Database
}
type Answer struct {
	User_id           int64
	FrequentPayments  []dtos.Payment `bson:"frequent_payments"`
	PredictedPayments []dtos.Payment `bson:"predicted_payments"`
}

func NewService(db *mongo.Database) *Service {
	return &Service{db: db}
}

func (s *Service) GetPayments(r *http.Request, id int64) ([]dtos.Payment, []dtos.Payment, error) {
	cursor, err := s.db.Collection("payments").Find(r.Context(),
		bson.D{{"user_id", id}})
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cerr := cursor.Close(r.Context()); cerr != nil {
			log.Print(cerr)
		}
	}()
	var frequentPayments []dtos.Payment
	var predictedPayments []dtos.Payment
	for cursor.Next(r.Context()) {
		var result Answer
		err = cursor.Decode(&result)
		if err != nil {
			return nil, nil, err
		}
		frequentPayments = result.FrequentPayments
		predictedPayments = result.PredictedPayments
	}
	if err = cursor.Err(); err != nil {
		return nil, nil, err
	}
	return frequentPayments, predictedPayments, nil
}

func (s *Service) AddPredictedPayment(r *http.Request, id int64, payment dtos.Payment) error {
	result, err := s.db.Collection("payments").UpdateOne(r.Context(), bson.M{"user_id": id}, bson.D{
		{"$push", bson.D{{"predicted_payments", bson.D{{"icon", payment.Icon},
			{"name", payment.Name}, {"link", payment.Link}}}}}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
