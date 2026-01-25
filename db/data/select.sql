--

select * from customer;

select * from invoice;

select * from invoiceentry;

-- https://stackoverflow.com/questions/55950437/how-to-join-my-invoice-and-invoice-product-table-with-sql-query

-- https://www.w3schools.com/sql/sql_groupby.asp
-- SELECT Shippers.ShipperName, COUNT(Orders.OrderID) AS NumberOfOrders FROM Orders
-- LEFT JOIN Shippers ON Orders.ShipperID = Shippers.ShipperID
-- GROUP BY ShipperName;

select
    COALESCE(c.customer_id,0) as CustomerID,
    COALESCE(c.first_name,'null') as FirstName,
    c.last_name,
    i.invoice_id,
    i.invoice_year,
    i.invoice_month,
    COALESCE(SUM(ie.work_hours),0) as MONTH_SUM
from
    invoice i
    left join customer c on i.customer_id = c.customer_id
    left join invoiceentry ie on ie.invoice_id = i.invoice_id
where
    c.first_name ilike '%'
group by
    ie.invoice_id,
    i.invoice_id,
    c.customer_id,
    c.first_name,
    c.last_name,
    i.invoice_month,
    i.invoice_year
order by i.invoice_month;

-- GetInvoice
-- why null when no invoiceentry exists , how to avoid 
SELECT 
    c.customer_id,
    c.first_name,
    c.last_name,
    i.invoice_id,
    i.invoice_year,
    i.invoice_month,
    ie.work_day,
    ie.work_hours
FROM invoice i
    LEFT JOIN invoiceentry ie ON ie.invoice_id = i.invoice_id
    LEFT JOIN customer c ON i.customer_id = c.customer_id
where 
    i.customer_id = 100002;

SELECT 
    c.customer_id,
    c.first_name,
    c.last_name,
    i.invoice_id,
    i.invoice_year,
    i.invoice_month,
    ie.work_day,
    ie.work_hours
FROM invoice i
    LEFT JOIN invoiceentry ie ON ie.invoice_id = i.invoice_id
    LEFT JOIN customer c ON i.customer_id = c.customer_id
where 
    i.invoice_id = 100102
order by 
    ie.work_day asc;


SELECT 
    c.customer_id,
    c.first_name,
    c.last_name,
    i.invoice_id,
    i.invoice_year,
    i.invoice_month,
    -- COALESCE(ie.work_day,0),
    COALESCE(ie.work_hours,0)
FROM invoice i
    LEFT JOIN invoiceentry ie ON ie.invoice_id = i.invoice_id
    LEFT JOIN customer c ON i.customer_id = c.customer_id
group by i.invoice_month;
-- where 
--     i.invoice_id = 100103;
--    and ie.work_day is not null;
select
        c.customer_id,
        c.first_name,
        c.last_name,
        i.invoice_id,
        i.invoice_year,
        i.invoice_month
        -- COALESCE(ie.work_day,0),
        -- COALESCE(ie.work_hours,0)
    from
        invoice i
        left join customer c on i.customer_id = c.customer_id
        -- left join invoiceentry ie on ie.invoice_id = i.invoice_id
    where
        c.customer_id = 100007;
        -- i.invoice_id = 100101;
--- total hours per invoice
select
    i.invoice_id,
    SUM(ie.work_hours) as total_hours
from
    invoice i
    left join invoiceentry ie on ie.invoice_id = i.invoice_id
where
    i.invoice_id = 100102
group by
    i.invoice_id;
    