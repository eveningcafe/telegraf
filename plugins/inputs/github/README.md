# GitHub Input Plugin

Gather repository information from [GitHub][] hosted repositories.

**Note:** Telegraf also contains the [webhook][] input which can be used as an
alternative method for collecting repository information.

### Configuration

```toml
[[inputs.github]]
  ## List of repositories to monitor
  repositories = ["influxdata/telegraf"]

  ## Github API access token.  Unauthenticated requests are limited to 60 per hour.
  # access_token = ""

  ## Timeout for HTTP requests.
  # http_timeout = "5s"
```

### Metrics

- github_repository
  - tags:
    - name - The repository name
    - owner - The owner of the repository
    - language - The primary language of the repository
    - license - The license set for the repository
  - fields:
    - stars (int)
    - forks (int)
    - open_issues (int)
    - size (int)

When the [internal][] input is enabled:

+ internal_github
  - tags:
    - access_token - An obfusticated reference to the configured access token or "Unauthenticated"
  - fields:
    - limit - How many requests you are limited to (per hour)
    - remaining - How many requests you have remaining (per hour)
    - blocks - How many requests have been blocked due to rate limit

### Example Output

```
github,full_name=influxdata/telegraf,name=telegraf,owner=influxdata,language=Go,license=MIT\ License stars=6401i,forks=2421i,open_issues=722i,size=22611i 1552651811000000000
internal_github,access_token=Unauthenticated rate_limit_remaining=59i,rate_limit_limit=60i,rate_limit_blocks=0i 1552653551000000000
```

[GitHub]: https://www.github.com
[internal]: /plugins/inputs/internal
[webhook]: /plugins/inputs/webhooks/github
