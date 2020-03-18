## Development

Make sure you have a configuration file under you home directory: `~/config/api.toml`. It should have an entry for DB connection:

```toml
[mysql]
    [mysql.dev]
    host = "" # string, you db ip
    port = 123 # number, you db port
    user = "" # string, your db user name
    pass = "" # string, your db password
```

### Front-end Assets

Run command: `cd client && npm install && npm run build`

### Run

`make build && make run`

You can use Windows Subsystem Linux to run Makefile. Or install [this tool](https://taskfile.dev/#/installation) to use the Taskfile.yml: `task build && task run`

## Deployment

The HTML files in `templates` are compiled to go binary using [go.rice](https://github.com/GeertJohan/go.rice). Install it to your machine and then run `rice embed-go`. It will generate `rice-box.go` file. After that you can run `make linux`.