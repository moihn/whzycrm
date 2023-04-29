#!/bin/sh
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# run table creation as dba user
DB_NAME=${DB_NAME:-appdb}
TMP_SCHEMA=${TMP_SCHEMA:-appuser}
PSQL_CMD=${PSQL_CMD:-docker compose exec -T db psql}

${PSQL_CMD} -U ${DB_NAME}_appuser -d ${DB_NAME} <<EOF
create table if not exists ${TMP_SCHEMA}.zx_catalog
(
Reference text,
Description_1 text,
Unit_L numeric,
Unit_W numeric,
Unit_H numeric,
Test_Performed text,
Materials text,
PriceU numeric,
Price numeric,
Description_2 text,
Price_Type integer,
Package text,
MOQ integer,
Unit text,
Qty_Per_Carton integer,
Carton_L numeric,
Carton_W numeric,
Carton_H numeric,
Net_Weight numeric
);
EOF

${PSQL_CMD} -U ${DB_NAME}_appuser -d ${DB_NAME} -c "COPY ${TMP_SCHEMA}.zx_catalog(Reference,Description_1,Materials,Test_Performed,PriceU,Price,Price_Type,Package,MOQ,Unit,Qty_Per_Carton,Carton_L,Carton_W,Carton_H,Net_Weight,Unit_L,Unit_W,Unit_H) FROM STDIN WITH (FORMAT CSV, HEADER true);" <${SCRIPT_DIR}/zx_spec.csv

${PSQL_CMD} -U ${DB_NAME}_appuser -d ${DB_NAME} <<EOF
do \$\$
declare
	vendor_id vendor.id%TYPE;
	product_type_id product_type.id%TYPE;
	cny_ccy_id currency.currency_id%TYPE;
	price_type_id price_type.price_type_id%TYPE;
	product_id vendor_product.vendor_product_id%TYPE;
	spec record;
begin
	select id
	into vendor_id
	from vendor
	where name = '展希';
	
	select id
	into product_type_id
	from product_type
	where name = '面具';
	
	select currency_id
	into cny_ccy_id
	from currency
	where iso_symbol = 'CNY';
	
	select price_type.price_type_id
	into price_type_id
	from price_type
	where name = 'INC_TAX';
	
	for spec in select * from ${TMP_SCHEMA}.zx_catalog
	loop
		insert into vendor_product
		(reference, vendor_id, product_type_id, description,
		length, width, height, weight)
		values
		(spec.reference, vendor_id, product_type_id, spec.description_1,
		spec.unit_l, spec.unit_w, spec.unit_h, NULL)
		returning vendor_product_id into product_id;
		
		if spec.price > 0 then
			insert into vendor_product_price
			(vendor_product_id, start_date, price, currency_id, price_type_id)
			values
			(product_id, '2023-01-01', spec.price, cny_ccy_id, price_type_id);
		end if;
		
		insert into vendor_product_pack_detail
		(vendor_product_id, carton_length, carton_width, carton_height,
		carton_gw, carton_nw, outer_quantity, inner_quantity,
		narrative, start_date)
		values
		(product_id, spec.carton_l, spec.carton_w, spec.carton_h,
		NULL, NULL, spec.qty_per_carton, NULL,
		NULL, '2023-01-01');
		
		if spec.moq > 0 then
			insert into vendor_product_moq
			(vendor_product_id, quantity, start_date)
			values
			(product_id, spec.moq, '2023-01-01');
		end if;
	end loop;
end; \$\$
EOF
