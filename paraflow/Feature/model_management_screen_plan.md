# Model Management
模型配置管理中心,提供AI模型的注册、配置、性能监控和API密钥管理能力,支持多种模型提供商。

Layout Hierarchy:
- Header (Full-width):
  - Top Navigation Bar
- Content Container (Positioned below the header):
  - Left Sidebar
  - Main Content Area

## Top Navigation Bar
- Logo和平台名称
- 全局搜索框
- 快速创建按钮
- 通知中心图标
- 用户头像和设置入口

## Left Sidebar
- Dashboard导航
- Agent Management导航
- Model Management导航(当前选中)
- Orchestration Canvas导航
- Schedule Tasks导航
- Execution History导航
- Chat Interface导航
- 系统设置入口

## Main Content Area

### 页面标题与操作栏
- 页面标题"Model Management"
- 添加新模型按钮(主操作按钮)
- 批量操作按钮(导入配置/导出配置)
- 视图切换(列表视图/卡片视图)

### 筛选与搜索栏
- 搜索框(按模型名称/提供商搜索)
- 筛选条件:
  - 模型提供商(全部/OpenAI/Anthropic/本地模型/其他)
  - 模型状态(全部/已激活/配置中/不可用)
  - 模型类型(全部/对话/嵌入/图像生成)
- 排序选项(添加时间/使用频率/名称)

### 模型列表区域
- 模型卡片(多个):
  - 模型名称和版本
  - 提供商logo和名称
  - 模型类型标签
  - 状态指示器(可用/不可用)
  - 使用统计:
    - 本月调用次数
    - 平均响应时间
    - Token消耗量
  - 关联的Agent数量
  - 快速操作按钮:
    - 配置按钮
    - 测试按钮
    - 更多操作菜单(禁用/删除/导出配置)
  - 最后更新时间

### 模型详情侧边抽屉(点击卡片展开)

#### 基本信息标签页
- 模型名称和版本
- 模型提供商
- 模型类型(对话/嵌入/图像生成)
- 模型描述
- 添加者和添加时间
- 状态开关(启用/禁用)

#### 连接配置标签页
- API端点URL输入
- API密钥管理:
  - API Key输入(遮蔽显示)
  - 测试连接按钮
  - 连接状态显示
- 超时设置:
  - 连接超时(秒)
  - 读取超时(秒)
- 高级配置:
  - 重试次数
  - 重试间隔
  - 代理设置

#### 模型参数标签页
- 默认参数配置:
  - Temperature滑块(0-2)
  - Max Tokens输入框
  - Top P滑块(0-1)
  - Frequency Penalty滑块(-2 to 2)
  - Presence Penalty滑块(-2 to 2)
  - Stop Sequences输入
- 参数说明提示
- 恢复默认值按钮

#### 成本配置标签页
- 计费模式选择(按Token/按请求)
- Token定价:
  - Input Token价格($/1M tokens)
  - Output Token价格($/1M tokens)
- 成本控制:
  - 每日成本上限
  - 单次调用成本上限
  - 告警阈值设置
- 当前月度成本统计

#### 性能监控标签页
- 实时性能指标:
  - 当前可用性状态
  - 平均响应时间
  - P95响应时间
  - P99响应时间
  - 当前并发数
- 性能趋势图表:
  - 过去24小时响应时间趋势
  - 调用成功率趋势
  - Token消耗趋势
- 错误统计:
  - 最近错误列表(时间、错误类型、错误信息)
  - 错误率图表

#### 使用情况标签页
- 使用统计:
  - 今日调用次数
  - 本周调用次数
  - 本月调用次数
  - 总Token消耗
  - 总成本
- 关联的Agent列表:
  - Agent名称
  - 使用频率
  - 跳转到Agent详情链接
- 最近调用记录:
  - 调用时间
  - 调用Agent
  - 请求Token数
  - 响应Token数
  - 响应时间
  - 查看详情链接

### 底部操作按钮(详情抽屉内)
- 保存配置按钮
- 取消按钮
- 测试模型连接按钮
- 查看API文档按钮
