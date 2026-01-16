--
delete from invoice cascade;

delete from invoiceentry cascade;

insert into invoice
(customer_id, invoice_month,invoice_year,remark,ahv,alv,quellensteuer,hours_total,salary_hour,bill_sum)
values
   (100002,1,2025,'remark1',5.3,1.1,5.0,0,0.,0.0),
   (100002,2,2025,'remark2',5.3,1.1,5.0,6.0,30,6.0),
   (100003,3,2025,'remark3',5.3,1.1,5.0,0,0.0,0.0);

insert into invoiceentry(invoice_id,work_day,work_hours)
values
   (100102,3,1.5),
   (100102,4,1.5),
   (100102,28,1.5),
   (100103,28,1.5);


