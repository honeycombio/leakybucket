# leakybucket

leakybucket is a go library for implementing rate limits using the [leaky bucket
algorithm](https://en.wikipedia.org/wiki/Leaky_bucket)

The rate limit is expressed in terms of a throughput and burst limit - over time
the average throughput can't exceed the throughput limit, but in any given time
period, traffic may be allowed to exceed the throughput limit in bursts.

This type of rate limit is often seen in network - you might buy 1Mb of
bandwidth with a 100Mb burst limit. Over the course of time, your average
throughput must be below 1Mb, but if you need to temporarily spike above it,
it's fine.

## Other implementations

There are quite a few leaky bucket implementations out there, each is slightly
different. Some of the alternatives written in go:

* https://github.com/Clever/leakybucket - only allows you to set a throughput
  limit; it has no burst limit. has 3 different storage engines
* https://github.com/joncalhoun/drip - uses a goroutine for emptying the bucket

## Contributions and license

Features, bug fixes and other changes  are gladly accepted. Please open issues
or a pull request with your change.

All contributions will be released under the Apache License 2.0.
