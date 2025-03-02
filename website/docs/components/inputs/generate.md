---
title: generate
type: input
status: stable
categories: ["Utility"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the contents of:
     lib/input/generate.go
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


Generates messages at a given interval using a [Bloblang](/docs/guides/bloblang/about)
mapping executed without a context. This allows you to generate messages for
testing your pipeline configs.

Introduced in version 3.40.0.

```yaml
# Config fields, showing default values
input:
  label: ""
  generate:
    mapping: ""
    interval: 1s
    count: 0
```

## Fields

### `mapping`

A [bloblang](/docs/guides/bloblang/about) mapping to use for generating messages.


Type: `string`  
Default: `""`  

```yaml
# Examples

mapping: root = "hello world"

mapping: root = {"test":"message","id":uuid_v4()}
```

### `interval`

The time interval at which messages should be generated, expressed either as a duration string or as a cron expression. If set to an empty string messages will be generated as fast as downstream services can process them. Cron expressions can specify a timezone by prefixing the expression with `TZ=<location name>`, where the location name corresponds to a file within the IANA Time Zone database.


Type: `string`  
Default: `"1s"`  

```yaml
# Examples

interval: 5s

interval: 1m

interval: 1h

interval: '@every 1s'

interval: 0,30 */2 * * * *

interval: TZ=Europe/London 30 3-6,20-23 * * *
```

### `count`

An optional number of messages to generate, if set above 0 the specified number of messages is generated and then the input will shut down.


Type: `int`  
Default: `0`  

## Examples

<Tabs defaultValue="Cron Scheduled Processing" values={[
{ label: 'Cron Scheduled Processing', value: 'Cron Scheduled Processing', },
{ label: 'Generate 100 Rows', value: 'Generate 100 Rows', },
]}>

<TabItem value="Cron Scheduled Processing">

A common use case for the generate input is to trigger processors on a schedule so that the processors themselves can behave similarly to an input. The following configuration reads rows from a PostgreSQL table every 5 minutes.

```yaml
input:
  generate:
    interval: '@every 5m'
    mapping: 'root = {}'
  processors:
    - sql_select:
        driver: postgres
        dsn: postgres://foouser:foopass@localhost:5432/testdb?sslmode=disable
        table: foo
        columns: [ "*" ]
```

</TabItem>
<TabItem value="Generate 100 Rows">

The generate input can be used as a convenient way to generate test data. The following example generates 100 rows of structured data by setting an explicit count. The interval field is set to empty, which means data is generated as fast as the downstream components can consume it.

```yaml
input:
  generate:
    count: 100
    interval: ""
    mapping: |
      root = if random_int() % 2 == 0 {
        {
          "type": "foo",
          "foo": "is yummy"
        }
      } else {
        {
          "type": "bar",
          "bar": "is gross"
        }
      }
```

</TabItem>
</Tabs>


