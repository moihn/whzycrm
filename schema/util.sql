SELECT df.tablespace_name "Tablespace",
  totalusedspace "Used MB",
  (df.totalspace - tu.totalusedspace) "Free MB",
  df.totalspace "Total MB",
  ROUND(100 * ( (df.totalspace - tu.totalusedspace)/ df.totalspace)) "% Free"
FROM
  (SELECT tablespace_name,
    ROUND(SUM(bytes) / 1048576) TotalSpace
  FROM dba_data_files
  GROUP BY tablespace_name
  ) df,
  (SELECT ROUND(SUM(bytes)/(1024*1024)) totalusedspace,
    tablespace_name
  FROM dba_segments
  GROUP BY tablespace_name
  ) tu
WHERE df.tablespace_name = tu.tablespace_name;


create or replace procedure add_client_order_item
(i_order_reference varchar2,
i_client_product_ref varchar2,
i_quantity number,
i_price number,
i_alternative_ship_date DATE default NULL,
i_ccy_iso varchar2 default 'USD',
i_added_date date default null)
as
    l_order_id client_order.order_id%TYPE;
    l_client_product_id client_product.client_product_id%TYPE;
    l_ccy_id currency.currency_id%TYPE;
    l_count number;
    l_added_date date;
begin
    begin
        select cp.client_product_id, co.order_id, co.creation_date
        into l_client_product_id, l_order_id, l_added_date
        from client_product cp, client_order co
        where cp.reference = i_client_product_ref
          and cp.client_id = co.client_id
          and co.order_reference = i_order_reference;
    exception
    when NO_DATA_FOUND then
        dbms_output.put_line('Missing product ' ||  i_client_product_ref || ' for order ' || i_order_reference);
        return;
    end;
    
    if i_added_date is not null then
        l_added_date := i_added_date;
    end if;

    select currency_id into l_ccy_id
    from currency
    where ISO_SYMBOL = i_ccy_iso;
    
    select count(*) into l_count
    from client_order_item
    where order_id = l_order_id
      and client_product_id = l_client_product_id;
    
    if l_count = 0 then
        insert into client_order_item
        (order_id, client_product_id, quantity, price, currency_id, added_date, alternative_ship_date)
        values
        (l_order_id, l_client_product_id, i_quantity, i_price, l_ccy_id, l_added_date, i_alternative_ship_date);
    else
        update client_order_item
        set quantity = i_quantity,
            price = i_price,
            currency_id = l_ccy_id,
            added_date = l_added_date,
            alternative_ship_date = i_alternative_ship_date
        where order_id = l_order_id
          and client_product_id = l_client_product_id;
    end if;
end;
/


create or replace procedure add_client_product_item
(i_vendor_id number,
i_vendor_product_ref varchar2,
i_client_product_item_narrative varchar2,
i_client_id number,
i_client_product_ref varchar2,
i_client_product_description varchar2,
i_client_product_narrative varchar2,
i_barcode varchar2 default null
)
as
    l_client_product_id client_product.client_product_id%TYPE;
    l_vendor_product_id vendor_product.vendor_product_id%TYPE;
    l_client_product_ref client_product.reference%TYPE;
    l_count number;
begin
    -- if client product reference is not specified, use same as vendor reference
    if i_client_product_ref is NULL then
        l_client_product_ref := i_vendor_product_ref;
    else
        l_client_product_ref := i_client_product_ref;
    end if;

    -- check if the client product is already defined
    -- create one if not, otherwise update existing information
    begin
        select client_product_id into l_client_product_id
        from client_product
        where client_id = i_client_id
          and reference = l_client_product_ref;
        
        if i_client_product_narrative <> NULL then
            update client_product
            set
              narrative = i_client_product_narrative
            where client_product_id = l_client_product_id;
        end if;
        
        if i_client_product_description <> NULL then
            update client_product
            set
              description = i_client_product_description
            where client_product_id = l_client_product_id;
        end if;
        
        if i_barcode <> NULL then
            update client_product
            set
              barcode = i_barcode
            where client_product_id = l_client_product_id;
        end if;
    exception
    when NO_DATA_FOUND then
        --insert client_product
        insert into client_product
        (client_id, reference, description, barcode)
        values
        (i_client_id, l_client_product_ref, i_client_product_description,
        i_barcode)
        returning client_product_id into l_client_product_id;
    end;
    
    -- find vendor product id
    begin
        select vendor_product_id into l_vendor_product_id
        from vendor_product
        where vendor_id = i_vendor_id
          and reference = i_vendor_product_ref;
        
        -- check if the vendor product is already under this client product
        select count(*) into l_count
        from client_product_item
        where client_product_id = l_client_product_id
          and vendor_product_id = l_vendor_product_id;
        
        -- if not yet, then add vendor product to client product
        if l_count = 0 then
            insert into client_product_item
            (client_product_id, vendor_product_id, narrative)
            values
            (l_client_product_id, l_vendor_product_id, i_client_product_item_narrative);
        end if;
    exception
    when NO_DATA_FOUND then
        dbms_output.put_line('Missing product ' ||  i_vendor_product_ref || ' from vendor ' || i_vendor_id);
        return;
    end;
