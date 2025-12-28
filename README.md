# jt

`jt`(json_table) is a CLI tool that converts JSON data to MySQL `JSON_TABLE` expressions.

Paste JSON, get SQL. Drop into CTEs and subqueries.

## Why jt?

As a backend engineer at a startup, I often need to perform ad-hoc data queries directly against the database.
Whether it's a one-time data fix for a customer, data request from management, or handling requests that fall outside the scope of our usual application features—these situations come up regularly.

MySQL's `JSON_TABLE` function is incredibly useful for these tasks.
When you start with a predefined set of data (often in JSON format), treating it as a relational table makes querying and transforming straightforward.

But writing `JSON_TABLE` expressions by hand is tedious, that's where `jt` comes in.

## Usage

```bash
# From stdin
echo '[{"id": 1, "name": "Alice"}]' | jt

# From file
jt data.json
```

Output:
```sql
JSON_TABLE(
  '[{"id":1,"name":"Alice"}]',
  '$[*]' COLUMNS (
    id INT PATH '$.id',
    name VARCHAR(255) PATH '$.name'
  )
) AS jt
```

### Flags

- `-c`, `--collation` : Specify a collation for string columns (e.g., `utf8mb4_unicode_ci`).

```bash
jt -c utf8mb4_unicode_ci data2.json
```

Output:
```sql
JSON_TABLE(
  '[{"id":1,"name":"divingbeetle"},{"id":2,"name":"John Smith"}]',
  '$[*]' COLUMNS (
    id INT PATH '$.id',
    name VARCHAR(255) COLLATE utf8mb4_unicode_ci PATH '$.name'
  )
) AS jt
```

## Installation

```bash
go install github.com/divingbeetle/jt@latest
```

Or build from source:

```bash
git clone https://github.com/divingbeetle/jt.git
cd jt
go build -o jt
```

## Limitations

Project is in early development stage, and might not cover all edge cases.

- **Performance**: My use case is not performance critical, so the SQL generation might not be optimized.
- **Type mapping**: JSON value type to SQL data type conversion is not well-defined; complex scenarios may not be handled perfectly.
- **MySQL only**: Other DBMSs or dven different MySQL versions are not targeted and may not work as expected.

## License

MIT
