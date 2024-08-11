
---

### 1. Create Coupons
**API Endpoint:** `POST /coupons`

This API endpoint allows the creation of a new coupon. The user sends JSON data containing coupon details.

**Assumptions:**
- The create coupon functionality is flexible enough to accommodate any type of coupon (e.g., cart-wise, product-wise coupons). The `discount_type` field is used to specify whether the discount is a “flat” or “percentage” type.

**Edge Cases:**
- Duplicate coupon codes and names can be created.
- Coupons may have invalid or missing fields.

**Limitations:**
- Consider adding uniqueness constraints to the name and coupon code to prevent duplicates.
- Added some fields that are currently unused but could be useful for future development, making coupons more unique for particular users.
- A subtype category could be added for specific coupon types. For example, for cart-wise coupons, a subtype could offer an additional 5% discount for online payments.

### 2. Applicable Coupons
**API Endpoint:** `POST /applicable-coupons`

This endpoint fetches all applicable coupons for a given cart and calculates the total discount each coupon will apply.

**Steps:**
- Validate the incoming JSON payload representing the cart.
- Calculate the total cart price using the `calculateCartPrice` function.
- Iterate over each coupon and determine its applicability:
  - **Cart-wise:** Compare the total cart price with the coupon's threshold value and calculate the discount based on the discount type.
  - **Product-wise:** Check if any product in the cart matches the `product_id` in the coupon details, and apply the discount accordingly.
  - **BxGy:** Check if the cart satisfies the buy conditions. If so, apply the corresponding get-products as a discount.

**Assumptions:**
- The cart structure is well-defined and contains all necessary fields for the calculations.
- The `details` field in the `coupons` struct is a map used to store relevant fields. For example:
  - For the product-wise type, `product_id` is stored in `details`.
  - For the BxGy type, the `buy_product` array containing `product_id` and quantity is stored, along with `get_product`, `product_id`, `quantity`, and the last reputation.

**Edge Cases:**
- An empty cart or no products in the cart match any coupon.
- Multiple applicable coupons may need prioritization.
- Invalid or unexpected data in the `details` field (e.g., missing product IDs, incorrect data types).

**Limitations:**
- The system does not prioritize between multiple applicable coupons, so it is unclear which would be the best offer for the user.
- The system only works for BxGy and not for BxGy with a Z% off or BxG y% off.
- More robust error handling should be added for unexpected data types in the `details` field.
- Consider adding a mechanism to limit the number of times a coupon can be applied to avoid abuse.
- There is no expiry date for the coupons, and only applicable coupons will be shown, not the best applicable coupons.
- For product-wise coupons, discounts can only be applied to products, not to brands, particular categories, bundles, or large-size discounts.

### 3. Applying a Coupon to a Cart
**API Endpoint:** `POST /coupons/:id`

This API endpoint applies specific coupons to user carts and returns the updated cart with the applied discount.

**Steps:**
- Validate the incoming JSON payload representing the cart.
- Retrieve the coupon using the provided coupon ID.
- Determine if the coupon applies to the cart:
  - **Product-wise:** Apply the discount to the matching product(s) in the cart.
  - **BxGy:** Check if the buy conditions are met, and if so, apply the corresponding get-products as a discount.
  - **Cart-wise:** Calculate the discount based on the total cart price and threshold value.
- Update the cart by applying the calculated discounts.
- Return the updated cart with the total discount and final price.

**Assumptions:**
- The coupon ID provided is valid and corresponds to an existing coupon in the database.
- The cart structure is valid and contains all necessary fields for calculations.

**Edge Cases:**
- Invalid coupon ID (coupon not found).
- Coupons are no longer valid or have been used up (e.g., max usage limit reached).
- Applying a coupon that doesn't match the cart (e.g., product not in the cart, cart total below the threshold value).

**Limitations:**
- Implement expiration dates and check if the coupon is still valid before applying it.
- Add the ability to apply multiple discounts or combine coupons if allowed.
- Improve the logic to handle scenarios where multiple discounts apply to the same cart, ensuring that the customer receives the best possible deal.
- Allow a single coupon to be used multiple times, so it doesn't expire after just one use.
- Track how many times different users have used a single coupon to enforce a usage limit.
- Prevent expired coupons from being added to the cart.
- Specify a particular coupon for a specific user, ensuring that no one else can use that coupon even if they have the coupon code.
- Restrict certain types of coupons from being applied to any type of cart order.
- If the order is canceled or refunded, the coupon should be reinstated or removed.