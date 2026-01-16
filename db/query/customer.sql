-- name: CustomersList :many
SELECT customer_id FROM customer;

-- name: CustomersListBySearch :many
SELECT customer_id FROM customer
		where
			first_name ilike '%' || @search::text || '%'
		or
 			last_name ilike '%' || @search::text || '%'
		or
 			street ilike '%' || @search::text || '%'
		or
 			town ilike '%' || @search::text || '%'
		or
 			remark ilike '%' || @search::text || '%'
		ORDER BY customer_id DESC
		LIMIT
			@page_size
		OFFSET
 			@page_count ;

-- name: CustomersListCount :one
SELECT count(*) FROM customer
where
	first_name ilike '%' || @search::text || '%'
or
	last_name ilike '%' || @search::text || '%';

-- name: CustomerGetById :one
SELECT * FROM customer
WHERE customer_id = @customerid;

-- name: CustomerDelete :exec
DELETE FROM customer WHERE customer_id = @customerid;

-- name: CustomerCreate :one
INSERT INTO Customer
	(first_Name, last_Name, street, town, remark, email, phone, ahv, alv, quellensteuer, salary_hour)
VALUES
	(@firstname, @lastname, @street, @town, @remark,@email, @phone, @ahv, @alv, @quellensteuer, @salaryhour)
RETURNING Customer_Id;


-- name: CustomerUpdate :exec
UPDATE Customer
SET
	first_Name = @firstname,
	last_Name = @lastname,
	street = @street,
	town = @town,
	remark = @remark,
	email = @email,
	phone = @phone,
	ahv = @ahv,
	alv = @alv,
	quellensteuer = @quellensteuer,
	salary_hour = @salaryhour
WHERE customer_Id = @customerid;

