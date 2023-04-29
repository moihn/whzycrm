-- create schema user
alter user whzy IDENTIFIED by "hn_@Orac1eC1oud";
ALTER USER whzy quota unlimited on data;

-- create sequence object
drop sequence whzy.seq_product_id;
create sequence whzy.seq_product_id;

drop sequence whzy.seq_vendor_product_price_id;
create sequence whzy.seq_vendor_product_price_id;

drop sequence whzy.seq_vendor_product_id;
create sequence whzy.seq_vendor_product_id;

drop sequence whzy.seq_vendor_id;
create sequence whzy.seq_vendor_id;

drop sequence whzy.seq_client_order_id;
create sequence whzy.seq_client_order_id;

-- create table object

create table whzy.vendor
(id number(38) primary KEY, name VARCHAR2(255) UNIQUE);

create table whzy.product_type
(id number(38) primary KEY, name VARCHAR2(255) UNIQUE);

-- Create VENDOR_PRODUCT table
CREATE TABLE WHZY.VENDOR_PRODUCT 
    (
     VENDOR_PRODUCT_ID NUMBER (38) PRIMARY KEY, 
     REFERENCE         VARCHAR2 (32) , 
     MATERIAL          VARCHAR2 (255) , 
     TEST_PERFORMED    VARCHAR2 (2) , 
     PRICE             NUMBER (38) , 
     PRICE_TYPE        VARCHAR2 (4000) , 
     PACKAGE           VARCHAR2 (4000) , 
     MOQ               VARCHAR2 (4000) , 
     VENDOR_ID         NUMBER (38) , 
     DESCRIPTION       VARCHAR2 (4000) , 
     PRODUCT_TYPE_ID   NUMBER , 
     UNIT_TYPE_ID      NUMBER ,
     LENGTH            NUMBER(38, 6),
     WIDTH             NUMBER(38, 6),
     HEIGHT            NUMBER(38, 6),
     WEIGHT            NUMBER(38, 6)
    ) 
;

CREATE UNIQUE INDEX WHZY.U_VENDOR_PRODUCT_1 ON WHZY.VENDOR_PRODUCT 
    (
     VENDOR_ID ASC,
     REFERENCE ASC 
    ) 
;

ALTER TABLE WHZY.VENDOR_PRODUCT 
    ADD CONSTRAINT VENDOR_PRODUCT_FK_UNIT FOREIGN KEY 
    ( 
     UNIT_TYPE_ID
    ) 
    REFERENCES WHZY.UNIT_TYPE ( UNIT_TYPE_ID ) 
    NOT DEFERRABLE 
;

-- Create VENDOR_PRODUCT_MOQ table
CREATE TABLE "WHZY"."VENDOR_PRODUCT_MOQ" 
(
    "VENDOR_PRODUCT_ID" NUMBER(38,0) NOT NULL ENABLE, 
	"QUANTITY" NUMBER(38,0), 
	"START_DATE" DATE NOT NULL ENABLE, 
	 CONSTRAINT "VENDOR_PRODUCT_MOQ_PK" PRIMARY KEY ("VENDOR_PRODUCT_ID", "START_DATE")
);

-- create vendor_product_pack_detail table
create table whzy.vendor_product_pack_detail
(
  vendor_product_id number(38),
  CARTON_LENGTH number(38,6),
  CARTON_WIDTH number(38, 6),
  CARTON_HEIGHT number(38, 6),
  CARTON_GW number(38, 6),
  CARTON_NW number(38, 6),
  OUTER_QUANTITY number(38),
  INNER_QUANTITY number(38),
  NARRATIVE varchar2(4000),
  start_date date,
  CONSTRAINT "vendor_product_pack_detail" PRIMARY KEY ("VENDOR_PRODUCT_ID", "START_DATE")
);

-- clean data
truncate table whzy.vendor_product_price;

-- upload static reference data

insert into whzy.product_type(id, name)
          select 0, '兵器玩具' from dual
