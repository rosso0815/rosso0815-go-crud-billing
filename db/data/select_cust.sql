--
select * from customer;
--
SELECT *
FROM customer
where
    first_name ilike '%' || '01' || '%'
    or last_name ilike '%' || '01' || '%';

INSERT INTO
    Customer (
        first_Name,
        last_Name,
        street,
        town,
        remark
    )
VALUES ('A', 'A', 'A', 'A', 'A')
RETURNING    Customer_Id;

INSERT INTO
    Customer (
        first_Name,
        last_Name,
        street,
        town,
        remark
    )
VALUES ('A', 'B', 'A', 'A', 'A')
RETURNING    Customer_Id;


update customer set street = 'C' where first_name = 'A' and last_name = 'B';