insert into PRODUCT_TYPE(NAME)
          select '兵器玩具' from dual
union all select '面具'    from dual
union all select '派对用品' from dual;

insert into VENDOR(NAME)
          select '义乌市鸿鲲有限公司' from dual
union all select '天友玩具'         from dual
union all select '展希工艺品'       from dual
union all select '泉州鑫彩'         from dual;

insert into PRICE_TYPE(NAME, INVOICE_RATE)
          select 'INC_TAX', 0 from dual
union all select 'EXC_TAX', 0.13 from dual;

insert into MATERIAL_TYPE(DESCRIPTION)
          select 'PU Foam' from dual
union all select 'PS' from dual
union all select 'Alloy' from dual;

insert into CURRENCY(ISO_SYMBOL, DESCRIPTION)
          select 'CNY', 'Chinese Yuan' from dual
union all select 'USD', 'United State Dollar' from dual;

insert into COUNTRY(NAME)
          select 'USA' from dual
union all select 'CANADA' from dual
union all select 'FRANCE' from dual
union all select 'DANMARK' from dual
union all select 'SPAIN' from dual;

insert into CLIENT(NAME, COUNTRY_ID)
          select 'Creative Education', c.COUNTRY_ID from COUNTRY c where c.NAME = 'CANADA'
union all select 'PTIT CLOWN', c.COUNTRY_ID from COUNTRY c where c.NAME = 'FRANCE'
union all select 'CONXION', c.COUNTRY_ID from COUNTRY c where c.NAME = 'DANMARK'
union all select 'ATOSA', c.COUNTRY_ID from COUNTRY c where c.NAME = 'SPAIN';

insert into CLIENT_ORDER_STATUS(STATUS_ID, DESCRIPTION)
          select 0, 'RECEIVED' from dual
union all select 1, 'PI_SENT' from dual
union all select 2, 'PI_SIGNED' from dual
union all select 3, 'DEPOSIT_PAID' from dual
union all select 4, 'IN_PRODUCTION' from dual
union all select 5, 'SHIPPED' from dual
union all select 6, 'BALANCE_PAID' from dual
union all select 7, 'CLOSED' from dual
union all select 8, 'CANCELLED' from dual;

insert into EXCHANGE_RATE(FROM_CCY_ID, TO_CCY_ID, MULT_RATE, START_DATE)
select c1.CURRENCY_ID, c2.CURRENCY_ID, 6.35, sysdate
from CURRENCY c1, CURRENCY c2
where c1.ISO_SYMBOL = 'USD'
  and c2.ISO_SYMBOL = 'CNY';