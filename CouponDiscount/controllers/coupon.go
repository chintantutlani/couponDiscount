package coupon

import (
	"log"
	"monk_commerce/models"
	services "monk_commerce/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CouponController struct {
	couponService *services.CouponServices
}

func NewCouponController(couponService *services.CouponServices) *CouponController {
	return &CouponController{
		couponService: couponService,
	}
}

func (cp *CouponController) CreateNewCoupon(ctx *gin.Context) {

	var coupon models.Coupon
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := cp.couponService.CreateCoupon(coupon)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": result.InsertedID})

}

func (cp *CouponController) GetAll(ctx *gin.Context) {
	users, err := cp.couponService.GetAll()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (cp *CouponController) GetCouponById(ctx *gin.Context) {
	var couponid string = ctx.Param("id")
	coupon, err := cp.couponService.GetCoupon(&couponid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, coupon)
}

func (cp *CouponController) UpdateCouponByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateFields map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateFields); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if len(updateFields) == 0 {
		log.Println("No fields to update")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "No fields provided for update"})
		return
	}

	allowedFields := map[string]bool{
		"name":                 true,
		"code":                 false,
		"maxuses":              false,
		"type":                 false,
		"details":              false,
		"discount":             false,
		"description":          true,
		"condition":            true,
		"discount_type":        false,
		"threshold_value":      false,
		"expiry_date":          false,
		"free_shipping":        false,
		"use_limit":            false,
		"usage_limit_per_user": false,
		"maximum_amount":       false,
	}

	for key := range updateFields {
		if !allowedFields[key] {
			log.Printf("Field '%s' is not allowed for update", key)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Field '" + key + "' is not allowed for update"})
			return
		}
	}

	err := cp.couponService.UpdateCouponByID(id, updateFields)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Coupon updated successfully"})
}

func (cp *CouponController) DeleteCouponByID(ctx *gin.Context) {
	var couponid string = ctx.Param("id")
	err := cp.couponService.DeleteCouponByID(&couponid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Coupon deleted successfully"})
}

func (cp *CouponController) GetApplicableCoupons(ctx *gin.Context) {
	var cart models.Cart
	if err := ctx.ShouldBindJSON(&cart); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	applicableCoupons, err := cp.couponService.GetApplicableCoupons(cart)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if applicableCoupons == nil {
		ctx.JSON(http.StatusOK, gin.H{"applicable_coupons": "Coupon not available for this cart"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"applicable_coupons": applicableCoupons})
}

func (cc *CouponController) ApplyCouponByID(ctx *gin.Context) {
	var request struct {
		Cart models.Cart `json:"cart"`
	}
	couponID := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	coupon, err := cc.couponService.GetCoupon(&couponID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Coupon not found"})
		return
	}

	updatedCart, err := cc.couponService.ApplyAllCoupons(request.Cart, []models.Coupon{*coupon})
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"updated_cart": updatedCart})
}

func (c CouponController) RegisterCouponRoutes(rg *gin.RouterGroup) {

	rg.POST("/coupons", c.CreateNewCoupon)
	rg.GET("/coupons", c.GetAll)
	rg.GET("/coupons/:id", c.GetCouponById)
	rg.PUT("/coupons/:id", c.UpdateCouponByID)
	rg.DELETE("/coupons/:id", c.DeleteCouponByID)
	rg.POST("/applicable-coupons", c.GetApplicableCoupons)
	rg.POST("/apply-coupon/:id", c.ApplyCouponByID)

}
