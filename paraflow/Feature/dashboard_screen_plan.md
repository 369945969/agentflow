# Dashboard
系统监控总览页面,提供实时运行状态、关键性能指标、资源使用情况和告警信息的可视化展示,支持快速访问核心功能。

Layout Hierarchy:
- Header (Full-width):
  - Top Navigation Bar
- Content Container (Positioned below the header):
  - Left Sidebar
  - Main Content Area

## Top Navigation Bar
- Logo和平台名称
- 全局搜索框
- 快速创建按钮(新建Agent/编排/任务)
- 通知中心图标(含未读数量)
- 用户头像和设置入口

## Left Sidebar
- Dashboard导航(当前选中)
- Agent Management导航
- Model Management导航
- Orchestration Canvas导航
- Schedule Tasks导航
- Execution History导航
- Chat Interface导航
- 系统设置入口

## Main Content Area

### 系统状态概览卡片组
- 实时运行任务数(含趋势对比)
- 今日执行总数(成功/失败分布)
- 活跃Agent数量
- 平均响应时间

### 实时执行监控面板
- 当前运行中的编排流程列表(流程名称、执行ID、当前节点、开始时间)
- 执行状态实时更新(运行中/暂停/完成)
- 快速操作按钮(查看详情/暂停/停止)

### 性能指标图表区域
- 过去24小时执行成功率趋势图
- 平均执行时长趋势图
- 资源使用率图表(CPU/内存/API调用量)
- 错误率统计图

### 告警与异常面板
- 最近告警列表(级别、时间、描述、关联任务)
- 异常执行记录(失败任务、错误信息、重试状态)
- 快速跳转到执行详情链接

### 快速访问区域
- 最近编辑的编排流程卡片(名称、最后修改时间、快速编辑按钮)
- 常用Agent角色卡片(名称、使用频率、快速配置按钮)
- 即将执行的定时任务列表(任务名、下次执行时间、编辑按钮)
