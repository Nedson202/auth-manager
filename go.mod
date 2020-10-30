module github.com/nedson202/user-service

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20201022194115-1af099fb3eca
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/pretty v0.1.0
	github.com/lib/pq v1.8.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/rs/cors v1.7.0
	github.com/subosito/gotenv v1.2.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	google.golang.org/appengine v1.6.7
)

go 1.13

replace github.com/nedson202/user-service => ./../user-service
