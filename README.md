sendit
======

A simple utility to make HTTP POST or GET requests.

Overview
--------

This utility is like `curl` and with that being said you should
use that instead.  I wrote sendit because I wanted an excuse to
write a program and learn more about Go.  Sendit skips security
and is fairly simple and limited in functuality.

Usage
-----

Make a GET request and print the body of the response to stdout:

```$ sendit -url 'https://www.example.com'```

POST some JSON data to some url (you can send any data, it doesn't
have to be JSON):

```$ sendit -url 'https://www.example.com/someapi' -data '{"foo": 13, "bar": 17}'```

POST the contents of a file to some url:

```$ sendit -url 'https://www.example.com/otherapi' -file boring.xml```

POST the contents of a file and some arbitrary data to some url:

```$ sendit -url 'https://www.example.com/megaapi' -file boring.json -data '[1,2,3]'```

In this last example two seperate POST requests are made: one with the
JSON file and the second with JSON array.

License
-------

Zlib
