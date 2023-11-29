# k8smultiarcher

## Description

## Configuration

| Environment Variable | Description |
| -------------------- | ----------- |
| CACHE_SIZE           | Sets the size of the cache. If not provided, a default size is used. |
| CACHE                | Determines the type of cache to use. Can be either 'inmemory' or 'redis'. If not provided or set to 'inmemory', an in-memory cache is used. |
| REDIS_ADDR           | Sets the address of the Redis server. Used when CACHE is set to 'redis'. If not provided, a default address is used. |
| HOST                 | Sets the host for the server. |
| PORT                 | Sets the port for the server. If not provided, the default is '8443' if TLS is enabled, '8080' otherwise. |
| TLS_ENABLED          | Determines whether TLS is enabled. If set to 'true', TLS is enabled. |
| CERT_PATH            | Sets the path to the TLS certificate. Used when TLS_ENABLED is set to 'true'. If not provided, the default is './certs/tls.crt'. |
| KEY_PATH             | Sets the path to the TLS key. Used when TLS_ENABLED is set to 'true'. If not provided, the default is './certs/tls.key'. |
