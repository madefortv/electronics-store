## Electronics Store

# Dependencies
go 1.4
sqlite

# set up
install [go]('https://golang.org/doc/install')
- clone repository or use `go get`
- compile the binary
- seed the database
- run the server
```bash
go build
chmod +x create_database.sh
./create_database.sh
./store
```

The serve is listening on localhost:8000


# run the tests
```bash
go test
```

# Structure
- main.go builds dependencies and injects into the server to run
- server.go provides a router for handling different endpoints like: `http://localhost:8000/{products,cart,offerings,deals}`
- service.go provides some abstraction to the database layer
- models.go hosts the datamodels and table building functions
- db.go is where the sql queries live
- config.go is the server/db config file
- server_test.go is the test file

# Assumptions
Products can only have one deal applied to them at any given time. Calculating Attribution is hard and costs extra.

# Approach
a vanilla Go web applcation minus the sqlite and decimal packages for money safety.
Utilize SQL to buld 4 tables, products, deals, offerings, and cart

Abstractly:

Products are items that the shop might carry, these could be out of stock or on sale

Deals are abstract modifiers for products, like "Coupon", "buyXgetY", "Bundle", or "Retail".

Offerings tie together a Deal with one or more Product (in the case of Bundles) and is represented as an additional row in the the offerings table.

A cart contains just products and quantities. I didn't want to mess around with session cookies and all that and this is easily extentible by introducing a user row/table



# Drawbacks/Learnings

Attribution is hard, and working with multiple deals on a single product and minimzing the amount can be complex.
Bundles should be moved to a seperate relationship and have a BundleID/Group to make determining if all the bundled items
are in the cart more simple.


