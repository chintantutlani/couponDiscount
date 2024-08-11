1. Creating Coupons

API Endpoint: POST /coupons

This API endpoint allows the creation of a new coupon. The user sends JSON data containing the coupon details.

Assumptions:

The coupon creation process is flexible enough to accommodate any type of coupon (e.g., cart-wise, product-wise).
The discount_type can be either "flat" or "percentage" to specify the type of discount.
Edge Cases:

Duplicate coupon codes and names can be created.
Coupons with invalid or missing fields.
Limitations:

Consider adding uniqueness constraints to the coupon name and code to prevent duplicates.
Some fields are included for future use, potentially making coupons more unique for specific users.
A subtype category can be added to certain coupon types. For example, for a cart-wise coupon, a subtype could offer an additional 5% discount for online payments.


2. Retrieving All Coupons

API Endpoint: GET /coupons

This API endpoint retrieves a list of all available coupons in the system. It returns an array of coupon objects.

Steps:

Call the service layer to fetch all coupons from the MongoDB collection.
If an error occurs (e.g., database connection issues), return an error response.
If no coupons are found, return an appropriate error message.
If coupons are found, return the list of coupons in the response.

Assumptions:

The coupon model is designed to handle various coupon types and can be retrieved without additional processing.
The service layer correctly handles empty or null results and returns a proper error message if no documents are found.

Edge Cases:

Database connection failures.
No coupons are found in the collection.
Inconsistent data, such as incomplete coupon records in the database.

Limitations-

The current implementation assumes that all coupons in the database are valid and properly formatted. Consider adding data validation checks when retrieving the coupons.
Consider adding pagination to handle a large number of coupons efficiently.

3. Retrieving a Coupon by ID

API Endpoint: GET /coupons/:id

This API endpoint retrieves a specific coupon based on its unique ID.

Steps:

Extract the id parameter from the request URL.
Convert the string id to a primitive.ObjectID.
Call the service layer to retrieve the coupon associated with this ObjectID from the MongoDB collection.
If the coupon is found, return it in the response.

Assumptions:

The id parameter passed in the URL is a valid MongoDB ObjectID.
Each coupon has a unique ObjectID in the database.
Limitations and Improvements:

The system currently assumes the id is always a valid ObjectID. Consider adding more detailed error handling for invalid ID formats.
Consider implementing a caching layer to reduce database load for frequently accessed coupons.

4. Updating a Coupon by ID

API Endpoint: PUT /coupons/:id

This API endpoint allows updating specific fields of an existing coupon. The user provides the fields to update in a JSON payload.

Steps:

Extract the id parameter from the request URL.
Bind the incoming JSON payload to a map of update fields.
Validate the update fields against a predefined list of allowed fields.
Convert the id parameter to a primitive.ObjectID.
Call the service layer to update the coupon in the MongoDB collection using the provided ObjectID and update fields.
If the update is successful, return a success message.

Assumptions:

The fields being updated are part of a predefined list of allowed fields.
The id parameter passed in the URL is a valid MongoDB ObjectID.

Limitations:

The system currently restricts updates to a specific list of fields. Consider making this list configurable to allow more flexibility in the future.
If the id is invalid or if no matching coupon is found, the system returns a generic error. Consider providing more specific error messages.

5. Deleting a Coupon by ID

API Endpoint: DELETE /coupons/:id

This API endpoint deletes a coupon from the database based on the provided coupon ID.

Steps:

Extract the id parameter from the request URL.
Convert the extracted id to a MongoDB primitive.ObjectID.
If the conversion fails, return an error response.
Delete the coupon document from the MongoDB collection using the provided ObjectID.
If the coupon is successfully deleted, return a success message.

Edge Cases:

The user does not have permission to delete the coupon.
Two requests to delete the same coupon are made simultaneously.

Limitations:

This implementation assumes hard deletes.
Consider implementing soft deletes, where the coupon is marked as inactive or deleted, but the record remains in the database for auditing or recovery.
The current implementation does not account for simultaneous operations on the same coupon.
Implement logging to record when and by whom the coupon was deleted.

6. Fetching Applicable Coupons

API Endpoint: POST /applicable-coupons

This API endpoint fetches all applicable coupons for a given cart and calculates the total discount that each coupon will apply.

Steps:

Validate the incoming JSON payload representing the cart.
Calculate the total cart price using the calculateCartPrice function.
Iterate over each coupon and determine its applicability:
Cart-wise: Compare the total cart price with the coupon's threshold value. Calculate the discount based on the discount type.
Product-wise: Check if any product in the cart matches the product_id in the coupon details. Apply the discount accordingly.
BxGy: Check if the cart satisfies the buy conditions. If so, apply the corresponding get products as a discount.

Assumptions:

The cart structure is well-defined and contains all necessary fields for the calculations.
The details field in the coupon struct is a map used to store relevant fields. For example:
For a product-wise coupon, product_id should be saved in details.
For a BxGy coupon, the details field should store a buy_product array containing product_id and quantity, along with a get_product array with product details.

Edge Cases:

Empty cart or no products in the cart that match any coupon.
Multiple applicable coupons; determining which to apply might require prioritization.
Invalid or unexpected data in the details field (e.g., missing product IDs, incorrect data types).

Limitations:

The system does not prioritize between multiple applicable coupons, so it may not apply the best offer for the user.
The system only supports BxGy type coupons and does not handle BxGy with Z% off or BxGy y% off.
Consider adding more robust error handling for unexpected data types in the details field.
Consider adding a mechanism to limit the number of times a coupon can be applied to avoid abuse.
The system currently lacks an expiry date for coupons, and only applicable coupons are shown, not the best applicable coupons.
For product-wise coupons, discounts can only be applied to products, not to brands, specific categories, bundles, or large-size discounts.

7. Applying a Coupon to a Cart

API Endpoint: POST /coupons/:id

This API endpoint applies a specific coupon to a user's cart and returns the updated cart with the applied discount.

Steps:

Validate the incoming JSON payload representing the cart.
Retrieve the coupon using the provided coupon ID.
Determine if the coupon applies to the cart:
Product-wise: Apply the discount to the matching product(s) in the cart.
BxGy: Check if the buy conditions are met, and if so, apply the corresponding get products as a discount.
Cart-wise: Calculate the discount based on the total cart price and threshold value.
Update the cart by applying the calculated discounts.
Return the updated cart with the total discount and final price.

Assumptions:

The coupon ID provided is valid and corresponds to an existing coupon in the database.
The cart structure is valid and contains all necessary fields for calculations.

Edge Cases:

Invalid coupon ID (coupon not found).
Coupons are no longer valid or have been used up (e.g., max uses limit reached).
Applying a coupon that doesn't match the cart (e.g., product not in the cart, cart total below threshold value).

Limitations:

Implement expiration dates and check if the coupon is still valid before applying it.
Add the ability to apply multiple discounts or combine coupons, if allowed.
Improve the logic to handle scenarios where multiple discounts apply to the same cart, ensuring that the customer receives the best possible deal.
Allow a single coupon to be used multiple times, so it doesn't expire after just one use.
Track how many times different users have used a single coupon to enforce usage limits.
Implement checks for expired coupons to prevent adding them to the cart.
Specify certain coupons for particular users, preventing others from using the coupon code even if they have it.
Define which types of coupons can be added to a cart, preventing all coupon types from being applied to any order.
If an order is canceled or refunded, the coupon should be reinstated or removed from the system.