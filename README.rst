====================
HTTP Trigger Trigger
====================

"HTTP Trigger Trigger" is an HTTP server application that take requests
and for every request it triggers an HTTP request to some other
upstream destination. The configuration is stored in a simple INI file.

The difference between HTTP Trigger Trigger and many other HTTP proxies
(such as nginx and Apache) is that its backend requests are fully
independent of their equivalent frontend request (except timing,
obviously). This means, among other things, that all GET and POST
parameters are stripped off.

Why would I want to use this?
-----------------------------
I've been using this to securily expose a build tool trigger to the web
so that Github could trigger builds when new pushes arrives.

In general, this tool is useful when you'd like to expose a certain
internal API endpoint to the web, but you'd like to do this in a very
secure fashion.

Installing/developing
----------
The easiest way to install is to `download a precompiled release`_. You
can also build the software yourself using a ``$GOPATH``/`Go
workspace`_::

    $ mkdir -p ~/src/http-trigger-trigger
    $ export GOPATH=~/src/http-trigger-trigger
    $ cd $GOPATH
    $ mkdir -p src/github.com/JensRantil
    $ cd src/github.com/JensRantil
    $ git clone git@github.com:JensRantil/http-trigger-trigger.git
    $ cd http-trigger-trigger
    $ go get
    $ go build

The binary built is ``http-trigger-trigger``.

.. _download a precompiled release: https://github.com/JensRantil/http-trigger-trigger/releases
.. _Go workspace: http://golang.org/doc/code.html

Configuration
-------------
The configuration consists of a single INI file::

    port=8080

    [/trigger/first]
    url = http://upstream.server.example.com/other-trigger

    [/trigger/second]
    url = http://upstream2.server.example.com/another-trigger

For every request all endpoints will be matched sequentially against the
INI section. The first match will trigger a request to its corresponding
``url``. If ``port`` is not defined, it will use port ``8080`` and
listen on all interfaces.

Running
-------
Create your configuration file, and start using::

    TRIGGER_TRIGGER_CONFIG=./config.ini ./http-trigger-trigger

Future improvements
-------------------
* Allow certain GET or POST parameters to pass through (and be renamed).
* Do input matching based on other things than path:
  * HTTP method
  * Host header.

