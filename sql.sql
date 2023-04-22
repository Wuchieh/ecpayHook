CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "public"."ecpay_paid_data"
(
    "id"                uuid                                        NOT NULL DEFAULT uuid_generate_v4(),
    "merchant_trade_no" varchar(20) COLLATE "pg_catalog"."default"  NOT NULL,
    "name"              varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "total_amount"      int8                                        NOT NULL,
    "donate_to"         varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "payment_time"      timestamp(6)                                NOT NULL DEFAULT now(),
    "message"           varchar(60) COLLATE "pg_catalog"."default",
    "simulate_paid"     bool                                                 DEFAULT false,
    CONSTRAINT "ecpayTable_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."ecpay_paid_data"
    OWNER TO "wuchieh";