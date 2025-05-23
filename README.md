i2p transport for apt
=====================

This is a simple transport for downloading debian packages from a repository
over i2p. It uses the SAM bridge and is pretty much just a modified version of
[diocles/apt-tranport-http-golang](https://github.com/diocles/apt-transport-http-golang),
a plain HTTP transport for apt, and by modified I mean like, a couple dozen
lines. This way, it uses an ephemeral destination instead of the HTTP proxy
which could be associated with your other traffic.

Besides that, I think most would agree that it is simpler to use an apt
transport to detect when a package should be retrieved from an i2p service.
Especially in cases where the user is mixing packages from Tor, I2P, and
Clearnet sources, this process can become confusing and involve configuring
multiple applications along with apt. Instead, apt-transport-i2p works with
other apt transports like apt-transport-tor and requires no configuration on
the vast majority of systems.

**WARNING: I'm introducing a change which will break existing installations**
**of this software. When you update, make sure to update your sources.list to**
**use the new namespace. Change all occurrences of "i2p://" to "i2psam://"**

*This application does not have outproxy support. If you need anonymous access*
*to Debian packages over the clearnet, consider apt-transport-tor or the*
*forthcoming apt-transport-i2phttp, which makes use of the HTTP proxy.*

This application is *beta* software, more-or-less. I don't plan to change it,
but I will respond to any issues that you submit. The branch tunnel-multiplexer
is pre-alpha, and will become the second version of apt-transport-i2p.

To install it:
--------------

If you're using Java I2P, you'll need to enable your SAM Application Bridge to
be able to use the SAM API. In the browser you use to visit the router console,
open [127.0.0.1:7657/configclient](http://127.0.0.1:7657/configclient) and check
the box to run the SAM Application bridge at startup, then save your client
configuration.

![Enable the SAM bridge](configclient.png)

Then click the ">" arrow to the right of the check box to
manually start the SAM bridge on the running router.

Once you've enabled the just:

        make build && sudo make install

to install bin/apt-transport-i2p to /usr/lib/apt/methods/i2psam

To use it:
---------

Adding this to your sources.list.d will configure apt to seek updates to
ppa.launchpad.net/i2p-maintainers from a caching proxy at the b32 address:
```h2knzawve56vtiimbdsl74bmbuw7xr65xhgrdjtjnbfxxw4hsqlq.b32.i2p```

        deb i2psam://h2knzawve56vtiimbdsl74bmbuw7xr65xhgrdjtjnbfxxw4hsqlq.b32.i2p/ppa.launchpad.net/i2p-maintainers/i2p/ubuntu bionic main
        deb-src i2psam://h2knzawve56vtiimbdsl74bmbuw7xr65xhgrdjtjnbfxxw4hsqlq.b32.i2p/ppa.launchpad.net/i2p-maintainers/i2p/ubuntu bionic main

