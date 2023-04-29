do $$
declare
	vendor_id integer := 0;
	cny_ccy_id currency.currency_id%TYPE;
	price_type_id price_type.price_type_id%TYPE;
	product_id vendor_product.vendor_product_id%TYPE;
	spec record;
begin
	select vendor_id
	into vendor_id
	from vendor
	where name = '鸿鲲饰品';
	
	select currency_id
	into cny_ccy_id
	from currency
	where iso_symbol = 'CNY';
	
	select price_type.price_type_id
	into price_type_id
	from price_type
	where name = 'INC_TAX';
	
	for spec in select * from appuser.lr_catalog
	loop
		select vendor_product_id
		into product_id
		from vendor_product
		where reference = spec.reference;
		
		insert into vendor_product_price
		(vendor_product_id, start_date, price, currency_id, price_type_id)
		values
		(product_id, '2023-01-01', spec.price, cny_ccy_id, price_type_id);
	end loop;
end; $$;

do $$
declare
	vendor_id integer := 0;
	product_type_id integer := 0;
	cny_ccy_id currency.currency_id%TYPE;
	price_type_id price_type.price_type_id%TYPE;
	product_id vendor_product.vendor_product_id%TYPE;
	spec record;
begin
	select vendor_id
	into vendor_id
	from vendor
	where name = '天友玩具';
	
	select id
	into product_type_id
	from product_type
	where name = '兵器玩具';
	
	select currency_id
	into cny_ccy_id
	from currency
	where iso_symbol = 'CNY';
	
	select price_type.price_type_id
	into price_type_id
	from price_type
	where name = 'INC_TAX';
	
	for spec in select * from appuser.ty_catalog
	loop
		select vendor_product_id
		into product_id
		from vendor_product
		where reference = spec.reference;
		
		insert into vendor_product_pack_detail
		(vendor_product_id, carton_length, carton_width, carton_height, carton_gw, carton_nw, outer_quantity, inner_quantity, narrative, start_date)
		values
		(product_id, spec.carton_l, spec.carton_w, spec.carton_h, spec.gross_weight, spec.net_weight, spec.qty_per_carton, spec.qty_per_inner, NULL, '2023-01-01');
	end loop;
end; $$;