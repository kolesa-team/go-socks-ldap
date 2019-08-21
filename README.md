# Basic Socks Proxy implementation with LDAP integration

# Quick start

```bash
docker run -d \
    -p 1080:1080 \
    -v $(pwd)/config.toml:/config.toml \
    kolesa/go-socks-ldap -config /config.toml -debug
```

# Features
 - Connecting to LDAP server (non TLS / TLS)
 - Local copy of user entries by filters and autoupdate from ldap

## Contributing:

Bug reports and pull requests are welcome!
