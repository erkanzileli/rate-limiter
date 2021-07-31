# rate-limiter

[![Coverage Status](https://coveralls.io/repos/github/erkanzileli/rate-limiter/badge.svg)](https://coveralls.io/github/erkanzileli/rate-limiter)

## Long story short

This is an app that you should put in front of your app to limiting incoming HTTP traffic. Works on L7 and acts as a
reverse proxy. I'm sure there are a lot of similar projects but this project is exactly for my need.

> Note: Putting in front of means that you should redirect the incoming traffic to this app instead of your app. This app will redirect the traffic to your app if it should.

## Motivation

At the end this app is an HTTP server. Acts as L7 single host reverse proxy.

## How to configure

I prefer telling these by using an example configuration. Here is an example configuration.

```yaml
app:
  port: 8081
  timeout: 3s
  hosts:
    - 12.34.56.78
    - 12.34.56.79
    - 12.34.56.80
server:
  addr: :8080
  readTimeout: 1000
  writeTimeout: 1000
cache:
  inMemory: true
  redis:
    addr: localhost:6379
    username: ""
    password: ""
    db: 0
algorithm:
  name: fixed-window-count
  options:
    # Fixed Window Count Options
    windowLengthInSeconds: 60
    # Another Algorithm Options..
defaultRuleScope: path # todo ?
rules:
  - scope: path
    pattern: GET /users/4/addresses
    windowLengthInSeconds: 30
    limit: 4
tracing:
  enabled: false
  provider: new-relic
  newRelic:
    appName: appName
    licenseKey: licenseKey
    distributedTracerEnabled: true
```

In the `app` object you fulfill your app's information.

| Name | Description | Required | Default |
|:---:|:---:|:---:|:---:|
|port | Your app port. You can specify this on each of your hosts  | yes | None |
|hosts | Your app instances. Can be `host` only or `host:port` | yes | None |
|timeout | Your app's timeout when redirecting | no | 1000ms |

In the `server` object you set up a server.

| Name | Description | Required | Default |
|:---:|:---:|:---:|:---:|
|addr | Address to listen  | yes | None |
|readTimeout | Server read timeout in millisecond format   | no | 1000ms |
|writeTimeout | Server read timeout  in millisecond format | no | 1000ms |

In the `cache` object you specify whether you want to use in-memory cache or Redis. If you use in-memory than the only
thing you have to do is set `cache.inMemory` field as `true`. However, if you use Redis than you only have to specify
the Redis settings.

In the `algorithm` object you specify which algorithm you want to use. Currently, only option is **fixed-window-count**.

| Name | Description | Required | Default |
|:---:|:---:|:---:|:---:|
|name | Algorithm's special name like enum  | no | fixed-window-count |
|options.windowLengthInSeconds | Valid for fixed-window-count and specifies the windows length.   | no | 60 |

In the `rules` array, you specify your rules. The pattern on the rule is basically a Regular Expression. _Rules are
processed sequentially._

| Name | Description | Required | Default |
|:---:|:---:|:---:|:---:|
|pattern | Regular expression to comparison pattern.   | yes | None |
|limit | Limit of the total requests   | yes | 0 |
|scope | Specifies the increment key. Can be `rule` or `path`. Rule pattern is used as increment key. Otherwise request's `METHOD PATH`   | no | `$defaultRuleScope` |
|windowLengthInSeconds | Custom period of just this rule   | no | `$algorithm.options.windowLengthInSeconds` |

In the `tracing` object, you specify a tracing method for this app. We support only NewRelic right now.

| Name | Description | Required | Default |
|:---:|:---:|:---:|:---:|
|enabled | Enabled the tracing   | no | false |
|provider | Your tracing provider. Available providers is: [`new-relic`]   | no | `new-relic` |
|newRelic.appName | App name on NewRelic   | yes | rate-limiter |
|newRelic.licenseKey | License key to access to NewRelic   | yes | None |
|newRelic.distributedTracerEnabled | Specifies whether Distributed Tracing is enabled or not   | no | false |

## homeless texts

Your regexp is compared with this format `METHOD PATH`, for example, you get a request like `DELETE /users?name=John`
but this is compared with your regexp as this format `DELETE /users`.

Let's make another example. Assume that we have already some rules as defined above. Here is an example `curl` and the
processing of this.

```shell
$ curl -X GET localhost:8000/users/4/addresses?country=tr
```

You can do this request 4 times on each 30 seconds.

## Todo

- [ ] Change config structure to give `app: { addr connectTimeout readTimeout }`
- [ ] Kubernetes Operator
- [ ] OpenTelemetry Support
- [ ] Rich algorithm support
