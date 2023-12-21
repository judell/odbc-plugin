# Table: {odbc_dsn_tablename}

Query data from ODBC sources. A table is automatically created to represent each DSN:tablename pair found in the configured `data_sources`.

For example, on Linux the `/etc/odbc.ini` file might define the DSNs (data source names) `SQLite` and `Postgres`.

```
$ cat /etc/odbc.ini
[SQLite]
Driver=/usr/local/lib/libsqlite3odbc.so
Database=/home/jon/sqlite/test.db

[PostgreSQL]
Driver=/usr/lib/x86_64-linux-gnu/odbc/psqlodbcw.so
Database=postgres
Servername=localhost
Port=5432
Protocol=11
ReadOnly=Yes
User=postgres
```

If a SQLite db has a table `names` and a Postgres db also has a table `names`, then this setup binds them to their respective identifiers using `DSN:table` format.


connection "odbc" {
  plugin = "odbc"

  data_sources = [
    "SQLite:names",
    "PostgreSQL:names"
  ]

}


### Inspect the table structure

Assuming your connection is called `odbc` (the default), list all tables with:

```bash
.inspect odbc
```

```
+------------------+------------------+
| table            | description      |
+------------------+------------------+
| postgresql_names | PostgreSQL:names |
| sqlite_names     | SQLite:names     |
+------------------+------------------+
```

To get details for a specific table, inspect it by name:

```bash
.inspect sqlite_names
```

```
+-----------+--------+-------------------------------------------------------+
| column    | type   | description                                           |
+-----------+--------+-------------------------------------------------------+
| _ctx      | jsonb  | Steampipe context in JSON form, e.g. connection_name. |
| _metadata | jsonb  | Metadata related to the data source                   |
| name      | text   | SQLite name                                           |
| number    | bigint | SQLite number                                         |
+-----------+--------+-------------------------------------------------------+
```

### Count names

```sql
select count(*) from postgresql_names;
```

### Find name #1

```sql
select name from sqlite_names where number = 1;
```




