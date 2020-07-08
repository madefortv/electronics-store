rm -f store.db
sqlite3 store.db 'CREATE TABLE products (
                    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                    name VARCHAR(32) NOT NULL,
                    description TEXT,
                    price VARCHAR(8) NOT NULL
                );'

sqlite3 store.db 'CREATE TABLE deals (
                    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                    name VARCHAR(32) NOT NULL DEFAULT "Regular Price",
                    type VARCHAR(16) NOT NULL DEFAULT "Retail",
                    coupon VARCHAR(8) NOT NULL DEFAULT "0.0",
                    percent VARCHAR(8) NOT NULL DEFAULT "0.0",
                    x INTEGER NOT NULL DEFAULT 0,
                    y INTEGER NOT NULL DEFAULT 0,
                    exclusive BOOLEAN NOT NULL DEFAULT 1
                );'

sqlite3 store.db 'CREATE TABLE offerings (
                    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
                    product_id INTEGER NOT NULL, 
                    deal_id INTEGER NOT NULL, 
                    active BOOLEAN, 
                    modified_price VARCHAR(8) NOT NULL DEFAULT "NAN",
                        FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE RESTRICT,
                        FOREIGN KEY (deal_id) REFERENCES deals (id) ON UPDATE RESTRICT);'

sqlite3 store.db 'CREATE TABLE cart (
                    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
                    product_id INTEGER NOT NULL,
                    quantity INTEGER NOT NULL DEFAULT 1,
                    FOREIGN KEY (product_id) REFERENCES products (id)
                 );'

#seed
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("laptop", "very fast", "1000.00");'
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("mouse", "much clicky", "10.00");'
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("monitor", "four kay", "100.00");'
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("usb", "type see", "5.00");'
sqlite3 store.db 'INSERT INTO products (name, description, price) VALUES ("keyboard", "mecha", "15.00");'

sqlite3 store.db 'INSERT INTO deals (name, type) VALUES ("Regular Price", "Retail");'
sqlite3 store.db 'INSERT INTO deals (name, type) VALUES ("Get a mouse with every laptop", "Bundle");'
sqlite3 store.db 'INSERT INTO deals (name, type, coupon) VALUES ("$10 off a monitor", "Bundle", "10.00");'
sqlite3 store.db 'INSERT INTO deals (name, type, x, y) VALUES ("Buy 2 usb get 1 free", "BuyXGetY", 2, 1);'
sqlite3 store.db 'INSERT INTO deals (name, type, percent) VALUES ("50% off keyboards", "percent", "50");'
sqlite3 store.db 'INSERT INTO deals (name, type) VALUES ("10% off any full price item", "percent");'
# bundle mouse / laptop
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active, modified_price) VALUES (1, 2, 1, "1000.00");'
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active, modified_price) VALUES (2, 2, 1, "1000.00");' 
# retail items
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (1, 1, 1);'
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (2, 1, 1);'
# coupon on monitor
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (3, 3, 1);'
# buy 2 get 1 freee usb
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (4, 4, 1);'
# 50% off keyboards
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (5, 5, 1);'
# an invalid offer 
sqlite3 store.db 'INSERT INTO offerings (product_id, deal_id, active) VALUES (1, 6, 0);'



sqlite3 store.db '.headers ON'
sqlite3 store.db 'INSERT INTO cart (product_id, quantity) VALUES (1, 1);'
sqlite3 store.db 'INSERT INTO cart (product_id, quantity) VALUES (2, 1);'
sqlite3 store.db 'INSERT INTO cart (product_id, quantity) VALUES (3, 1);'
sqlite3 store.db 'INSERT INTO cart (product_id, quantity) VALUES (4, 3);'
sqlite3 store.db 'INSERT INTO cart (product_id, quantity) VALUES (5, 1);'


sqlite3 store.db -header -column 'SELECT PID, DID, PNAME, DNAME, price, quantity, type, coupon, percent, x, y, modified_price  FROM (SELECT products.id AS PID, products.name AS PNAME, products.price, deals.id AS DID, deals.name AS DNAME, deals.type, deals.x, deals.y, deals.coupon, deals.percent, offerings.modified_price FROM offerings INNER JOIN products on products.id = offerings.product_id INNER JOIN deals on deals.id = offerings.deal_id WHERE active = 1) INNER JOIN cart on cart.product_id = pid WHERE cart.quantity > 0;'


#sqlite3 store.db 'SELECT products.id, products.name, products.price, cart. quantity FROM 
#    (SELECT products.id as pid, products.name, products.price, deals.id, deals.name, deals.type FROM offerings INNER JOIN products on products.id = offerings.product_id INNER JOIN deals on deals.id = offerings.deal_id WHERE active = 1) INNER JOIN offers.pid on cart.product_id;'
