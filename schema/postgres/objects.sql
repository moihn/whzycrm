create table if not exists VENDOR
(
    ID serial PRIMARY KEY,
    NAME text UNIQUE
);

create table if not exists PRODUCT_TYPE
(
    ID serial PRIMARY KEY,
    NAME text UNIQUE
);

create table if not exists UNIT_TYPE 
(
    UNIT_TYPE_ID serial PRIMARY KEY,
    NAME text NOT NULL
);

create table if not exists MATERIAL_TYPE 
(
    TYPE_ID serial PRIMARY KEY,
    DESCRIPTION text
);

create table if not exists PRICE_TYPE 
(
    PRICE_TYPE_ID serial PRIMARY KEY, 
    NAME          text, 
    INVOICE_RATE  numeric NOT NULL
);

create table if not exists CURRENCY 
(
    CURRENCY_ID serial PRIMARY KEY,
    ISO_SYMBOL  text UNIQUE,
    DESCRIPTION text
);

create table if not exists COUNTRY
(
    COUNTRY_ID serial PRIMARY KEY, 
    NAME       text
);

create table if not exists EXCHANGE_RATE 
(
    RATE_ID     serial PRIMARY KEY, 
    FROM_CCY_ID integer NOT NULL, 
    TO_CCY_ID   integer NOT NULL, 
    MULT_RATE   numeric NOT NULL, 
    START_DATE  date NOT NULL, 
        CONSTRAINT "U_EXCHANGE_RATE_1" UNIQUE (FROM_CCY_ID, TO_CCY_ID, START_DATE), 
        CONSTRAINT "FK_EXCHANGE_RATE_1" FOREIGN KEY (FROM_CCY_ID) REFERENCES CURRENCY (CURRENCY_ID), 
        CONSTRAINT "FK_EXCHANGE_RATE_2" FOREIGN KEY (TO_CCY_ID) REFERENCES CURRENCY (CURRENCY_ID)
);

create table if not exists VENDOR_PRODUCT 
(
    VENDOR_PRODUCT_ID serial PRIMARY KEY, 
    REFERENCE         text NOT NULL, 
    TEST_PERFORMED    boolean, 
    VENDOR_ID         integer NOT NULL, 
    DESCRIPTION       text, 
    MATERIAL_TYPE_ID  integer, 
    PRODUCT_TYPE_ID   integer NOT NULL, 
    UNIT_TYPE_ID      integer,
    LENGTH            numeric,
    WIDTH             numeric,
    HEIGHT            numeric,
    WEIGHT            numeric,
        CONSTRAINT "U_VENDOR_PRODUCT_1" UNIQUE (VENDOR_ID, REFERENCE),
        CONSTRAINT "VENDOR_PRODUCT_FK_UNIT" FOREIGN KEY (UNIT_TYPE_ID) REFERENCES UNIT_TYPE (UNIT_TYPE_ID),
        CONSTRAINT "VENDOR_PRODUCT_FK_MATERIAL" FOREIGN KEY (MATERIAL_TYPE_ID) REFERENCES MATERIAL_TYPE (TYPE_ID),
        CONSTRAINT "VENDOR_PRODUCT_FK_VENDOR" FOREIGN KEY (VENDOR_ID) REFERENCES VENDOR (ID)
);

create table if not exists VENDOR_PRODUCT_MOQ
(
    VENDOR_PRODUCT_ID integer NOT NULL, 
	QUANTITY          integer, 
	START_DATE        date NOT NULL, 
        CONSTRAINT "VENDOR_PRODUCT_MOQ_PK" PRIMARY KEY (VENDOR_PRODUCT_ID, START_DATE)
);

create table if not exists VENDOR_PRODUCT_PACK_DETAIL
(
    VENDOR_PRODUCT_ID integer,
    CARTON_LENGTH     numeric,
    CARTON_WIDTH      numeric,
    CARTON_HEIGHT     numeric,
    CARTON_GW         numeric,
    CARTON_NW         numeric,
    OUTER_QUANTITY    integer,
    INNER_QUANTITY    integer,
    NARRATIVE         text,
    START_DATE        date,
        CONSTRAINT        "VENDOR_PRODUCT_PACK_DETAIL" PRIMARY KEY (VENDOR_PRODUCT_ID, START_DATE)
);


create table if not exists VENDOR_PRODUCT_PRICE
(
    PRICE_ID          serial PRIMARY KEY, 
    VENDOR_PRODUCT_ID integer NOT NULL,
    START_DATE        date NOT NULL, 
    PRICE             numeric NOT NULL CHECK (PRICE > 0), 
    CURRENCY_ID       integer NOT NULL,
    PRICE_TYPE_ID     integer NOT NULL, 
        CONSTRAINT "U_VENDOR_PRODUCT_PRICE_1" UNIQUE (VENDOR_PRODUCT_ID, START_DATE),
        CONSTRAINT "FK_VENDOR_PRODUCT_PRICE_CCY" FOREIGN KEY (CURRENCY_ID) REFERENCES CURRENCY (CURRENCY_ID),
        CONSTRAINT "FK_VENDOR_PRODUCT_PRICE_PID" FOREIGN KEY (VENDOR_PRODUCT_ID) REFERENCES VENDOR_PRODUCT (VENDOR_PRODUCT_ID), 
        CONSTRAINT "FK_VENDOR_PRODUCT_PRICE_TYPE" FOREIGN KEY (PRICE_TYPE_ID) REFERENCES PRICE_TYPE (PRICE_TYPE_ID)
);

