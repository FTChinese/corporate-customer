# FT Academy

This is the server-side implementation backing the [Reader](https://github.com/FTChinese/reader-react) and [B2B](https://github.com/FTChinese/b2b-react) apps. This app relies on the [API](https://github.com/FTChinese/subscription-api) to run since it forwards most reader-related requests to the API.

## Development

To build the binary, `./build/api.toml` file must exist. This file will be embedded into the resulting binary using go 1.16 embed package. Make sure your have go above 1.16.

### Go get

Enable GOPROXY environment variable so that `go build` command could visit an accessible module server.

Linux/Max: `export GOPROXY=https://goproxy.io`

For windows, open your terminal `$Env:GOPROXY=https://goproxy.io`

### Run

`make build && make run`

You can use Windows Subsystem Linux to run Makefile. Or install [this tool](https://taskfile.dev/#/installation) to use the Taskfile.yml: `task build && task run`

## Test Payment

The payment is divided into live/test mode. They are determined differently depending on the payment method chosen. Those modes, however, only applicable to server in production mode. For local development, you are always in testing mode.

### One-off payment

The mode is solely determined by user's logged in account. When you are using a testing account issued by Superyard, you are in test mode; otherwise in live.

### Stripe

Stripe payment is determined by the publishable key the server returned to client. Stripe.js requires it to be initialized upon page loading, outside any React components, and should never be re-created afterwards. This literally restrict you from getting publishable key from backend dynamically. For production, the only solution could be setting up different backends for live/test mode respectively.

My current solution is to give the app a command-line option `-livemode=<true|false>`. When `true`, the server given client Stripe's live publishable key; otherwise test key. Default value is `true` to keep backward compatible since initially the server's Supervisor is not configured with this option.

This way you could test Stripe locally by launching this app `ftacademy -production=false livemode=false`. If you want to use online API locally, run `ftcademy -produciton=true -livemode=false`.

## B2B Licence Rules

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

### Priorities of various payment source

From highest to lowest:

* IAP/Stripe
* B2B
* Alipay/Wechat

## API

The API consists of 2 parts based on request validation methods：

* Handle AJAX request using Json Web Token;
* Handle restful request using authorization key.

### AJAX

The AJAX part is divided into multiple sections based on features.

#### Paywall

This section is publicly available since they must be accessed unconditionally.

* GET `/api/paywall` Output paywall data.
* GET `/api/paywall/stripe/prices`
* GET `/api/paywall/stripe/prices/:id`
* GET `/api/paywall/stripe/publishable-key`

#### Reader Section

This section simply forwards requests for subscription-api.

* GET `/api/reader/auth/email/exists` Check if an email is signed up.
* POSt `/api/reader/auth/email/login` Login using email.
* POST `/api/reader/auth/email/signup` Create a new email account. It also goes here when a new mobile is trying to create a new email account and link to it.
* POST `/api/reader/auth/email/verification/:token` Verify email.
* PUT `/api/reader/auth/mobile/verification` Send user a SMS.
* POST `/api/reader/auth/mobile/verication` Verify the SMS in the above step.
* POST `/api/reader/auth/mobile/link` A new mobile links to existing email account
* POST `/api/reader/auth/mobile/signup` Create a new mobile-only account.
* POST `/api/reader/password-reset` Reset password.
* POST `/api/reader/password-reset/letter` Send a password-reset letter to user.
* GET `/api/reader/password-reset/tokens/:token` Verify a password reset letter.
* GET `/api/reader/account` Get user account, either ftc or wechat.
* PATCH `/api/reader/account/email` Change email for ftc account.
* POST `/api/reader/account/email/request-verification` Re-send a verification email.
* PATCH `/api/reader/account/name` Change username
* PATCH `/api/reader/account/password` Change password
* PATCH `/api/reader/account/mobile` Switch mobile
* PUT `/api/reader/account/mobile/verification` Send an SMS before permitting mobile switch.
* GET `/api/reader/account/address` Load address.
* PATCH `/api/reader/account/address` Update address
* GET `/api/reader/account/profile` Load profile
* PATCH `/api/reader/account/profile` Update profile
* POST `/api/reader/account/wx/signup` A wechat-login user creates a new email account and link to it.
* POST `/api/reader/account/wx/link` A wechat-login user links to existing email account.
* POST `/api/reader/account/wx/unlink` A wechat-login user, with email account linked, unlinks the email account.

#### B2B Section

* POST /api/b2b/auth/login
* POST /api/b2b/auth/signup
* GET /api/b2b/auth/verify/:token
* POST /api/b2b/auth/password-reset
* POST /api/b2b/auth/password-reset/letter
* GET /api/b2b/auth/password-reset/token/:token
* GET /api/b2b/account/jwt
* POST /api/b2b/account/request-verification
* PATCH /api/b2b/account/display-name
* PATCH /api/b2b/account/password
* GET /api/b2b/team
* POST /api/b2b/team
* PATCH /api/b2b/team
* GET /api/b2b/search/membership?email=<string>
* GET /api/b2b/orders
* POST /api/b2b/orders
* GET /api/b2b/orders/:id
* GET /api/b2b/licences
* GET /api/b2b/licences/:id
* POST /api/b2b/licences/:id/revoke
* GET /api/b2b/invitations
* POST /api/b2b/invitations
* POST /api/b2b/invitations/:id/revoke
* GET /api/b2b/licence/invitation/verification/:token
* POST /api/b2b/licence/grant

### Restful API

Used by another backend app "superyard" to access B2B data.

* GET /api/cms/profile/:id
* GET /api/cms/teams/:id
* GET /api/cms/orders/:id
* POST /api/cms/orders/:id