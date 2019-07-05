module github.com/Nedson202/user-service

go 1.12

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/context v1.1.1
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/lib/pq v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/nedson202/user-service v0.0.0-00010101000000-000000000000
	github.com/subosito/gotenv v1.1.1
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
)

replace github.com/nedson202/user-service => ./../user-service
