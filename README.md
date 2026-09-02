# qzone-history

一个使用 Go 编写的 QQ 空间历史活动备份与重建工具。程序通过 QQ 二维码登录，读取当前账号可访问的活动记录，保存到本地 SQLite 数据库，并导出 JSON 备份。

> 本项目仅供学习、研究和个人数据备份使用。请遵守腾讯服务条款与适用法律，不要收集或传播无权处理的数据。

## 当前能力

- 在仅监听 `127.0.0.1` 的本地网页中展示 QQ 登录二维码
- 缓存本地登录信息，并在复用前验证登录状态
- 分页获取 QQ 空间活动记录并保存到 SQLite
- 从活动记录重建说说与留言板数据
- 导出说说、留言和好友数据为 `<QQ号>_export.json`

Excel 和 HTML 导出接口尚未实现；当前稳定的自动导出格式只有 JSON。

## 环境要求

- Go 1.25.2 或更高版本
- Windows、macOS 或带桌面浏览器的 Linux
- 手机 QQ（用于扫码登录）

## 快速开始

```bash
git clone https://github.com/MarcusDog/qzone-history.git
cd qzone-history
go mod download
go run ./cmd
```

启动后程序会自动打开本地二维码页面。扫码完成后不要立即关闭终端；程序会继续抓取、重建并导出数据。生成的数据库和导出文件默认不会被 Git 跟踪。

## 配置

默认配置位于 `config/config.yaml`。当前运行入口只支持 SQLite；数据库路径可通过配置中的 `database.dbName` 修改。

配置也可用环境变量覆盖，变量名以 `QZONE_HISTORY_` 开头，并以 `_` 代替层级分隔符，例如：

```bash
QZONE_HISTORY_DATABASE_DBNAME=./backup.db go run ./cmd
```

程序日志只输出数据库类型和文件路径，不会打印数据库密码等完整配置。

## 验证

```bash
go test ./...
go vet ./...
go test -race ./...
```

## 项目结构

- `cmd`：应用入口与依赖装配
- `internal/domain`：实体、仓库与用例接口
- `internal/usecase`：业务流程
- `internal/infrastructure`：QQ 空间客户端、配置与持久化
- `internal/delivery/app`：命令行应用编排
- `pkg`：数据库、二维码页面和解析工具

## 数据与安全说明

- 登录 Cookie 和历史数据会保存在本机数据库中，请妥善保护该文件。
- JSON 导出文件使用仅当前用户可读写的权限创建（操作系统支持时）。
- 二维码网页仅绑定回环地址，并设置禁用缓存和基础浏览器安全响应头。
- QQ 网页接口可能随时调整；若接口返回异常，程序会明确报错，不会把失败响应当成有效数据。

## License

[MIT](LICENSE)
