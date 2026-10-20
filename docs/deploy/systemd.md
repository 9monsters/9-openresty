# systemd 守护

## Hub Service

`/etc/systemd/system/openresty-admin.service`:

```ini
[Unit]
Description=OpenResty Admin Hub
After=network.target postgresql.service

[Service]
Type=simple
User=openresty
Environment="JWT_SIGNING_KEY=prod-secret-key"
Environment="DATABASE_URL=postgres://openresty:pwd@localhost:5432/openresty?sslmode=disable"
ExecStart=/opt/openresty-admin/bin/openresty-admin --host 0.0.0.0 --port 3000 --grpc-addr :9090
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Agent Service

`/etc/systemd/system/openresty-agent.service`:

```ini
[Unit]
Description=OpenResty Agent
After=network.target openresty.service

[Service]
Type=simple
User=openresty
ExecStart=/opt/openresty-admin/bin/openresty-agent \
  --hub-addr hub.internal:9090 \
  --metrics-port 9101 \
  --nginx-bin /usr/local/bin/openresty \
  --nginx-conf-dir /etc/nginx \
  --nginx-log-dir /var/log/nginx
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## 启用

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now openresty-admin
sudo systemctl enable --now openresty-agent

# 状态检查
sudo systemctl status openresty-admin
sudo journalctl -fu openresty-admin
```