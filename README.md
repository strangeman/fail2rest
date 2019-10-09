# fail2rest

## Overview
fail2rest is a small REST server that aims to allow full administration of a fail2ban server via HTTP

fail2rest can be run using Docker. It must mount the fail2ban socket in `/var/run/fail2ban` and the database in `/var/lib/fail2ban`.

## Endpoints

| Endpoint         | Method | Sample Body   | Description                   |
| ---------------- | ------- | ------------- | ----------------------------- |
| /global/ping     | GET | No data       | Verify fail2ban can be pinged |
| /global/bans     | GET | No data       | List all banned IPs           |
| /global/status   | GET | No data       | Get status of fail2ban        |
| /jail/{jail}/      | GET | No data       | Get all data about a jail       |
| /jail/{jail}/ban   | POST | `{"ip": "127.0.0.1"}` | Ban an IP in the jail   |
| /jail/{jail}/unban | POST | `{"ip": "127.0.0.1"}` | Unban an IP in the jail |
| /jail/{jail}/failregex | POST | `{"fail_regex": "sample regex"}` | Add a fail regex to the jail |
| /jail/{jail}/failregex | DELETE | `{"fail_regex": "sample regex"}` | Remove a fail regex from the jail |

## License
The MIT License (MIT)

Copyright (c) 2014 Sean DuBois

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
