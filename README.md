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

### Go get

Enable GOPROXY environment variable so that `go build` command could visit an accessible module server.

Linux/Max: `export GOPROXY=https://goproxy.io`

For windows, open your terminal `$Env:GOPROXY=https://goproxy.io`

### Front-end Assets

Run command: `cd client && npm install && npm run build`

### Run

`make build && make run`

You can use Windows Subsystem Linux to run Makefile. Or install [this tool](https://taskfile.dev/#/installation) to use the Taskfile.yml: `task build && task run`

## Deployment

The HTML files in `templates` are compiled to go binary using [go.rice](https://github.com/GeertJohan/go.rice). Install it to your machine and then run `rice embed-go`. It will generate `rice-box.go` file. After that you can run `make linux`.

## Licence Rules

Ideally a licence user should have a clean account without membership and this account is dedicated to b2b licence. However, reality does not always go that way. We should follow these rules to add b2b to existing subscription:

* For valid auto-renewal subscription, always deny granting licence to user account;
* Expired user is treated as a clean account;
* A valid B2B membership should be denied of any means of self-subscription.

Excluding the above cases, we are left with switching between alipay/wechat/b2b:

### Alipay/wechat to B2B

When a valid alipay/wechat subscription is trying to use b2b licence, turn current remaining subscription period to add-on and grant the licence to this user immediately.

### Licence renewal prior to expiration

When a B2B licence is renewed before it is expired, the corresponding membership is updated.

### B2B Expired

When a b2b membership expired and licence is no longer renewed:

* If it has addon, addon will be re-enabled, and the linked licence should be reset to available state;
* With or without addon, the user is free to make purchase as normal.

In the above 2 cases, to ensure data integrity, always **reset** the linked licence to clean state by removing the `assignee_id` field and `current_status` to `available`.

When a licence is renewed after expired, the licence user will be renewed if it is not touched after expiration. If the user, however, changed subscription via alipay/wechat either by a purchase out of its own pocket or previous addon, the licence should already have been reset and does not have a membership linked to it now. In such case we are safe to update the licence only.

There's special case under the above case if data integrity is broken: user's payment method is changed from b2b to alipay/wechat, but the licence still contain assignee_id pointing to this membership. In such case licence renewal should never change its linked user's membership. Shall we simply reset licence to clean state?