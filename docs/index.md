---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/odbc.svg"
brand_color: "#000000"
display_name: "ODBC"
short_name: "odbc"
description: "Steampipe plugin to query data from ODBC data sources."
og_description: "Query ODBC files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/odbc-social-graphic.png"
---

# ODBC + Steampipe

Open Database Connectivity (ODBC) is a standard application programming interface (API) for accessing SQL databases.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

Query data from the `postgresql_names` table:

```sql
select name, _metadata from postgresql_names;
```

```
+------+-----------------------------------------------+
| name | _metadata                                     |
+------+-----------------------------------------------+
| jon  | {"connection_name":"odbc","dsn":"postgresql"} |
+------+-----------------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/odbc/tables)**

## Get started

### Install

Download and install the latest ODBC plugin:

```bash
steampipe plugin install odbc
```

### Credentials

Your system's ODBC configuration holds any required credentials for any databases you connect to. The plugin requires none, it just needs to know the data source names (DSNs).

### Configuration

Installing the latest ODBC plugin will create a config file (`~/.steampipe/config/odbc.spc`) with a single connection named `odbc`:

```hcl
connection "odbc" {
  plugin = "odbc"

  # fill the structure with DSN:table pairs.

  # data_sources = [
  #  "SQLite:names",
  #  "PostgreSQL:names"
  # ]

}  

```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-odbc
- Community: [Slack Channel](https://steampipe.io/community/join)
