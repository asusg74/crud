CREATE TABLE IF NOT EXISTS customers (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			phone TEXT NOT NULL UNIQUE,
			active BOOLEAN NOT NULL DEFAULT TRUE,
			created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)

INSERT INTO customers (name, phone, active) VALUES('harry', '+992000000002', FALSE)
UPDATE customers SET active = TRUE WHERE id = 2
DELETE FROM customers WHERE 
UPDATE customers SET phone = '+99212', name = 'not billy' WHERE id = 1