create table if not exists CLIENT
(
    CLIENT_ID  serial PRIMARY KEY, 
    NAME       text, 
    COUNTRY_ID integer
);

create table if not exists CLIENT_QUOTATION
(
    QUOTATION_ID serial PRIMARY KEY, 
    CLIENT_ID    integer,
    CURRENCY_ID  integer, 
    UPDATED_DATE date, 
    SENT         boolean,
        CONSTRAINT "U_CLIENT_QUOTATION_1" UNIQUE (CLIENT_ID, UPDATED_DATE),
        CONSTRAINT "CLIENT_QUOTATION_FK1" FOREIGN KEY (CURRENCY_ID) REFERENCES CURRENCY (CURRENCY_ID),
        CONSTRAINT "CLIENT_QUOTATION_FK2" FOREIGN KEY (CLIENT_ID) REFERENCES CLIENT (CLIENT_ID)
);

create table if not exists CLIENT_QUOTATION_ITEM
(
    QUOTATION_ID      integer NOT NULL, 
    VENDOR_PRODUCT_ID integer NOT NULL, 
    PRICE             numeric NOT NULL, 
    MOQ               integer,
    NARRATIVE         text,
        CONSTRAINT "CLIENT_QUOTATION_ITEM_PK" PRIMARY KEY (QUOTATION_ID, VENDOR_PRODUCT_ID),
        CONSTRAINT "CLIENT_QUOTATION_ITEM_FK1" FOREIGN KEY (QUOTATION_ID) REFERENCES CLIENT_QUOTATION (QUOTATION_ID),
        CONSTRAINT "CLIENT_QUOTATION_ITEM_FK2" FOREIGN KEY (VENDOR_PRODUCT_ID) REFERENCES VENDOR_PRODUCT (VENDOR_PRODUCT_ID)
);

create table if not exists CLIENT_PRODUCT
(
    CLIENT_PRODUCT_ID serial PRIMARY KEY,
    CLIENT_ID         integer NOT NULL,
    REFERENCE         text NOT NULL,
    DESCRIPTION       text,
    NARRATIVE         text,
    BARCODE           text,
        CONSTRAINT "U_CLIENT_PRODUCT_1" UNIQUE (CLIENT_ID, REFERENCE), 
        CONSTRAINT "FK_CLIENT_PRODUCT_1" FOREIGN KEY (CLIENT_ID) REFERENCES CLIENT (CLIENT_ID)
);

create table if not exists CLIENT_PRODUCT_ITEM
(
    CLIENT_PRODUCT_ITEM_ID serial PRIMARY KEY,
    CLIENT_PRODUCT_ID integer NOT NULL,
    VENDOR_PRODUCT_ID integer NOT NULL,
    NARRATIVE text,
        CONSTRAINT "FK_CLIENT_PRODUCT_ITEM_1" FOREIGN KEY (CLIENT_PRODUCT_ID) REFERENCES CLIENT_PRODUCT (CLIENT_PRODUCT_ID),
        CONSTRAINT "FK_CLIENT_PRODUCT_ITEM_2" FOREIGN KEY (VENDOR_PRODUCT_ID) REFERENCES VENDOR_PRODUCT (VENDOR_PRODUCT_ID)
);

create table if not exists CLIENT_ORDER_STATUS
(
    STATUS_ID   integer PRIMARY KEY,
    DESCRIPTION text NOT NULL
);

create table if not exists CLIENT_ORDER
(
    ORDER_ID               serial PRIMARY KEY,
    ORDER_REFERENCE        text NOT NULL, 
    CLIENT_ID              integer NOT NULL, 
    CLIENT_ORDER_REFERENCE text, 
    CREATION_DATE          date NOT NULL,
    SHIPMENT_DATE          date,
    STATUS_ID              integer NOT NULL, 
        CONSTRAINT "U_CLIENT_ORDER_REF" UNIQUE (ORDER_REFERENCE), 
        CONSTRAINT "U_CLIENT_ORDER_C_REF" UNIQUE (CLIENT_ID, CLIENT_ORDER_REFERENCE), 
        CONSTRAINT "FK_CLIENT_ORDER_C_ID" FOREIGN KEY (CLIENT_ID) REFERENCES CLIENT (CLIENT_ID),
        CONSTRAINT "FK_CLIENT_ORDER_S_ID" FOREIGN KEY (STATUS_ID) REFERENCES CLIENT_ORDER_STATUS (STATUS_ID)
);

create table if not exists CLIENT_ORDER_ITEM
(
    ORDER_ITEM_ID     serial PRIMARY KEY, 
    ORDER_ID          integer NOT NULL, 
    CLIENT_PRODUCT_ID integer NOT NULL, 
    QUANTITY          integer NOT NULL, 
    PRICE             numeric NOT NULL, 
    CURRENCY_ID       integer NOT NULL, 
    ADDED_DATE        date NOT NULL, 
    ALTERNATIVE_SHIP_DATE date, 
        CONSTRAINT "U_CLIENT_ORDER_ITEM_1" UNIQUE (ORDER_ID, CLIENT_PRODUCT_ID, ADDED_DATE, ALTERNATIVE_SHIP_DATE, PRICE, CURRENCY_ID), 
        CONSTRAINT "FK_CLIENT_ORDER_ITEM_1" FOREIGN KEY (ORDER_ID) REFERENCES CLIENT_ORDER (ORDER_ID), 
        CONSTRAINT "FK_CLIENT_ORDER_ITEM_2" FOREIGN KEY (CLIENT_PRODUCT_ID) REFERENCES CLIENT_PRODUCT (CLIENT_PRODUCT_ID)
);