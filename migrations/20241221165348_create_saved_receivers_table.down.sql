ALTER TABLE saved_receivers DROP FOREIGN KEY fk_customer_id;
ALTER TABLE saved_receivers DROP FOREIGN KEY fk_bank_id;

DROP TABLE IF EXISTS saved_receivers;