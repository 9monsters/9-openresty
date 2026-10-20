# gRPC 协议参考

## 服务定义

Proto 文件：`proto/openresty/admin/v1/services.proto`

生成代码：`gen/go/openresty/admin/v1/`

| Service | RPC | 方向 | 说明 |
|---------|-----|------|------|
| `AgentRegistration` | `Register` | Agent → Hub | 注册新 Agent |
| `AgentRegistration` | `Heartbeat` | Agent → Hub | 心跳上报 |
| `ConfigDeploy` | `Push` | Hub → Agent | 推送配置变更 |
| `ConfigDeploy` | `Validate` | Hub → Agent | 验证配置 |
| `ConfigDeploy` | `Reload` | Hub → Agent | Reload Nginx |
| `ConfigDeploy` | `Rollback` | Hub → Agent | 回滚配置 |
| `StatsCollector` | `StreamStats` | Agent → Hub | 流式统计上报 |
| `TerminalProxy` | `Shell` | 双向流 | 终端代理 |
| `LogReader` | `ReadLog` | Hub → Agent | 远程日志 |

## PushConfig 消息

```protobuf
message PushConfigRequest {
  string release_id = 1;
  repeated ConfigChange changes = 2;
}

message ConfigChange {
  string path = 1;
  string content = 2;
  bool delete = 3;
  string checksum_sha256 = 4;
}
```

## 注册消息

```protobuf
message RegisterRequest {
  string hostname = 1;
  string nginx_version = 2;
  bool is_openresty = 3;
  map<string, string> labels = 4;
  int32 metrics_port = 5;
}

message RegisterResponse {
  string node_id = 1;
  string session_token = 2;
}
```