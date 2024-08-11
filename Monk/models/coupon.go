package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Coupon struct {
	Id                primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string                 `json:"name" bson:"name"`
	CuponCode         string                 `json:"code" bson:"code"`
	Maxuses           int                    `json:"maxuses" bson:"maxuses"`
	Type              string                 `json:"type" bson:"type"`
	Details           map[string]interface{} `json:"details" bson:"details"`
	Discount          float64                `json:"discount" bson:"discount"`
	Description       string                 `json:"description" bson:"description"`
	Condition         string                 `json:"condition" bson:"condition"`
	DiscountType      string                 `json:"discount_type" bson:"discount_type"`
	ThresholdValue    float64                `json:"threshold_value" bson:"threshold_value"`
	UserId            string                 `json:"user_id" bson:"user_id"`                           // unused
	ExpiryDate        string                 `json:"expiry_date" bson:"expiry_date"`                   // unused
	FreeShipping      bool                   `json:"free_shipping" bson:"free_shipping"`               // unused
	Uselimity         bool                   `json:"use_limit" bson:"use_limit"`                       // unused
	UsageLimitPerUser int                    `json:"usage_limit_per_user" bson:"usage_limit_per_user"` // unused
	MaximunAmount     float64                `json:"maximun_amount" bson:"maximum_amount"`             // unused
}

type Cart struct {
	Items []CartItem `json:"items" bson:"items"`
}

type CartItem struct {
	ProductID     string  `json:"product_id" bson:"product_id"`
	Quantity      int     `json:"quantity" bson:"quantity"`
	Price         float64 `json:"price" bson:"price"`
	TotalDiscount float64 `json:"total_discount" bson:"total_discount"`
}

type ApplicableCoupon struct {
	CouponID   string  `json:"coupon_id" bson:"coupon_id"`
	Type       string  `json:"type" bson:"type"`
	Discount   float64 `json:"discount" bson:"discount"`
	CouponCode string  `json:"coupon_code" bson:"coupon_code"`
}

type UpdatedCart struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"total_price"`
	TotalDiscount float64    `json:"total_discount"`
	FinalPrice    float64    `json:"final_price"`
}
