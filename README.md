Rubens
======

Rubens is a simple Go server to manipulate images (resize) built on top of `resize <https://github.com/nfnt/resize>`

The server as a proxy of your storage engine (Amazon S3) and serves from a redis cache. Configuration and image sizes can be specified via config.json


Installation
============

Build it
--------

1. Make sure you have a Go language compiler >= 1.3 (required) and git installed.
2. Ensure your GOPATH is properly set. Eg: If your workspace is $HOME/work then GOPATH could be $HOME/work:$HOME/work/github.com/itachi3/rubens
3. Download it:
::

    go get github.com/itachi3/rubens

You now have rubens and the libraries it uses for image handling.
4. Build using 
::
    go build main.go

You now have the rubens executable ready to serve

Configuration
=============

Configuration should be stored in a readable file and in JSON format.

``config.json``

.. code-block:: json

    {
    "logs" : {
    "errorLog" : "/var/log/go/error.log",
    "accessLog" : "/var/log/go/access.log"
    },
    "dataStores" : {
        "redis" : {
            "protocol" : "tcp",
            "port" : "6379"
        },
        "amazonS3" : {
            "region" : "us-east-1",
            "bucketName" : "image.agentdesks.com"
        }
    },
    "image" : {
    "width" : ["900", "380", "240"],
    "height" : ["1280", "640", "340"]
    }
    }
