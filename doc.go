/*
Package leakybucket provides an in-memory implementation of the leaky bucket
algorithm for rate limiting.

Summary

The leaky bucket algorithm (wikilink) is often used for rate limits or to create
an even stream of traffic from a variable source. This package is only for rate
limits.

When configuring a leaky bucket, you specify the allowed throughput (volume and
time) and a burst limit.

Throughput limit is the rate at which the bucket drains. In other words, with a
full bucket, this is the rate of events that will leak out. The throughput limit
should be specified with two parameters, the DrainAmount is the number of
entries to leak, and the DrainPeriod is the frequency to drain the bucket. For
example, you would specify a throughput limit of 5 entries per second with a
DrainAmount of 5 and a DrainPeriod of time.Second

Burst limit is the number of entries you want to allow before starting to reject
new entries. This is the size of the bucket - as entries arrive the bucket
fills, and when it reaches capacity new entries will be rejected. Burst limit
should not be less than the throughput limit (or the throughput limit will never
be reached, as the bucket will reject entries when it hits capacity).

When the rate limit is reached and adding a new entry would exceed the burst
limit, a BucketOverflow error is returned.

Usage

Create one Bucket for each resource you wish to protect with a rate limit. When
a new entry arrives for that resource, attempt to Add() the event to the bucket.
If successful, Add() will return nil. If the rate limit has been exceeded, Add()
will return a BucketOverflow error.

Buckets have internal locking so Add()ing to a bucket is thread safe.
*/
package leakybucket
