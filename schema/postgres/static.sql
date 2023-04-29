insert into currency
(iso_symbol, description)
values
('CNY', 'Chinese Yuan'),
('USD', 'US Dollar');

insert into country
(name)
values
('Spain'),
('Danmark'),
('France'),
('Germany'),
('Mexico'),
('Brazil');

insert into material_type
(description)
values
('PU Foam'),
('PS'),
('Alloy');

insert into price_type
(name, invoice_rate)
values
('INC_TAX', 0),
('EXC_TAX', 0.13);

insert into CLIENT_ORDER_STATUS
(STATUS_ID, DESCRIPTION)
values
(0, 'RECEIVED'),
(1, 'PI_SENT'),
(2, 'PI_SIGNED'),
(3, 'DEPOSIT_PAID'),
(4, 'IN_PRODUCTION'),
(5, 'SHIPPED'),
(6, 'BALANCE_PAID'),
(7, 'CLOSED'),
(8, 'CANCELLED');

insert into PRODUCT_TYPE
(NAME)
values
('兵器玩具'),
('面具'),
('派对用品');
