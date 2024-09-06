# Gotel Reservation

## Start MongoDB as a container
```
docker run --name mongodb -d -p 27017:27017 mongodb/mongodb-community-server:6.0-ubi8
```

## Seed db
```
make seed
```

## Test units
```
make test
```

## TODO

- [x] Init
- [x] Set Makefile
- [x] Initialize DB
- [x] Set Users CRUD
- [x] Set User Validation
- [x] Test User API
- [x] Set Hotel API
- [x] Set JWT Auth
- [x] Set Booking API
- [x] Set Booking Validiation
- [x] Admin Auth
- [x] Set Error Handling Middleware
- [x] Clean Up