end;
/

create or replace procedure update_vendor_product_price
(i_vendor_id number,
i_reference varchar2,
i_new_price number,
i_price_type_id number,
i_price_start_date date default trunc(sysdate),
i_price_ccy varchar2 default 'CNY')
as 
    l_vendor_product_id vendor_product.vendor_product_id%TYPE;
    l_currency_id vendor_product_price.currency_id%TYPE;
    l_price_date vendor_product_price.start_date%TYPE;
    l_last_price vendor_product_price.price%TYPE;
    l_last_price_ccy_id vendor_product_price.currency_id%TYPE;
    l_last_price_type_id vendor_product_price.price_type_id%TYPE;
    l_insert_new number := 0;
    l_update_old number := 0;
begin
    select vendor_product_id into l_vendor_product_id
    from vendor_product
    where vendor_id = i_vendor_id
      and reference = i_reference;
      
    begin
        select currency_id into l_currency_id
        from currency
        where UPPER(iso_symbol) = UPPER(i_price_ccy);
    exception
    when NO_DATA_FOUND then
        dbms_output.put_line('Missing currency ' ||  i_price_ccy);
        return;
    end;
      
    -- find latest price
    begin
        select max(start_date) into l_price_date
        from vendor_product_price
        where vendor_product_id = l_vendor_product_id
          and start_date <= i_price_start_date;
        
        select price, currency_id, price_type_id
        into l_last_price, l_last_price_ccy_id, l_last_price_type_id
        from vendor_product_price
        where vendor_product_id = l_vendor_product_id
          and start_date = l_price_date;
          
        if l_last_price <> i_new_price
            or l_last_price_ccy_id <> l_currency_id
            or l_last_price_type_id <> i_price_type_id then
            if l_price_date = i_price_start_date then
                l_update_old := 1;
            else
                l_insert_new := 1;
            end if;
        end if;
    exception
    when NO_DATA_FOUND then
        l_insert_new := 1;
    end;
    
    -- insert new
    if l_insert_new = 1 then
        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        values
        (l_vendor_product_id, i_price_start_date, i_new_price, l_currency_id, i_price_type_id);
    end if;
    
    if l_update_old = 1 then
        update vendor_product_price
        set
            price = i_new_price,
            currency_id = l_currency_id,
            price_type_id = i_price_type_id
        where vendor_product_id = l_vendor_product_id
          and start_date = i_price_start_date;
    end if;
exception
when NO_DATA_FOUND then
    dbms_output.put_line('Missing product ' ||  i_reference || ' for vendor ' || i_vendor_id);
end;
/


declare
    l_order_reference client_order.order_reference%TYPE := 'HP211123-T';
    l_client_id client.client_id%TYPE := 1;
begin
    add_client_product_item(2, 'DP-TY1038', NULL, l_client_id, '14425', 'Lion Shield', NULl, '771877144257');
    add_client_product_item(2, 'J-TYL002', NULL, l_client_id, '14430', 'Snake Sword', NULL, '771877144301');
    add_client_product_item(2, 'J-TY1043', NULL, l_client_id, '14470', 'Gladius Long Dagger', NULL, '771877144707');
    add_client_product_item(2, 'J-TY1002', NULL, l_client_id, '14480', 'Sting Sword', NULL, '771877144806');
    
    insert into client_order
    (order_reference, client_id, client_order_reference, creation_date, status_id)
    values
    (l_order_reference, l_client_id, 'PO7405', trunc(sysdate), 0);
    
    add_client_order_item(l_order_reference, '14470', 1000, 1.77);
    add_client_order_item(l_order_reference, '14480', 1000, 2.4);
