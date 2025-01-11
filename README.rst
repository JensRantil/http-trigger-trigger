====================
HTTP Trigger Trigger
====================
"HTTP Trigger Trigger" is an HTTP server application that take requests
and for every request it either triggers an HTTP request to some other
upstream destination, or it executes a shell command. The configuration
is stored in a simple INI file.

The difference between HTTP Trigger Trigger and many other HTTP proxies
(such as nginx and Apache) is that its backend requests are fully
independent of their equivalent frontend request (except timing,
obviously). This means, among other things, that all GET and POST
parameters are stripped off.

HTTP Trigger Trigger also supports rate limitting. This makes it very
useful to make sure that your backend infrastructure will not be
DDoS:ed/overloaded in case of a DDoS attack against HTTP Trigger
Trigger.

Why would I want to use this?
-----------------------------
The initial usecase was to securely expose to the web a build tool
trigger so that Github could trigger builds when new pushes arrives.

In general, this tool is useful when you'd like to expose a certain
internal API endpoint to an insecure network, and you'd like to do this
in a very secure fashion.

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
See ``setup.ini.example`` for an example configuration file containing
all possible key/values. The file uses standard INI format. Here's a
stripped down example::

    listen=:8080

    [/trigger/backend/http]
    url = http://upstream.server.example.com/other-trigger

    [/trigger/script]
    command = echo Hello World

    [/rate/limited/trigger]
    url = http://upstream2.server.example.com/another-trigger
    queue_size=2
    hit_delay=1s

For every request all endpoints will be matched sequentially against the
INI section. The first match will trigger a request to its corresponding
``url``. If ``listen`` is not defined, it will use port ``8080`` and
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

* Modify the default response content.

