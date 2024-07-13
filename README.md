run postgresql server in the background
```bash
sudo service postgresql start
```

connect to psql
```bash
sudo -u postgres psql -d blogator
```


protocol://username:password@host:port/database
```bash
postgres://postgres:postgres@localhost:5432/blogator
```

goose migration
```bash
// up
goose -dir=./sql/schema postgres "postgres://username:password@localhost:5432/your_database" up

// down
goose -dir=./sql/schema postgres "postgres://username:password@localhost:5432/your_database" down
```