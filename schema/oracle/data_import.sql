-- import TY product, price, MOQ, packaging
DECLARE
    my_cur_vendor_product_id vendor_product.VENDOR_PRODUCT_ID%TYPE;
    my_vendor_name vendor.NAME%TYPE := '天友玩具';
    my_product_type_name product_type.NAME%TYPE := '兵器玩具';
    my_material_type_name material_type.DESCRIPTION%TYPE := 'PU Foam';
BEGIN
    for row in (
        select * from ADMIN.TY_SPEC
        order by reference DESC
    ) LOOP
        insert into vendor_product
        (vendor_id, product_type_id, material_type_id, reference, description, length, width, height, weight)
        select v.id, pt.id, mt.type_id, row.reference, row.description_1,
             row.unit_l, row.unit_w, row.unit_h, row.unit_weight
        from VENDOR v, PRODUCT_TYPE pt, MATERIAL_TYPE mt
        where v.NAME = my_vendor_name
          and pt.NAME = my_product_type_name
          and mt.DESCRIPTION = my_material_type_name;
        
        select vp.vendor_product_id into my_cur_vendor_product_id
        from VENDOR_PRODUCT vp, VENDOR v
        where v.NAME = my_vendor_name
          and v.id = vp.vendor_id
          and vp.reference = row.reference;

        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        select my_cur_vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from currency ccy, price_type pt
        where ccy.iso_symbol = 'CNY'
          and pt.name = 'INC_TAX';

        IF row.MOQ IS NOT NULL then
            insert into vendor_product_moq
            (vendor_product_id, quantity, start_date)
            values
            (my_cur_vendor_product_id, row.moq, '01-JAN-2020');
        END IF;

        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.gross_weight IS NOT NULL or
           row.net_weight IS NOT NULL or
           row.qty_per_carton IS NOT NULL
         then
            insert into vendor_product_pack_detail
            (VENDOR_PRODUCT_ID,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             CARTON_GW,
             CARTON_NW,
             OUTER_QUANTITY,
             INNER_QUANTITY,
             NARRATIVE,
             start_date)
            values (my_cur_vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.gross_weight,
                   row.net_weight,
                   row.qty_per_carton,
                   NULL,
                   NULL,
                   '01-JAN-2020');
        END IF;
    END LOOP;
END;
/

-- import L product, price, MOQ, packaging
DECLARE
    my_cur_vendor_product_id vendor_product.VENDOR_PRODUCT_ID%TYPE;
    my_vendor_name vendor.NAME%TYPE := '义乌市鸿鲲有限公司';
    my_product_type_name product_type.NAME%TYPE := '派对用品';
BEGIN
  for row in (
         select * from ADMIN.L_SPEC
         order by reference DESC
  ) LOOP
        insert into vendor_product
        (vendor_id, product_type_id, reference, description, length, width, height)
        select v.id, pt.id, row.reference, row.description_1,
             row.unit_l, row.unit_w, row.unit_h
        from VENDOR v, PRODUCT_TYPE pt
        where v.NAME = my_vendor_name
          and pt.NAME = my_product_type_name;
          
        select vp.vendor_product_id into my_cur_vendor_product_id
        from VENDOR_PRODUCT vp, VENDOR v
        where v.NAME = my_vendor_name
          and v.id = vp.vendor_id
          and vp.reference = row.reference;

        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        select my_cur_vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from currency ccy, price_type pt
        where ccy.iso_symbol = 'CNY'
          and pt.name = 'INC_TAX';

        IF row.MOQ IS NOT NULL then
            insert into vendor_product_moq
            (vendor_product_id, quantity, start_date)
            values
            (my_cur_vendor_product_id, row.moq, '01-JAN-2020');
        END IF;

        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.qty_per_carton IS NOT NULL or
           row.description_2 IS NOT NULL
         then
            insert into vendor_product_pack_detail
            (VENDOR_PRODUCT_ID,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             OUTER_QUANTITY,
             NARRATIVE,
             start_date)
            values (my_cur_vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.qty_per_carton,
                   row.description_2,
                   '01-JAN-2020');
        END IF;
  END LOOP;
END;
/

-- import Z product, price, MOQ, packaging
DECLARE
    my_cur_vendor_product_id vendor_product.VENDOR_PRODUCT_ID%TYPE;
    my_vendor_name vendor.NAME%TYPE := '展希工艺品';
    my_product_type_name product_type.NAME%TYPE := '面具';
BEGIN
  for row in (
         select * from ADMIN.Z_SPEC
         order by reference DESC
  ) LOOP
       insert into vendor_product
       (vendor_id, product_type_id, reference, description, length, width, height)
        select v.id, pt.id, row.reference, row.description_1,
             row.unit_l, row.unit_w, row.unit_h
        from VENDOR v, PRODUCT_TYPE pt
        where v.NAME = my_vendor_name
          and pt.NAME = my_product_type_name;
          
        select vp.vendor_product_id into my_cur_vendor_product_id
        from VENDOR_PRODUCT vp, VENDOR v
        where v.NAME = my_vendor_name
          and v.id = vp.vendor_id
          and vp.reference = row.reference;

        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        select my_cur_vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from currency ccy, price_type pt
        where ccy.iso_symbol = 'CNY'
          and pt.name = 'INC_TAX';

        IF row.MOQ IS NOT NULL then
            insert into vendor_product_moq
            (vendor_product_id, quantity, start_date)
            values
            (my_cur_vendor_product_id, row.moq, '01-JAN-2020');
        END IF;

        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.qty_per_carton IS NOT NULL
         then
            insert into vendor_product_pack_detail
            (VENDOR_PRODUCT_ID,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             OUTER_QUANTITY,
             start_date)
            values (my_cur_vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.qty_per_carton,
                   '01-JAN-2020');
        END IF;
  END LOOP;
END;
/

-- import F product, price, MOQ, packaging
DECLARE
    my_cur_vendor_product_id vendor_product.VENDOR_PRODUCT_ID%TYPE;
    my_vendor_name vendor.NAME%TYPE := '泉州鑫彩';
    my_product_type_name product_type.NAME%TYPE := '面具';
BEGIN
  for row in (
         select * from ADMIN.F_SPEC
         order by reference DESC
  ) LOOP
       insert into vendor_product
       (vendor_id, product_type_id, reference, description, length, width, height)
        select v.id, pt.id, row.reference, row.description_1,
             row.unit_l, row.unit_w, row.unit_h
        from VENDOR v, PRODUCT_TYPE pt
        where v.NAME = my_vendor_name
          and pt.NAME = my_product_type_name;
          
        select vp.vendor_product_id into my_cur_vendor_product_id
        from VENDOR_PRODUCT vp, VENDOR v
        where v.NAME = my_vendor_name
          and v.id = vp.vendor_id
          and vp.reference = row.reference;

        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        select my_cur_vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from currency ccy, price_type pt
        where ccy.iso_symbol = 'CNY'
          and pt.name = 'EXC_TAX';

        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.qty_per_carton IS NOT NULL
         then
            insert into vendor_product_pack_detail
            (VENDOR_PRODUCT_ID,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             OUTER_QUANTITY,
             start_date)
            values (my_cur_vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.qty_per_carton,
                   '01-JAN-2020');
        END IF;
  END LOOP;
END;
/