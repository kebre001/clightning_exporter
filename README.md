# clightning_exporter

### Arguments

```
./lightning_exporter [PATH TO lightning-cli]
```


## Systemd service

```
[Unit]
Description=Prometheus c-lightning exporter
After=syslog.target network.target

[Service]
User=bitcoin
Group=bitcoin
ExecStart=[PATH TO EXPORTER BINARY] [PATH TO LIGHTNING-CLI]
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
