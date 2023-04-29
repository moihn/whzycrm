DECLARE
    my_cur_vendor_product_id number(38);
    my_vendor_id         number(38);
    my_cur_price         number(38, 6);

    -- update below price information
    my_vendor_name       varchar2(255) := '泉州鑫彩';
    my_price_date        date        := '15-NOV-2021';
    my_price_type_name   varchar2(255) := 'EXC_TAX';

BEGIN
    -- get the my_vendor_id
    select v.id into my_vendor_id
    from VENDOR v
    where v.name = my_vendor_name;

    for row in (
        -- update the table name
        select *
        from ADMIN.F_QUOTE_20211115
    ) LOOP
    
        IF row.price IS NULL THEN
            CONTINUE;
        END IF;
        
        BEGIN
            select vp.vendor_product_id into my_cur_vendor_product_id
            from VENDOR_PRODUCT vp
            where vp.vendor_id = my_vendor_id
              and vp.reference = row.reference;
              
        EXCEPTION
        WHEN NO_DATA_FOUND THEN
            -- insert new product
            insert into vendor_product v
            (vendor_id, reference)
            values
            (my_vendor_id, row.reference)
            returning v.vendor_product_id into my_cur_vendor_product_id;
            
            dbms_output.put_line( 'ADDED NEW VENDOR PRODUCT: ' || row.reference);
        END;
            
        -- check if price already exists
        BEGIN
            select price into my_cur_price
            from vendor_product_price vpp
            where vpp.vendor_product_id = my_cur_vendor_product_id
              and vpp.start_date = my_price_date;
              
            update vendor_product_price vpp
            set vpp.price = row.price,
                vpp.currency_id = (select ccy.currency_id from currency ccy where ccy.iso_symbol = 'CNY'),
                vpp.price_type_id = (select pt.price_type_id from price_type pt where pt.name = my_price_type_name)
            where vpp.vendor_product_id = my_cur_vendor_product_id
              and vpp.start_date = my_price_date;
            dbms_output.put_line( 'UPDATED NEW PRICE FOR VENDOR PRODUCT: ' || row.reference || ' : ' || my_cur_vendor_product_id  || ' : ' || my_price_date || ' : ' || my_cur_price || ' to ' || row.price);
        EXCEPTION
        WHEN NO_DATA_FOUND THEN 
            insert into vendor_product_price
            (vendor_product_id, start_date, price, currency_id, price_type_id)
            select my_cur_vendor_product_id, my_price_date, row.price, ccy.currency_id, pt.price_type_id
            from currency ccy, price_type pt
            where ccy.iso_symbol = 'CNY'
              and pt.name = my_price_type_name;
            dbms_output.put_line( 'ADDED NEW PRICE FOR VENDOR PRODUCT: ' || row.reference || ' : ' || my_cur_vendor_product_id || ' : ' || row.price);
        END;
    END LOOP;
END;