union all select 1, '面具'    from dual
union all select 2, '派对用品' from dual;

insert into whzy.vendor(id, name)
          select 0, '义乌市鸿鲲有限公司' from dual
union all select 1, '天友玩具'         from dual
union all select 2, '展希工艺品'       from dual;

insert into whzy.price_type(price_type_id, name)
          select 0, 'INC_TAX' from dual
union all select 1, 'EXC_TAX' from dual;

insert into whzy.material_type(type_id, description)
          select 0, 'PU Foam' from dual
union all select 1, 'PS' from dual
union all select 2, 'Alloy' from dual;

insert into whzy.CURRENCY(currency_id, iso_symbol, description)
          select 0, 'CNY', 'Chinese Yuan' from dual
union all select 1, 'USD', 'United State Dollar' from dual;

insert into whzy.COUNTRY(COUNTRY_ID, name)
          select 0, 'USA' from dual
union all select 1, 'CANADA' from dual
union all select 2, 'FRANCE' from dual
union all select 3, 'DANMARK' from dual
union all select 4, 'SPAIN' from dual;

insert into whzy.CLIENT(client_id, name, country_id)
          select 0, 'Creative Education', 1 from dual
union all select 1, 'PTIT CLOWN', 2 from dual
union all select 2, 'CONXION', 3 from dual
union all select 3, 'ATOSA', 4 from dual;

-- import TY product, price, MOQ
DECLARE
    vendor_product_id number(38);
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
        where v.NAME = '天友玩具'
          and pt.NAME = '兵器玩具'
          and mt.DESCRIPTION = 'PU Foam';
        
        select vp.vendor_product_id into vendor_product_id
        from VENDOR_PRODUCT vp, VENDOR v
        where v.NAME = '天友玩具'
          and v.id = vp.vendor_id
          and vp.reference = row.reference;

        insert into vendor_product_price
        (vendor_product_id, start_date, price, currency_id, price_type_id)
        select vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from currency ccy, price_type pt
        where ccy.iso_symbol = 'CNY'
          and pt.name = 'INC_TAX';

        IF row.MOQ IS NOT NULL then
            insert into vendor_product_moq
            (vendor_product_id, quantity, start_date)
            values
            (vendor_product_id, row.moq, '01-JAN-2020');
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
            values (vendor_product_id,
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

-- import L product and price
BEGIN
  for row in (
         select * from L_SPEC
         order by reference DESC
  ) LOOP
       insert into whzy.vendor_product
       (vendor_product_id, vendor_id, product_type_id, material_type_id, reference, description)
       values (whzy.seq_vendor_product_id.nextval, 0, 2, 3, row.reference, row.description_1);
  END LOOP;
END;

BEGIN
    for row in (
        select * from L_SPEC
        order by reference DESC
    ) LOOP
        insert into whzy.vendor_product_price
        (price_id, vendor_product_id, start_date, price, currency_id, price_type_id)
        select whzy.seq_vendor_product_price_id.nextval, vp.vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from whzy.vendor_product vp, whzy.currency ccy, whzy.price_type pt
        where vp.VENDOR_ID = 0
          and vp.REFERENCE = row.REFERENCE
          and ccy.iso_symbol = 'CNY'
          and pt.name = 'INC_TAX';
    END LOOP;
END;

BEGIN
    for row in (
        select * from L_SPEC
        order by reference DESC
    ) LOOP
        IF row.MOQ IS NOT NULL then
            insert into whzy.vendor_product_moq
            (vendor_product_id, quantity, start_date)
            select vp.vendor_product_id, row.moq, '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 0
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from L_SPEC
        order by reference DESC
    ) LOOP
        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.qty_per_carton IS NOT NULL or
           row.description_2 IS NOT NULL
         then
            insert into whzy.vendor_product_pack_detail
            (vendor_product_id,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             OUTER_QUANTITY,
             NARRATIVE,
             start_date)
            select vp.vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.qty_per_carton,
                   row.description_2,
                   '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 0
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from L_SPEC
        order by reference DESC
    ) LOOP
        IF row.unit_l IS NOT NULL or
           row.unit_w IS NOT NULL or
           row.unit_h IS NOT NULL
         then
            update whzy.vendor_product
            set length = row.unit_l,
                width = row.unit_w,
                height = row.unit_h
            where VENDOR_ID = 0
            and REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

