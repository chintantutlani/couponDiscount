package coupon

import (
	"context"
	"errors"
	"fmt"
	"monk_commerce/models"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CouponServices struct {
	couponcollection *mongo.Collection
	ctx              context.Context
}

func NewCouponService(couponcollection *mongo.Collection, ctx context.Context) *CouponServices {
	return &CouponServices{
		couponcollection: couponcollection,
		ctx:              ctx,
	}
}

func (cp *CouponServices) CreateCoupon(coupon models.Coupon) (*mongo.InsertOneResult, error) {
	return cp.couponcollection.InsertOne(cp.ctx, coupon)
}

func (cp *CouponServices) GetAll() ([]*models.Coupon, error) {

	cursor, err := cp.couponcollection.Find(cp.ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(cp.ctx)
	var coupons []*models.Coupon

	if err := cursor.All(cp.ctx, &coupons); err != nil {
		return nil, err
	}

	if len(coupons) == 0 {
		return nil, errors.New("documents not found")
	}

	return coupons, nil
}

func (cp *CouponServices) GetCoupon(id *string) (*models.Coupon, error) {
	var coupon *models.Coupon
	objectID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}

	query := bson.D{bson.E{Key: "_id", Value: objectID}}
	err = cp.couponcollection.FindOne(cp.ctx, query).Decode(&coupon)
	if err != nil {
		return nil, err
	}

	return coupon, err
}

func (cp *CouponServices) UpdateCouponByID(id string, updateFields map[string]interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objectID}}
	update := bson.D{primitive.E{Key: "$set", Value: updateFields}}

	result, err := cp.couponcollection.UpdateOne(cp.ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (cs *CouponServices) GetApplicableCoupons(cart models.Cart) ([]models.ApplicableCoupon, error) {

	var applicableCoupons []models.ApplicableCoupon
	totalCartPrice := calculateCartPrice(cart.Items)

	coupons, err := cs.GetAll()
	if err != nil {
		return applicableCoupons, err
	}

	for _, coupon := range coupons {
		switch coupon.Type {
		case "cart-wise":

			if coupon.Maxuses > 0 {
				if totalCartPrice >= coupon.ThresholdValue {
					var discountAmount float64

					if coupon.DiscountType == "percentage" {
						discountAmount = totalCartPrice * (coupon.Discount / 100)
					} else if coupon.DiscountType == "flat" {
						discountAmount = coupon.Discount
					}

					if discountAmount > 0 {

						if coupon.CuponCode == "" {
							if coupon.Name == "" {
								coupon.Name = "DIS"
							}
							coupon.CuponCode = GenerateCouponCode(coupon.Name, coupon.Discount)
						}

						applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
							CouponID:   coupon.Id.Hex(),
							Type:       "cart-wise",
							Discount:   discountAmount,
							CouponCode: coupon.CuponCode,
						})
					}
				}
			}
		case "product-wise":

			for _, item := range cart.Items {
				if productID, ok := coupon.Details["product_id"].(string); ok && item.ProductID == productID {
					var discountAmount float64

					// Calculate discount based on DiscountType (percentage or flat)
					if coupon.DiscountType == "percentage" {
						discountAmount = item.Price * (coupon.Discount / 100) * float64(item.Quantity)
					} else if coupon.DiscountType == "flat" {
						discountAmount = coupon.Discount * float64(item.Quantity)
					}

					if discountAmount > 0 {

						// Generate Coupon Code if not available
						if coupon.CuponCode == "" {
							if coupon.Name == "" {
								coupon.Name = "DIS"
							}
							coupon.CuponCode = GenerateCouponCode(coupon.Name, coupon.Discount)
						}

						applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
							CouponID:   coupon.Id.Hex(),
							Type:       "product-wise",
							Discount:   discountAmount,
							CouponCode: coupon.CuponCode,
						})
					}
				}
			}

		case "bxgy":
			buyProducts, ok := coupon.Details["buy_products"].([]interface{})
			if !ok {
				if arr, ok := coupon.Details["buy_products"].(primitive.A); ok {
					buyProducts = []interface{}(arr)
				} else {
					return nil, fmt.Errorf("invalid type for buy_products")
				}
			}

			getProducts, ok := coupon.Details["get_products"].([]interface{})
			if !ok {
				if arr, ok := coupon.Details["get_products"].(primitive.A); ok {
					getProducts = []interface{}(arr)
				} else {
					return nil, fmt.Errorf("invalid type for get_products")
				}
			}

			repetitionLimit, ok := coupon.Details["repetition_limit"].(float64)
			if !ok {

				return nil, fmt.Errorf("invalid type for repetition_limit")
			}

			buyCount := make(map[string]int)
			for _, buyProduct := range buyProducts {
				buyProductMap := buyProduct.(map[string]interface{})
				productID := buyProductMap["product_id"].(string)
				quantity := int(buyProductMap["quantity"].(float64))
				buyCount[productID] = quantity
			}

			repetitions := GetMinRepetitions(buyCount, cart)
			fmt.Printf("Repetitions calculated: %v\n", repetitions)

			if repetitions > 0 {
				if int(repetitions) > int(repetitionLimit) {
					repetitions = int(repetitionLimit)
				}

				discount := 0.0
				for _, getProduct := range getProducts {
					getProductMap := getProduct.(map[string]interface{})
					productID := getProductMap["product_id"].(string)
					quantity := int(getProductMap["quantity"].(float64)) * repetitions

					for _, item := range cart.Items {
						if item.ProductID == productID {
							discount += item.Price * float64(quantity)
							fmt.Printf("ProductID: %s, Quantity: %d, Item Price: %f, Total Discount for this product: %f\n", productID, quantity, item.Price, item.Price*float64(quantity))
						}
					}
				}

				if coupon.CuponCode == "" {
					if coupon.Name == "" {
						coupon.Name = "DIS"
					}
					coupon.CuponCode = GenerateCouponCode(coupon.Name, coupon.Discount)
				}

				applicableCoupons = append(applicableCoupons, models.ApplicableCoupon{
					CouponID:   coupon.Id.Hex(),
					Type:       "bxgy",
					Discount:   discount,
					CouponCode: coupon.CuponCode,
				})
			}
		}
	}
	return applicableCoupons, nil
}

