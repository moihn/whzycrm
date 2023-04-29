alter table vendor modify name not null;
alter table product_type modify name not null;
alter table material_type modify description not null;
alter table price_type modify name not null;
alter table currency modify (iso_symbol not null, description not null);
alter table country modify name not null;
alter table vendor_product_moq modify vendor_product_id drop identity;
alter table vendor_product_moq add foreign key (vendor_product_id) references vendor_product(vendor_product_id);
alter table vendor_product_moq modify quantity not null;
alter table vendor_product_pack_detail add foreign key (vendor_product_id) references vendor_product(vendor_product_id);
alter table client modify (name not null, country_id not null);
alter table client_quotation modify (client_id not null, currency_id not null, updated_date not null, sent not null);