-- import F product and price
BEGIN
  for row in (
         select * from F_SPEC
         order by reference DESC
  ) LOOP
       insert into whzy.vendor_product
       (vendor_product_id, vendor_id, product_type_id, material_type_id, reference, description)
       values (whzy.seq_vendor_product_id.nextval, 3, 1, 4, row.reference, row.description_1);
  END LOOP;
END;

BEGIN
    for row in (
        select * from F_SPEC
        order by reference DESC
    ) LOOP
        insert into whzy.vendor_product_price
        (price_id, vendor_product_id, start_date, price, currency_id, price_type_id)
        select whzy.seq_vendor_product_price_id.nextval, vp.vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from whzy.vendor_product vp, whzy.currency ccy, whzy.price_type pt
        where vp.VENDOR_ID = 3
          and vp.REFERENCE = row.REFERENCE
          and ccy.iso_symbol = 'CNY'
          and pt.name = 'EXC_TAX';
    END LOOP;
END;

BEGIN
    for row in (
        select * from F_SPEC
        order by reference DESC
    ) LOOP
        IF row.MOQ IS NOT NULL then
            insert into whzy.vendor_product_moq
            (vendor_product_id, quantity, start_date)
            select vp.vendor_product_id, row.moq, '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 3
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from F_SPEC
        order by reference DESC
    ) LOOP
        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.qty_per_carton IS NOT NULL or
           row."INNER" IS NOT NULL
         then
            insert into whzy.vendor_product_pack_detail
            (vendor_product_id,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             CARTON_GW,
             CARTON_NW,
             OUTER_QUANTITY,
             INNER_QUANTITY,
             NARRATIVE,
             start_date)
            select vp.vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   NULL,
                   NULL,
                   row.qty_per_carton,
                   row."INNER",
                   NULL,
                   '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 3
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from F_SPEC
        order by reference DESC
    ) LOOP
        IF row.unit_l IS NOT NULL or
           row.unit_w IS NOT NULL or
           row.unit_h IS NOT NULL
         then
            update whzy.vendor_product
            set length = row.unit_l,
                width = row.unit_w,
                height = row.unit_h
            where VENDOR_ID = 3
            and REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

-- import Z product and price
BEGIN
  for row in (
         select * from Z_SPEC
         order by reference DESC
  ) LOOP
       insert into whzy.vendor_product
       (vendor_product_id, vendor_id, product_type_id, material_type_id, reference, description)
       values (whzy.seq_vendor_product_id.nextval, 2, 1, 1, row.reference, row.description_1);
  END LOOP;
END;

BEGIN
    for row in (
        select * from Z_SPEC
        order by reference DESC
    ) LOOP
        insert into whzy.vendor_product_price
        (price_id, vendor_product_id, start_date, price, currency_id, price_type_id)
        select whzy.seq_vendor_product_price_id.nextval, vp.vendor_product_id, '01-JAN-2020', row.price, ccy.currency_id, pt.price_type_id
        from whzy.vendor_product vp, whzy.currency ccy, whzy.price_type pt
        where vp.VENDOR_ID = 2
          and vp.REFERENCE = row.REFERENCE
          and ccy.iso_symbol = 'USD'
          and pt.name = 'INC_TAX';
    END LOOP;
END;