func GenerateCouponCode(name string, discount float64) string {

	couponname := strings.ToUpper(name)
	if len(couponname) > 3 {
		couponname = couponname[:3]
	}
	discountname := fmt.Sprintf("%02.0f", discount)
	couponcode := fmt.Sprintf("%s%s", couponname, discountname)
	return couponcode
}

// coupon_service.go

// coupon_service.go
func (cs *CouponServices) ApplyAllCoupons(cart models.Cart, coupons []models.Coupon) (models.UpdatedCart, error) {
	var updatedCart models.UpdatedCart
	updatedCart.Items = cart.Items
	totalDiscount := 0.0

	// Apply Product-wise Discounts
	for _, coupon := range coupons {
		if coupon.Type == "product-wise" {
			for i, item := range updatedCart.Items {
				if productID, ok := coupon.Details["product_id"].(string); ok && item.ProductID == productID {
					itemDiscount := item.Price * (coupon.Details["discount"].(float64) / 100) * float64(item.Quantity)
					updatedCart.Items[i].TotalDiscount += itemDiscount
					totalDiscount += itemDiscount
				}
			}
		}
	}

	// Apply BXGY Discounts
	for _, coupon := range coupons {
		if coupon.Type == "bxgy" {
			buyProducts, ok := coupon.Details["buy_products"].([]interface{})
			if !ok {
				if arr, ok := coupon.Details["buy_products"].(primitive.A); ok {
					buyProducts = []interface{}(arr)
				} else {
					return updatedCart, fmt.Errorf("invalid type for buy_products")
				}
			}

			getProducts, ok := coupon.Details["get_products"].([]interface{})
			if !ok {
				if arr, ok := coupon.Details["get_products"].(primitive.A); ok {
					getProducts = []interface{}(arr)
				} else {
					return updatedCart, fmt.Errorf("invalid type for get_products")
				}
			}

			repetitionLimit, ok := coupon.Details["repetition_limit"].(float64)
			if !ok {
				return updatedCart, fmt.Errorf("invalid type for repetition_limit")
			}

			// Calculate Repetitions
			buyCount := make(map[string]int)
			for _, buyProduct := range buyProducts {
				buyProductMap := buyProduct.(map[string]interface{})
				productID := buyProductMap["product_id"].(string)
				quantity := int(buyProductMap["quantity"].(float64))
				buyCount[productID] = quantity
			}

			repetitions := GetMinimunRepetitions(buyCount, updatedCart.Items)
			if repetitions > int(repetitionLimit) {
				repetitions = int(repetitionLimit)
			}

			// Apply Get Products Discounts
			for _, getProduct := range getProducts {
				getProductMap := getProduct.(map[string]interface{})
				productID := getProductMap["product_id"].(string)
				quantity := int(getProductMap["quantity"].(float64)) * repetitions

				for i, item := range updatedCart.Items {
					if item.ProductID == productID {
						updatedCart.Items[i].Quantity += quantity
						itemDiscount := item.Price * float64(quantity)
						updatedCart.Items[i].TotalDiscount += itemDiscount
						totalDiscount += itemDiscount
					}
				}
			}
		}
	}

	// Apply Cart-wise Discounts
	for _, coupon := range coupons {
		if coupon.Type == "cart-wise" {
			totalCartPrice := calculateCartPrice(updatedCart.Items)
			if totalCartPrice >= coupon.ThresholdValue {
				cartDiscount := totalCartPrice * (coupon.Discount / 100)
				totalDiscount += cartDiscount
			}
		}
	}

	// Calculate final price and total discount
	updatedCart.TotalPrice = calculateCartPrice(updatedCart.Items)
	updatedCart.TotalDiscount = totalDiscount
	updatedCart.FinalPrice = updatedCart.TotalPrice - totalDiscount

	return updatedCart, nil
}

// Helper function to calculate the total price of items in the cart
// func calculateCartPrice(items []models.CartItem) float64 {
// 	total := 0.0
// 	for _, item := range items {
// 		total += float64(item.Quantity) * item.Price
// 	}
// 	return total
// }

// Helper function to calculate the minimum number of repetitions for BXGY discounts
func GetMinimunRepetitions(buyCount map[string]int, items []models.CartItem) int {
	minReps := int(^uint(0) >> 1) // Set to maximum integer value
	for productID, requiredQty := range buyCount {
		for _, item := range items {
			if item.ProductID == productID {
				availableReps := item.Quantity / requiredQty
				if availableReps < minReps {
					minReps = availableReps
				}
			}
		}
	}
	return minReps
}

func calculateCartPrice(items []models.CartItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func GetMinRepetitions(buyCount map[string]int, cart models.Cart) int {
	totalquantity := 0
	for productID, requiredquantity := range buyCount {
		availableQuantity := 0
		for _, item := range cart.Items {
			if item.ProductID == productID {
				availableQuantity += item.Quantity
			}
		}
		if availableQuantity < requiredquantity {
			return 0
		}
		totalquantity += availableQuantity / requiredquantity
	}
	return totalquantity
}
