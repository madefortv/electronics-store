rm -f store.db
sqlite3 store.db 'CREATE TABLE products (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name VARCHAR(32) NOT NULL, description TEXT, price VARCHAR(8) NOT NULL);'

sqlite3 store.db 'CREATE TABLE deals (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name VARCHAR(32) NOT NULL DEFAULT "Regular Price", type VARCHAR(16) NOT NULL DEFAULT "Retail", coupon VARCHAR(8) NOT NULL DEFAULT "0.0", percent VARCHAR(8) NOT NULL DEFAULT "0.0", x INTEGER NOT NULL DEFAULT 0, y INTEGER NOT NULL DEFAULT 0, exclusive BOOLEAN NOT NULL DEFAULT 1);'

sqlite3 store.db 'CREATE TABLE offerings (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, product_id INTEGER NOT NULL, deal_id INTEGER NOT NULL, active BOOLEAN, modifed_price VARCHAR(8),  FOREIGN KEY (product_id) REFERENCES products (id), FOREIGN KEY (deal_id) REFERENCES deals (id) );'


sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("laptop", "very fast", "1000.00");'
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("mouse", "much clicky", "10.00");'
sqlite3 store.db 'INSERT INTO deals (name, type) VALUES ("Regular Price", "Retail");'
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (1, 1, 1);'

sqlite3 store.db 'SELECT * FROM deals;'

