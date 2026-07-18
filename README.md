# Cursor API Bridge

一个 macOS 桌面工具，让 Cursor Agent 使用自定义的 OpenAI-compatible API。

## 功能

- 图形界面配置 API URL 和 API Key
- 自动读取 `/v1/models` 模型列表
- 支持多模型选择和推理强度
- 本地 Cursor Agent 协议适配
- SSE 流式响应与工具调用
- 每台电脑独立生成本地 CA
- 启停时自动应用或恢复 Cursor 代理配置
- 不包含任何默认 API URL、Key 或用户凭据

## 安全说明

程序通过本地 HTTPS 代理连接 Cursor，因此首次启动服务时需要把应用生成的本地 CA 加入 macOS 登录钥匙串。CA 私钥只保存在用户本机：

```text
~/.liufugui-cursor-bridge/data/ca.key
```

停止服务会清理 Cursor 代理设置。卸载前建议先在软件中停止服务。

## 构建

需要 Go、Node.js、Protocol Buffers 和 Wails v3。

```bash
cd frontend
npm install
npm run build
cd ..
go build -o CursorAgentBridge .
```

## 开源许可与来源

本项目基于 MIT 许可的软件进行二次开发，并保留原始版权声明。界面、独立 CA、模型自动发现、默认隐私配置和 macOS 打包流程经过重新设计。

详见 [LICENSE](LICENSE)。
