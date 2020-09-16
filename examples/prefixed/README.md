# Prefixed

This example shows how to execute a child process and prefix each line of output
with a string of your choosing. The string will be colored randomly using logic
defined by <https://github.com/brad-jones/goprefix/v2>. This is especially useful
for running many child processes concurrently and creating output similar to
tools like `docker-compose`.

This example also makes use of the functionality provided by
<https://github.com/brad-jones/goasync> to run the 2 ping
processes concurrently.

It also uses some functionality from <https://github.com/brad-jones/goerr>
to handle errors.

## Expected Output

_Depending on OS & network conditions of course_

```
ip2 |
ip2 | Pinging 127.0.0.2 with 32 bytes of data:
ip1 |
ip1 | Pinging 127.0.0.1 with 32 bytes of data:
ip1 | Reply from 127.0.0.1: bytes=32 time=11ms TTL=59
ip2 | Reply from 127.0.0.2: bytes=32 time=12ms TTL=59
ip2 | Reply from 127.0.0.2: bytes=32 time=10ms TTL=59
ip1 | Reply from 127.0.0.1: bytes=32 time=10ms TTL=59
ip1 | Reply from 127.0.0.1: bytes=32 time=11ms TTL=59
ip2 | Reply from 127.0.0.2: bytes=32 time=11ms TTL=59
ip2 | Reply from 127.0.0.2: bytes=32 time=11ms TTL=59
ip2 |
ip2 | Ping statistics for 127.0.0.2:
ip2 | Packets: Sent = 4, Received = 4, Lost = 0 (0% loss),
ip1 | Reply from 127.0.0.1: bytes=32 time=11ms TTL=59
ip1 |
ip1 | Ping statistics for 127.0.0.1:
ip1 | Packets: Sent = 4, Received = 4, Lost = 0 (0% loss),
ip1 | Approximate round trip times in milli-seconds:
ip1 | Minimum = 10ms, Maximum = 11ms, Average = 10ms
ip2 | Approximate round trip times in milli-seconds:
ip2 | Minimum = 10ms, Maximum = 12ms, Average = 11ms
```
