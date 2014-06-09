# Git Hook Proxy

## Motivation

To help integrate Gitlab post receive hook and do further CI with something like Jenkins.

Gitlab sends a webhook as an 'application/json' post, with the [JSON as part of the post body](http://grab.by/qrKw).
This makes it harder for Jenkins to read, as it usually expects parameters.
This is slightly different from Github, which posts it as a ['payload' parameter](https://help.github.com/articles/post-receive-hooks).  

Git hook proxy takes the Gitlab web hook and translates it into something more easily workable by Jenkins.

## To Run proxy

    go build proxy.go

This will generate a 'proxy' executable.  

To run the proxy:

    ./proxy -listen <listen address>

For Example:

    ./proxy -listen 127.0.0.1:9999

## To use proxy

Add this to your Gitlab webook: 

    http://[proxy_listen_url]?url=[target_url]

Make sure to specify 'url' parameter to tell the proxy where to forward requests to.
  
The proxy will take the webhook request, and translate it to a request to the target_url in the format of:

- payload: JSON body
- START: Start commit hash
- END: End commit hash
- REFNAME: Ref name

