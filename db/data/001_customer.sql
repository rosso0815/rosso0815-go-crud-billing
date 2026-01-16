--

delete from customer;

insert into customer ( first_name,last_name,street,town,remark,email,phone,ahv,alv,quellensteuer,salary_hour)
    values
        ('customer_01','last_name1','street1','town1','remark1','mail1','phone1',5.3,1.1,5,30),
        ('customer_02','last_name2','street2','town2','remark2','mail2','phone2',5.3,1.1,5,30),
        ('customer_03','last_name3','street3','town3','remark3','mail3','phone3',5.3,1.1,5,35);

insert into customer ( first_name,last_name)
    values
        ('customer_11','last_name11'),
        ('customer_12','last_name12'),
        ('customer_13','last_name13'),
        ('customer_14','last_name14');

--