BEGIN
    for row in (
        select * from Z_SPEC
        order by reference DESC
    ) LOOP
        IF row.MOQ IS NOT NULL then
            insert into whzy.vendor_product_moq
            (vendor_product_id, quantity, start_date)
            select vp.vendor_product_id, row.moq, '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 2
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from Z_SPEC
        order by reference DESC
    ) LOOP
        IF row.carton_l IS NOT NULL or
           row.carton_w IS NOT NULL or
           row.carton_h IS NOT NULL or
           row.net_weight IS NOT NULL or
           row.qty_per_carton IS NOT NULL or
           row."package" IS NOT NULL
         then
            insert into whzy.vendor_product_pack_detail
            (vendor_product_id,
             CARTON_LENGTH,
             CARTON_WIDTH,
             CARTON_HEIGHT,
             CARTON_NW,
             OUTER_QUANTITY,
             NARRATIVE,
             start_date)
            select vp.vendor_product_id,
                   row.carton_l,
                   row.carton_w,
                   row.carton_h,
                   row.net_weight,
                   row.qty_per_carton,
                   row."package",
                   '01-JAN-2019'
            from whzy.vendor_product vp
            where vp.VENDOR_ID = 2
            and vp.REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

BEGIN
    for row in (
        select * from TY_SPEC
        order by reference DESC
    ) LOOP
        IF row.unit_l IS NOT NULL or
           row.unit_w IS NOT NULL or
           row.unit_h IS NOT NULL
         then
            update whzy.vendor_product
            set length = row.unit_l,
                width = row.unit_w,
                height = row.unit_h
            where VENDOR_ID = 2
            and REFERENCE = row.REFERENCE;
        END IF;
    END LOOP;
END;

-- create QUOTATION table
  CREATE TABLE "CLIENT_QUOTATION" 
   (
    "QUOTATION_ID" NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY START WITH 1 NOCACHE NOT NULL, 
	"CLIENT_ID" NUMBER, 
	"UPDATED_DATE" DATE, 
	"SENT" CHAR(1),
	 CONSTRAINT "QUOTATION_PK" PRIMARY KEY ("QUOTATION_ID")
   );

-- create PRODUCT_IN_QUOTATION table
CREATE TABLE "CLIENT_QUOTATION_ITEM" 
   (
   	"QUOTATION_ID" NUMBER, 
	"VENDOR_PRODUCT_ID" NUMBER, 
	"PRICE" NUMBER(38,6), 
	"CURRENCY_ID" NUMBER, 
	"NARRATIVE" VARCHAR2(4000) ,
	 CONSTRAINT "PRODUCT_IN_QUOTE_PK" PRIMARY KEY ("QUOTATION_ID", "VENDOR_PRODUCT_ID"),
     CONSTRAINT "PRODUCT_IN_QUOTE_FK1" FOREIGN KEY ("QUOTATION_ID") REFERENCES "CLIENT_QUOTATION" ("QUOTATION_ID"),
     CONSTRAINT "PRODUCT_IN_QUOTE_FK2" FOREIGN KEY ("VENDOR_PRODUCT_ID") REFERENCES "VENDOR_PRODUCT" ("VENDOR_PRODUCT_ID")
   );

-- create CLIENT_ORDER table
  CREATE TABLE "CLIENT_ORDER" 
   (
    "ORDER_ID" NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY START WITH 1 NOCACHE NOT NULL,
	"ORDER_REFERENCE" VARCHAR2(255) NOT NULL, 
	"CLIENT_ID" NUMBER NOT NULL, 
	"CLIENT_ORDER_REFERENCE" VARCHAR2(255), 
	"CREATION_DATE" DATE NOT NULL,
    "SHIPMENT_DATE" DATE,
	"STATUS" NUMBER DEFAULT ON NULL 0 NOT NULL, 
	 CONSTRAINT "CLIENT_ORDER_PK" PRIMARY KEY ("ORDER_ID"), 
	 CONSTRAINT "U_CLIENT_ORDER_REF" UNIQUE ("ORDER_REFERENCE"), 
	 CONSTRAINT "U_CLIENT_ORDER_C_REF" UNIQUE ("CLIENT_ID", "CLIENT_ORDER_REFERENCE"), 
	 CONSTRAINT "FK_CLIENT_ORDER_C_ID" FOREIGN KEY ("CLIENT_ID") REFERENCES "CLIENT" ("CLIENT_ID")
   );

