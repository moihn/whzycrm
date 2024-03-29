grant select,insert,update,delete on
  client,
  client_order,
  client_order_item,
  client_order_status,
  client_product,
  client_product_item,
  client_quotation,
  client_quotation_item,
  country,
  currency,
  exchange_rate,
  material_type,
  price_type,
  product_type,
  unit_type,
  vendor,
  vendor_product,
  vendor_product_moq,
  vendor_product_pack_detail,
  vendor_product_price
  to :appuser;
  
grant usage, select on sequence
  vendor_id_seq,
  product_type_id_seq,
  unit_type_unit_type_id_seq,
  material_type_type_id_seq,
  price_type_price_type_id_seq,
  currency_currency_id_seq,
  country_country_id_seq,
  exchange_rate_rate_id_seq,
  vendor_product_vendor_product_id_seq,
  vendor_product_price_price_id_seq,
  client_client_id_seq,
  client_quotation_quotation_id_seq,
  client_product_client_product_id_seq,
  client_product_item_client_product_item_id_seq,
  client_order_order_id_seq,
  client_order_item_order_item_id_seq
  to :appuser;
