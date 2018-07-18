# XMPP Prekey Server for OTRv4 and forward

This project implements an XMPP Component (XEP-114) that allows XMPP server
providers to deploy an OTRv4 Prekey Server
(https://github.com/otrv4/otrv4-prekey-server) in an easy way.

In order to use this component, you also need to run the Raw prekey server (https://github.com/otrv4/otrng-prekey-server/tree/master/server/raw).

If you're using Prosody, you can add configuration options like these to add the server:

```
component_ports = { 42335 }
component_interface = "127.0.0.1"

Component "prekeys.example.com"
        component_secret = "changeme"
        name = "OTR Prekey Server"
```

You also need to add a DNS s2s entry for the domain you're using. 

For ejabberd you can add an equivalent ejabberd_service.

Assuming you have the raw server running on localhost with the default port you
can now start the xmpp server like this:

```
$ ./prekey-server-xmpp -xmpp-address 127.0.0.1 -xmpp-port 42335 -shared-secret "changeme" -name prekeys.example.com -fingerprint 123AAAAAAAAAAAAAAAAA123CDCDCDABABABABABA88797676556456421132425673535467575765AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
```

This component will now be discoverable as prekeys.example.com inside the XMPP
server - it will respond to disco#info and disco#items, and messages sent to it
will be forwarded to the Raw prekey server.