end;
/

select co.client_order_reference, coi.price, ccy.iso_symbol, coi.quantity, cp.reference, vp.reference, v.name
from client_order co, client_order_item coi, client_product cp, client_product_item cpi, vendor_product vp, vendor v, currency ccy
where co.order_reference = 'HP211123-T'
  and co.order_id = coi.order_id
  and coi.client_product_id = cp.client_product_id
  and coi.currency_id = ccy.currency_id
  and cp.client_product_id = cpi.client_product_id
  and cpi.vendor_product_id = vp.vendor_product_id
  and vp.vendor_id = v.id;

declare
    l_order_reference client_order.order_reference%TYPE := 'HP210929-T';
    l_client_id client.client_id%TYPE := 3;
begin
    add_client_product_item(2, 'BS-TYL100', l_client_id, NULL,        'Lille dolk',         '5707978102280');
    add_client_product_item(2, 'D-TYL053',  l_client_id, 'D-TYL053B', 'Samuraisværd ',   '5707978102266');
    add_client_product_item(2, 'C-TY1007',  l_client_id, NULL,  'Alm hammer',         '5707978102303');
    add_client_product_item(2, 'C-TY1008',  l_client_id, NULL,  'Hammer m. løve',         '5707978102310');
    add_client_product_item(2, 'C-TYL018',  l_client_id, NULL,  'Hammer rå look',         '5707978102327');
    add_client_product_item(2, 'D-TY1055',  l_client_id, NULL,     'Lilla sværd', '5707978102235');
    add_client_product_item(2, 'D-TYL013',  l_client_id, NULL,     'Sværd m. rød plet',         '5707978102242');
    add_client_product_item(2, 'D-TYL021',  l_client_id, NULL,     'Lille piratsværd',         '5707978102259');
    add_client_product_item(2, 'D-TYL034',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'D-TYL044',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'DP-TY1043', l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'F-TYL042',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'F-TYL045',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'F-TYL095-1',l_client_id, NULL,     'Lille økse',         '5707978102273');
    add_client_product_item(2, 'J-TYL083',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'Q-TY1030-13',  l_client_id, NULL,     NULL,         NULL);
    add_client_product_item(2, 'Q-TY1030-23',  l_client_id, NULL,     'Stor pistol',         '5707978102341');
    add_client_product_item(2, 'J-TYL046',  l_client_id, NULL,     'Dragesværd',         '5707978102228');
    add_client_product_item(2, 'B-TYL010',  l_client_id, NULL,     'Kølle med dødningehoved',         '5707978102297');
    add_client_product_item(2, 'D-TYL026',  l_client_id, NULL,     NULL,        NULL);
    
    insert into client_order
    (order_reference, client_id, client_order_reference, creation_date, status_id)
    values
    (l_order_reference, l_client_id, 'PO05368', '29-SEP-2021', 0);
    
    add_client_order_item(l_order_reference, 'BS-TYL100', 1000, 0.79);
    add_client_order_item(l_order_reference, 'C-TY1007', 504,  4.54);
    add_client_order_item(l_order_reference, 'D-TY1055', 300, 2.5);
    add_client_order_item(l_order_reference, 'D-TYL013', 600, 2.9);
    add_client_order_item(l_order_reference, 'D-TYL021', 300, 2.35);
    add_client_order_item(l_order_reference, 'D-TYL026', 360, 2.74);
    add_client_order_item(l_order_reference, 'D-TYL034', 500, 2.74);
    add_client_order_item(l_order_reference, 'D-TYL044', 504, 2.5);
    add_client_order_item(l_order_reference, 'D-TYL053B', 2016, 2.9);
    add_client_order_item(l_order_reference, 'DP-TY1043', 240, 4.54);
    add_client_order_item(l_order_reference, 'F-TYL042', 300, 2.9);
    add_client_order_item(l_order_reference, 'F-TYL045', 400, 2.35);
    add_client_order_item(l_order_reference, 'F-TYL095-1', 500, 1.57);
    add_client_order_item(l_order_reference, 'J-TYL083', 500, 2.9);
    add_client_order_item(l_order_reference, 'J-TYL046', 300, 2.66);
end;
/