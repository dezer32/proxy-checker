# Proxy checker

## Usage

### Cmd line help

```shell
Usage of ./proxy-checker:
  -i string
    	Path to file with proxies. (default "proxies.json")
  -o string
    	Path to file with checked json. (default "proxies.checked.1673272869.json")
```

### Proxy file format

```json
[
  {
    "ip": "185.43.249.132",
    "port": 39316,
    "protocol": "socks4"
  }
]
```