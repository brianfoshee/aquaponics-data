Go app to collect and read back sensor data.

Current Postgres schema:
```sql
CREATE TABLE readings (
  ph numeric(18,2),
  tds numeric(18,2),
  temperature numeric(18,2),
  created_at timestamp
);
```

The app requires two ENV variables to be set before it'll run:
```
PORT=5000 DATABASE_URL=postgres://user@localhost/database go run main.go
```

If any dependencies are added, use [Godep] [godep] to package them into the repo:
```
godep save
```

[godep]: https://github.com/tools/godep
