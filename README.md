Tirek
=====
_**Small and simple DNS load tester**_

Tirek is a lightweight utility for producing masive DNS load, this is done by spawning a set of
worker threads. Each one of those threads is responsible for producing requests at the specified rate.

**NOTE**: This is defiantly **NOT** a production tool and should not be used for benchmarking
 
Usage
=====

```sh
Usage of ./tirek:
  -cpus int
    	Sets the number of CPUs to be used during test (default _Number of CPU cores on your system_)
  -duration duration
    	Sets how long the attack should last (default 1m0s)
  -rate int
    	Sets the number of requests to attempt per second (default 50)
  -target string
    	Sets the DNS server to target (default "127.0.0.1:53")
  -workers int
    	Sets the number of worker threads to use (default 5)
```

Licence and Copyright
=====================

MIT License

Copyright (c) 2017, Liam Haworth.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.