-- create CLIENT_ORDER_ITEM table
  CREATE TABLE "CLIENT_ORDER_ITEM" 
   (
    "ORDER_ITEM_ID" NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY START WITH 1 NOCACHE NOT NULL, 
	"ORDER_ID" NUMBER NOT NULL, 
	"VENDOR_PRODUCT_ID" NUMBER NOT NULL, 
	"QUANTITY" NUMBER NOT NULL, 
	"PRICE" NUMBER(38,6) NOT NULL, 
	"CURRENCY_ID" NUMBER NOT NULL, 
	"ADDED_DATE" DATE NOT NULL, 
	"ALTERNATIVE_SHIP_DATE" DATE, 
	 CONSTRAINT "CLIENT_ORDER_ITEM_PK" PRIMARY KEY ("ORDER_ITEM_ID"), 
	 CONSTRAINT "U_ORDER_ITEM_1" UNIQUE ("ORDER_ID", "VENDOR_PRODUCT_ID", "ADDED_DATE", "ALTERNATIVE_SHIP_DATE", "PRICE", "CURRENCY_ID"), 
	 CONSTRAINT "FK_CLIENT_ORDER_ITEM_1" FOREIGN KEY ("ORDER_ID") REFERENCES "CLIENT_ORDER" ("ORDER_ID"), 
	 CONSTRAINT "FK_CLIENT_ORDER_ITEM_2" FOREIGN KEY ("VENDOR_PRODUCT_ID") REFERENCES "VENDOR_PRODUCT" ("VENDOR_PRODUCT_ID")
   );

-- create EXCHANGE_RATE table
  CREATE TABLE "EXCHANGE_RATE" 
   (
    "RATE_ID" NUMBER GENERATED BY DEFAULT ON NULL AS IDENTITY START WITH 1 NOCACHE NOT NULL, 
	"FROM_CCY_ID" NUMBER NOT NULL, 
	"TO_CCY_ID" NUMBER NOT NULL, 
	"MULT_RATE" NUMBER NOT NULL, 
	"START_DATE" DATE NOT NULL, 
	 CONSTRAINT "EXCHANGE_RATE_PK" PRIMARY KEY ("RATE_ID"), 
	 CONSTRAINT "U_EXCHANGE_RATE_1" UNIQUE ("FROM_CCY_ID", "TO_CCY_ID", "START_DATE"), 
	 CONSTRAINT "FK_EXCHANGE_RATE_1" FOREIGN KEY ("FROM_CCY_ID") REFERENCES "CURRENCY" ("CURRENCY_ID"), 
	 CONSTRAINT "FK_EXCHANGE_RATE_2" FOREIGN KEY ("TO_CCY_ID") REFERENCES "CURRENCY" ("CURRENCY_ID")
   );

-- calculate quotation price
DECLARE
    rate NUMBER(38, 6);
    out_price NUMBER(38, 6);
begin

        select 1/MULT_RATE
        into rate
        from EXCHANGE_RATE
        where FROM_CCY_ID = 1
          and TO_CCY_ID = 0;
DBMS_OUTPUT.put_line (rate);
    if rate is not NULL then
        select vpp.price * rate * (1 + pt.invoice_rate) into out_price
        from VENDOR_PRODUCT vp, VENDOR_PRODUCT_PRICE vpp, PRICE_TYPE pt
        where vp.vendor_id = 1
          and vp.reference = 'D-TYL062'
          and vpp.vendor_product_id = vp.vendor_product_id
          and vpp.price_type_id = pt.price_type_id;
        if out_price is not NULL then
           DBMS_OUTPUT.put_line (out_price);
        else
           DBMS_OUTPUT.put_line ('no price found');
        end if;
    else
        DBMS_OUTPUT.put_line ('no fx rate found: ' || rate);
    end if;
end;