## Go app to collect and read back sensor data.

### DB Setup
Current Postgres schema is in db.

To get the database setup locally (dev or test environment), ensure that Goose
is installed locally and then run its setup:
```
$ go get bitbucket.org/liamstask/goose/cmd/goose
$ DATABASE_URL="postgres://user@localhost/database?sslmode=disable" goose --env development up
```

See [Goose][goose] docs for more info on rolling back migrations.

To run migrations on Heroku, first push them up and then run:
```
$ heroku run goose --env production up
```

### Setting up dependencies
To update all of the latest dependencies run
```
$ godep restore
```

If any dependencies are added, use [Godep] [godep] to package them into the repo:
```
$ godep save
```
* For updating Goose, be sure to temporarily remove the line ``// +build heroku`` from `install_goose.go` 
so that godep will pick up the package.

## Running the app
The app requires two ENV variables to be set before it'll run:
```
$ PORT=5000 DATABASE_URL="postgres://user@localhost/database?sslmode=disable" go run main.go
```



[godep]: https://github.com/tools/godep
[goose]: https://bitbucket.org/liamstask/goose
