![image](https://hub.steampipe.io/images/plugins/turbot/odbc-social-graphic.png)

# ODBC Plugin for Steampipe

Use SQL to query data from ODBC data sources.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/odbc)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/odbc/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-odbc/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install odbc
```

[Configure](https://hub.steampipe.io/plugins/turbot/odbc#configuration) the data source names that you'll query. 

Run steampipe:

```shell
steampipe query
```

Count names:

```sql
select count(*) from postgresql_names;
```

Find a name by number:

```sql
select name from sqlite_names where number = 1;
```


## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-odbc.git
cd steampipe-plugin-odbc
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/odbc.spc
```

Try it!

```
steampipe query
> .inspect odbc
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-odbc/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [ODBC Plugin](https://github.com/turbot/steampipe-plugin-odbc/labels/help%20wanted)
