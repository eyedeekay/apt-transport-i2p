i2p transport for apt
=====================

This is a simple transport for downloading debian packages from a repository
over i2p. It uses the SAM bridge and is pretty much just a modified version of
[diocles/apt-tranport-http-golang](https://github.com/diocles/apt-transport-http-golang),
a plain HTTP transport for apt, and by modified I mean like, a couple dozen
lines. This way, it uses an ephemeral destination instead of the HTTP proxy
which could be associated with your other traffic.

Also I just think it's easier.

To install it:
--------------

just:

        make build && sudo make install

to install bin/apt-transport-i2p to /usr/lib/apt/methods/i2p

To use it:
---------

To add an eepSite to your sources.list, for example:

        deb i2p://http://wnhxwrq4fkn3cov6bnqsdaniubeo3625rmsm53yaz336bxvtiqeq.b32.i2p/deb-pkg rolling main

