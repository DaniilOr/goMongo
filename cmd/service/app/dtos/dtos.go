package dtos
type Payment struct{
	Icon string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Link string `json:"link" bson:"link"`
}
type User struct{
	Id int64 `json:"id" bson:"user_id"`
	Predicted []Payment `json:"predicted_payments" bson:"predicted_payments"`
	Frequent []Payment `json:"Frequent_payments" bson:"frequent_payments"`
}

type TokenDTO struct {
	Token *string `json:"token"`
}
type ResultDTO struct {
	Result string `json:"result"`
	Description string `json:"description,omitempty"`

}
type Response struct{
	Frequent_paymenys []Payment`json:"frequent_payments"`
	Predicted_paymenys []Payment`json:"predicted_payments"`
}