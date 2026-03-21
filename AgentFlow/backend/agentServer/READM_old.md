我来为你设计一个完整的 DeepAgent 风格的 Agent 管理核心处理程序。这是一个复杂的分布式 Agent 系统，我会从架构设计、核心模块、数据流、状态管理等多个维度进行详细设计。
一、系统架构总览
plain
复制
┌─────────────────────────────────────────────────────────────────┐
│                        WebSocket Gateway                         │
│  (连接管理 | 心跳检测 | 消息路由 | 负载均衡 | 认证授权)            │
└────────────────────┬────────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────────┐
│                      AgentCore (核心调度层)                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐    │
│  │ 消息路由器    │  │ 会话管理器    │  │ 状态机管理器         │    │
│  │ MessageRouter│  │SessionManager│  │    StateMachine      │    │
│  └──────────────┘  └──────────────┘  └──────────────────────┘    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐    │
│  │ 智能分发器    │  │ 日志持久化    │  │ 配置管理中心         │    │
│  │  Dispatcher  │  │LogPersistence│  │   ConfigManager      │    │
│  └──────────────┘  └──────────────┘  └──────────────────────┘    │
└────────────────────┬────────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
┌──────────────┐ ┌──────────┐ ┌──────────────┐
│  SubAgent    │ │ SubAgent │ │  SubAgent    │
│  OpenCode-1  │ │OpenCode-2│ │  OpenCode-N  │
│ (代码执行)    │ │(数据分析) │ │  (其他技能)   │
└──────────────┘ └──────────┘ └──────────────┘
        │            │            │
        └────────────┴────────────┘
                     ▼
            ┌────────────────┐
            │  结果聚合器      │
            │ResultAggregator│
            └────────────────┘
二、核心数据模型设计
2.1 用户会话配置模型 (UserProfile)
TypeScript
复制
interface UserProfile {
  userId: string;                    // 用户唯一标识
  currentSessionId: string | null;   // 当前活跃会话ID
  config: UserInteractionConfig;     // 用户交互配置
  metadata: UserMetadata;            // 用户元数据
  createdAt: Date;
  updatedAt: Date;
}

interface UserInteractionConfig {
  // 实时返回控制
  streamMode: 'realtime' | 'batch' | 'adaptive';  // 实时/批量/自适应
  enableThinking: boolean;          // 是否显示思考过程
  thinkingDepth: 'none' | 'light' | 'deep'; // 思考深度
  
  // 结果返回策略
  returnStrategy: ReturnStrategy;
  
  // 消息处理偏好
  groupMessagePolicy: GroupMessagePolicy;
  
  // 超时配置
  timeoutMs: number;
  maxRetries: number;
}

interface ReturnStrategy {
  type: 'immediate' | 'final_only' | 'chunked';
  // immediate: 每个chunk都返回
  // final_only: 只返回最终结果
  // chunked: 按语义块返回
  
  finalResultMarker: string;        // 标记最终结果的关键词
  chunkSize: number;               // 分块大小（字符数）
  debounceMs: number;              // 防抖时间
}
2.2 会话模型 (Session)
TypeScript
复制
interface Session {
  sessionId: string;               // 会话唯一ID (UUID)
  userId: string;                  // 所属用户
  type: 'single' | 'group';        // 单聊或群聊
  status: SessionStatus;           // 会话状态
  
  // 群聊特有
  groupInfo?: GroupInfo;
  
  // 执行上下文
  executionContext: ExecutionContext;
  
  // 消息历史
  messages: Message[];
  
  createdAt: Date;
  updatedAt: Date;
  expiresAt: Date;                // 会话过期时间
}

interface GroupInfo {
  groupId: string;
  groupName: string;
  members: GroupMember[];
  // 智能路由配置
  routingRules: RoutingRule[];
}

interface GroupMember {
  userId: string;
  role: 'admin' | 'member' | 'bot';
  skills: Skill[];                // 成员技能标签
  description: string;            // 角色描述
  priority: number;               // 处理优先级
  isActive: boolean;              // 是否在线
}
2.3 消息模型 (Message)
TypeScript
复制
interface Message {
  messageId: string;
  sessionId: string;
  userId: string;                  // 发送者
  
  type: MessageType;
  content: MessageContent;
  metadata: MessageMetadata;
  
  // 路由信息
  routing: RoutingInfo;
  
  // 处理状态
  processingStatus: ProcessingStatus;
  
  timestamp: Date;
}

type MessageType = 
  | 'text'           // 纯文本
  | 'command'        // 命令消息
  | 'system'         // 系统消息
  | 'multimodal'     // 多模态(图文等)
  | 'group_mention'; // 群聊@消息

interface MessageContent {
  text: string;
  attachments?: Attachment[];
  mentions?: string[];              // @的用户ID列表
  command?: CommandPayload;       // 命令载荷
}

interface RoutingInfo {
  targetAgentId?: string;          // 指定目标Agent
  targetSkills?: string[];         // 目标技能要求
  routingDecision?: RoutingDecision; // 路由决策记录
}

interface ProcessingStatus {
  state: 'pending' | 'routing' | 'assigned' | 'processing' | 'completed' | 'failed';
  assignedAgent?: string;          // 分配的SubAgent
  startedAt?: Date;
  completedAt?: Date;
  error?: ErrorInfo;
}
2.4 SubAgent (OpenCode) 配置模型
TypeScript
复制
interface SubAgentConfig {
  agentId: string;                 // Agent唯一ID
  name: string;
  type: 'opencode' | 'custom';    // 后端类型
  endpoint: string;                // OpenCode服务地址
  
  // 能力定义
  capabilities: AgentCapability[];
  
  // 执行配置
  executionConfig: AgentExecutionConfig;
  
  // 健康状态
  health: HealthStatus;
}

interface AgentCapability {
  skillId: string;                 // 技能ID
  skillName: string;               // 技能名称
  description: string;             // 技能描述
  keywords: string[];              // 触发关键词
  inputSchema: JSONSchema;         // 输入参数Schema
  outputSchema: JSONSchema;        // 输出参数Schema
  priority: number;                // 技能优先级
}

interface AgentExecutionConfig {
  // 与OpenCode通信配置
  requestTimeout: number;
  maxConcurrentRequests: number;
  
  // 结果标记配置
  resultMarkers: {
    finalPrefix: string;           // 如 "[FINAL]"
    thinkingPrefix: string;        // 如 "[THINKING]"
    errorPrefix: string;           // 如 "[ERROR]"
    chunkSeparator: string;        // 分块分隔符
  };
  
  // 流式处理配置
  streaming: {
    enabled: boolean;
    bufferSize: number;
    flushInterval: number;
  };
}
三、核心模块详细设计
3.1 AgentCore 主控制器
TypeScript
复制
class AgentCore {
  private wsGateway: WebSocketGateway;
  private messageRouter: MessageRouter;
  private sessionManager: SessionManager;
  private dispatcher: SmartDispatcher;
  private logPersistence: LogPersistence;
  private subAgentPool: SubAgentPool;
  private resultAggregator: ResultAggregator;
  
  constructor(config: AgentCoreConfig) {
    this.initializeComponents(config);
    this.setupEventHandlers();
  }
  
  // 主入口：处理WebSocket消息
  async handleWebSocketMessage(
    connectionId: string, 
    rawMessage: WebSocketRawMessage
  ): Promise<void> {
    try {
      // 1. 解析和验证消息
      const message = await this.parseMessage(rawMessage);
      
      // 2. 获取或创建会话
      const session = await this.sessionManager.getOrCreateSession(
        message.userId, 
        message.sessionId,
        message.type
      );
      
      // 3. 判断消息类型并路由
      if (message.type === 'group_mention' || session.type === 'group') {
        await this.handleGroupMessage(message, session);
      } else {
        await this.handleSingleMessage(message, session);
      }
      
    } catch (error) {
      await this.handleError(connectionId, error);
    }
  }
  
  // 处理单聊消息
  private async handleSingleMessage(
    message: Message, 
    session: Session
  ): Promise<void> {
    // 1. 确定目标SubAgent（基于历史、技能匹配等）
    const targetAgent = await this.dispatcher.selectAgent(message, session);
    
    // 2. 创建执行上下文
    const context = this.createExecutionContext(message, session, targetAgent);
    
    // 3. 发送给SubAgent处理
    await this.dispatchToSubAgent(context);
  }
  
  // 处理群聊消息（核心智能路由逻辑）
  private async handleGroupMessage(
    message: Message, 
    session: Session
  ): Promise<void> {
    const groupInfo = session.groupInfo!;
    
    // 1. 解析消息意图和所需技能
    const intent = await this.analyzeIntent(message.content.text);
    const requiredSkills = intent.requiredSkills;
    
    // 2. 智能匹配群成员
    const candidates = this.findQualifiedMembers(
      groupInfo.members, 
      requiredSkills,
      message.content.mentions
    );
    
    // 3. 选择最佳处理者
    const selectedMember = this.selectBestCandidate(candidates, intent);
    
    if (!selectedMember) {
      // 无匹配成员，使用默认Agent或拒绝
      await this.handleNoMatch(message, session);
      return;
    }
    
    // 4. 构建角色化提示
    const rolePrompt = this.buildRolePrompt(selectedMember, message);
    
    // 5. 创建带角色信息的执行上下文
    const context: ExecutionContext = {
      ...this.createExecutionContext(message, session, selectedMember.agentConfig),
      roleIdentity: {
        userId: selectedMember.userId,
        roleDescription: selectedMember.description,
        skills: selectedMember.skills,
        persona: selectedMember.persona || 'default'
      },
      systemPrompt: rolePrompt
    };
    
    // 6. 发送给对应的SubAgent
    await this.dispatchToSubAgent(context);
  }
  
  // 构建角色化系统提示
  private buildRolePrompt(member: GroupMember, message: Message): string {
    return `
# 角色设定
你是群聊成员 "${member.userId}"，角色描述: ${member.description}

# 你的技能
${member.skills.map(s => `- ${s.name}: ${s.description}`).join('\n')}

# 当前群聊上下文
群ID: ${message.sessionId}
触发消息: "${message.content.text}"
${message.content.mentions ? `你被以下成员提及: ${message.content.mentions.join(', ')}` : ''}

# 任务要求
请基于你的角色设定和技能，处理上述消息。保持角色一致性，使用符合角色身份的语气和专业知识。
    `.trim();
  }
}
3.2 智能分发器 (SmartDispatcher)
TypeScript
复制
class SmartDispatcher {
  private agentPool: SubAgentPool;
  private skillRegistry: SkillRegistry;
  private loadBalancer: LoadBalancer;
  private contextAnalyzer: ContextAnalyzer;
  
  // 核心：选择最合适的SubAgent
  async selectAgent(
    message: Message, 
    session: Session
  ): Promise<SubAgentConfig> {
    
    // 1. 如果消息指定了目标Agent，直接返回
    if (message.routing.targetAgentId) {
      return this.agentPool.getAgent(message.routing.targetAgentId);
    }
    
    // 2. 分析消息意图和所需技能
    const intent = await this.contextAnalyzer.analyze(message, session);
    
    // 3. 基于技能匹配度筛选候选Agent
    const candidates = this.skillRegistry.findAgentsBySkills(intent.requiredSkills);
    
    // 4. 多维度评分排序
    const scoredCandidates = candidates.map(agent => ({
      agent,
      score: this.calculateScore(agent, message, session, intent)
    })).sort((a, b) => b.score - a.score);
    
    // 5. 负载均衡选择
    return this.loadBalancer.select(scoredCandidates.map(c => c.agent));
  }
  
  private calculateScore(
    agent: SubAgentConfig, 
    message: Message,
    session: Session,
    intent: IntentAnalysis
  ): number {
    let score = 0;
    
    // 技能匹配度 (40%)
    const skillMatch = this.calculateSkillMatch(agent.capabilities, intent);
    score += skillMatch * 0.4;
    
    // 历史会话连续性 (30%)
    const continuity = this.calculateSessionContinuity(agent.agentId, session);
    score += continuity * 0.3;
    
    // 当前负载情况 (20%)
    const loadScore = 1 - (agent.health.currentLoad / agent.health.maxCapacity);
    score += loadScore * 0.2;
    
    // 响应时间历史 (10%)
    const latencyScore = 1 / (1 + agent.health.averageLatency);
    score += latencyScore * 0.1;
    
    return score;
  }
}
3.3 SubAgent 执行器 (OpenCode 适配层)
TypeScript
复制
class SubAgentExecutor {
  private httpClient: AxiosInstance;
  private streamHandlers: Map<string, StreamHandler>;
  private logPersistence: LogPersistence;
  
  // 执行主方法
  async execute(context: ExecutionContext): Promise<void> {
    const { message, session, targetAgent, userConfig } = context;
    
    // 1. 构建OpenCode请求
    const request = this.buildOpenCodeRequest(context);
    
    // 2. 根据用户配置选择执行模式
    if (userConfig.streamMode === 'realtime') {
      await this.executeStreaming(context, request);
    } else {
      await this.executeBatch(context, request);
    }
  }
  
  // 流式执行（实时返回）
  private async executeStreaming(
    context: ExecutionContext,
    request: OpenCodeRequest
  ): Promise<void> {
    const { message, session, userConfig } = context;
    const streamId = `${session.sessionId}_${message.messageId}`;
    
    // 创建流处理器
    const handler = new StreamHandler({
      onChunk: (chunk: StreamChunk) => {
        // 实时返回给WebSocket
        this.sendToWebSocket(context.connectionId, {
          type: 'stream_chunk',
          sessionId: session.sessionId,
          messageId: message.messageId,
          content: chunk.content,
          isThinking: chunk.isThinking,
          timestamp: new Date()
        });
        
        // 同时写入日志（追加模式）
        this.logPersistence.appendLog(
          context.userId,
          session.sessionId,
          chunk
        );
      },
      
      onFinal: (finalResult: FinalResult) => {
        // 发送最终结果标记
        this.sendToWebSocket(context.connectionId, {
          type: 'stream_end',
          sessionId: session.sessionId,
          messageId: message.messageId,
          content: finalResult.content,
          isFinal: true,
          usage: finalResult.usage
        });
        
        // 完成日志写入
        this.logPersistence.finalizeLog(
          context.userId,
          session.sessionId,
          finalResult
        );
      },
      
      markers: context.targetAgent.executionConfig.resultMarkers
    });
    
    this.streamHandlers.set(streamId, handler);
    
    try {
      // 发送SSE请求到OpenCode
      const response = await this.httpClient.post(
        `${context.targetAgent.endpoint}/v1/chat/completions`,
        {
          ...request,
          stream: true,
          // 添加AgentCore的标记要求
          response_format: {
            type: 'text',
            // 要求OpenCode在最终结果前添加标记
            final_marker: context.targetAgent.executionConfig.resultMarkers.finalPrefix
          }
        },
        { responseType: 'stream' }
      );
      
      // 处理SSE流
      for await (const line of response.data) {
        handler.processLine(line);
      }
      
    } finally {
      this.streamHandlers.delete(streamId);
    }
  }
  
  // 批量执行（只返回最终结果）
  private async executeBatch(
    context: ExecutionContext,
    request: OpenCodeRequest
  ): Promise<void> {
    const { message, session, userConfig } = context;
    
    // 非流式请求
    const response = await this.httpClient.post(
      `${context.targetAgent.endpoint}/v1/chat/completions`,
      {
        ...request,
        stream: false,
        // 明确要求返回完整结果
        return_full_response: true
      }
    );
    
    const result = response.data.choices[0].message.content;
    
    // 一次性写入日志
    this.logPersistence.writeFullLog(
      context.userId,
      session.sessionId,
      {
        request: message.content,
        response: result,
        timestamp: new Date(),
        usage: response.data.usage
      }
    );
    
    // 发送最终结果
    this.sendToWebSocket(context.connectionId, {
      type: 'message_complete',
      sessionId: session.sessionId,
      messageId: message.messageId,
      content: result,
      isFinal: true
    });
  }
}
3.4 日志持久化模块 (核心需求实现)
TypeScript
复制
class LogPersistence {
  private basePath: string;
  private fileHandles: Map<string, FileHandle>;  // 活跃文件句柄缓存
  private writeQueues: Map<string, WriteQueue>;  // 写入队列
  
  constructor(config: { basePath: string }) {
    this.basePath = config.basePath;
    this.ensureDirectoryStructure();
  }
  
  // 追加写入日志（流式模式）
  async appendLog(
    userId: string,
    sessionId: string,
    chunk: StreamChunk
  ): Promise<void> {
    const filePath = this.getLogFilePath(userId, sessionId);
    
    // 确保目录存在
    await this.ensureUserDirectory(userId);
    
    // 格式化日志条目
    const logEntry = this.formatLogEntry(chunk);
    
    // 使用队列保证写入顺序
    const queue = this.getWriteQueue(filePath);
    await queue.add(() => this.doAppend(filePath, logEntry));
  }
  
  // 完成日志（添加结束标记）
  async finalizeLog(
    userId: string,
    sessionId: string,
    finalResult: FinalResult
  ): Promise<void> {
    const filePath = this.getLogFilePath(userId, sessionId);
    
    const finalEntry = `
================================================================================
[SESSION_END] ${new Date().toISOString()}
Final Result: ${finalResult.content}
Usage: ${JSON.stringify(finalResult.usage)}
================================================================================
`;
    
    await this.doAppend(filePath, finalEntry);
    
    // 关闭文件句柄
    await this.closeFileHandle(filePath);
  }
  
  // 完整写入（批量模式）
  async writeFullLog(
    userId: string,
    sessionId: string,
    data: CompleteLogData
  ): Promise<void> {
    const filePath = this.getLogFilePath(userId, sessionId);
    await this.ensureUserDirectory(userId);
    
    const content = this.formatCompleteLog(data);
    await fs.writeFile(filePath, content, 'utf-8');
  }
  
  // 路径生成：/logs/{userId}/{sessionId}.log
  private getLogFilePath(userId: string, sessionId: string): string {
    // 安全处理：防止路径遍历
    const safeUserId = userId.replace(/[\\/:*?"<>|]/g, '_');
    const safeSessionId = sessionId.replace(/[\\/:*?"<>|]/g, '_');
    
    return path.join(this.basePath, safeUserId, `${safeSessionId}.log`);
  }
  
  private formatLogEntry(chunk: StreamChunk): string {
    const timestamp = new Date().toISOString();
    const type = chunk.isThinking ? 'THINKING' : 
                 chunk.isFinal ? 'FINAL' : 'CHUNK';
    
    return `[${timestamp}] [${type}] ${chunk.content}\n`;
  }
  
  private async doAppend(filePath: string, content: string): Promise<void> {
    let handle = this.fileHandles.get(filePath);
    
    if (!handle) {
      handle = await fs.open(filePath, 'a');
      this.fileHandles.set(filePath, handle);
    }
    
    await handle.write(content);
    // 可选：定期fsync保证持久化
  }
  
  // 读取历史日志（用于上下文构建）
  async readSessionHistory(
    userId: string, 
    sessionId: string,
    options: { limit?: number; offset?: number } = {}
  ): Promise<LogEntry[]> {
    const filePath = this.getLogFilePath(userId, sessionId);
    
    try {
      const content = await fs.readFile(filePath, 'utf-8');
      return this.parseLogContent(content, options);
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
        return []; // 文件不存在返回空
      }
      throw error;
    }
  }
}
四、群聊智能路由算法
TypeScript
复制
class GroupRoutingEngine {
  // 核心路由决策
  async routeGroupMessage(
    message: Message,
    groupInfo: GroupInfo
  ): Promise<RoutingDecision> {
    
    // 1. 提取消息中的@提及
    const mentions = message.content.mentions || [];
    
    // 2. 如果明确@了某人，优先路由给被@者
    if (mentions.length > 0) {
      const mentionedMembers = groupInfo.members.filter(m => 
        mentions.includes(m.userId)
      );
      
      if (mentionedMembers.length > 0) {
        return {
          type: 'direct',
          target: this.selectFromMentioned(mentionedMembers, message),
          reason: 'explicit_mention'
        };
      }
    }
    
    // 3. 分析消息意图和所需技能
    const analysis = await this.analyzeMessageRequirements(message);
    
    // 4. 基于技能匹配查找候选
    const candidates = this.findSkillMatches(
      groupInfo.members, 
      analysis.requiredSkills
    );
    
    // 5. 多因素决策
    if (candidates.length === 0) {
      return { type: 'none', reason: 'no_skill_match' };
    }
    
    if (candidates.length === 1) {
      return {
        type: 'skill_match',
        target: candidates[0],
        reason: 'single_match'
      };
    }
    
    // 多候选时，使用决策算法
    return this.resolveMultipleCandidates(candidates, message, analysis);
  }
  
  // 多候选冲突解决
  private resolveMultipleCandidates(
    candidates: GroupMember[],
    message: Message,
    analysis: RequirementAnalysis
  ): RoutingDecision {
    
    // 评分维度：
    const scored = candidates.map(member => {
      let score = 0;
      
      // 1. 技能匹配度 (精确匹配加分，部分匹配减分)
      const skillScore = member.skills.reduce((sum, skill) => {
        const match = analysis.requiredSkills.find(rs => 
          rs.id === skill.id || skill.keywords.some(k => 
            analysis.keywords.includes(k)
          )
        );
        return sum + (match ? skill.priority : 0);
      }, 0);
      score += skillScore * 0.4;
      
      // 2. 角色活跃度 (最近发言时间)
      const recencyScore = this.calculateRecencyScore(member.userId, message.sessionId);
      score += recencyScore * 0.2;
      
      // 3. 负载均衡 (避免总是同一人处理)
      const loadScore = 1 - (member.recentTaskCount / 10); // 假设最近10条
      score += loadScore * 0.2;
      
      // 4. 角色描述语义相似度
      const semanticScore = this.calculateSemanticSimilarity(
        member.description,
        message.content.text
      );
      score += semanticScore * 0.2;
      
      return { member, score };
    });
    
    scored.sort((a, b) => b.score - a.score);
    
    return {
      type: 'intelligent_routing',
      target: scored[0].member,
      alternatives: scored.slice(1, 3).map(s => s.member),
      reason: `skill_match:${scored[0].score.toFixed(2)}`,
      confidence: scored[0].score
    };
  }
}
五、WebSocket 消息协议设计
TypeScript
复制
// 客户端 -> AgentCore
interface ClientToServerMessage {
  // 基础字段
  messageId: string;        // 客户端生成的消息ID
  sessionId?: string;       // 现有会话ID（空则新建）
  userId: string;          // 发送者ID
  
  // 内容
  content: {
    text: string;
    attachments?: Attachment[];
    mentions?: string[];    // 群聊@列表
  };
  
  // 控制指令
  command?: {
    type: 'new_session' | 'close_session' | 'update_config';
    payload: any;
  };
  
  // 临时覆盖配置（可选）
  overrideConfig?: Partial<UserInteractionConfig>;
  
  timestamp: number;
}

// AgentCore -> 客户端
interface ServerToClientMessage {
  type: 'stream_start' | 'stream_chunk' | 'stream_end' | 
        'thinking' | 'error' | 'system';
  
  messageId: string;        // 对应请求的消息ID
  sessionId: string;
  
  payload: {
    content?: string;       // 文本内容
    isFinal?: boolean;     // 是否最终结果
    isThinking?: boolean;   // 是否思考过程
    progress?: number;      // 进度百分比
    
    // 元数据
    metadata?: {
      agentId: string;      // 处理该消息的SubAgent
      processingTime: number;
      tokenUsage?: TokenUsage;
    };
    
    error?: ErrorDetails;
  };
  
  timestamp: number;
}
六、目录结构与部署架构
plain
复制
deep-agent-core/
├── src/
│   ├── core/
│   │   ├── AgentCore.ts           # 主控制器
│   │   ├── MessageRouter.ts       # 消息路由器
│   │   └── SessionManager.ts      # 会话管理器
│   ├── agents/
│   │   ├── SubAgentPool.ts        # Agent池管理
│   │   ├── SubAgentExecutor.ts    # 执行器
│   │   └── adapters/
│   │       └── OpenCodeAdapter.ts  # OpenCode适配器
│   ├── routing/
│   │   ├── SmartDispatcher.ts     # 智能分发
│   │   ├── GroupRoutingEngine.ts  # 群聊路由
│   │   └── IntentAnalyzer.ts      # 意图分析
│   ├── persistence/
│   │   ├── LogPersistence.ts      # 日志持久化
│   │   └── SessionStore.ts        # 会话存储(Redis)
│   ├── websocket/
│   │   ├── WebSocketGateway.ts    # WS网关
│   │   └── ConnectionManager.ts   # 连接管理
│   └── types/
│       └── index.ts               # 类型定义
├── config/
│   ├── agents.yaml                # SubAgent配置
│   ├── routing-rules.yaml         # 路由规则
│   └── default-configs.yaml       # 默认用户配置
├── logs/                          # 日志根目录
│   └── {userId}/
│       └── {sessionId}.log
├── docker-compose.yml
└── Dockerfile
七、关键流程时序图
plain
复制
用户A (单聊)              AgentCore              OpenCode              日志系统
   |                         |                      |                     |
   |--- 发送消息 ------------>|                      |                     |
   |                         |--- 1. 解析消息       |                     |
   |                         |--- 2. 获取会话配置     |                     |
   |                         |--- 3. 选择SubAgent    |                     |
   |                         |--- 4. 创建执行上下文   |                     |
   |                         |                      |                     |
   |                         |--- 5. 发送请求 ------>|                     |
   |                         |                      |--- 6. 处理          |
   |                         |<-- 7. SSE流开始 ------|                     |
   |                         |                      |                     |
   |                         |--- 8. 实时返回给WS -->|                     |
   |<-- 9. 收到chunk ---------|                      |                     |
   |                         |--- 10. 写入日志 ---->|                     |
   |                         |                      |                     |-- 11. 追加文件
   |                         |                      |                     |
   |<-- 12. 继续收chunk ------|                      |                     |
   |                         |                      |                     |
   |                         |<-- 13. 最终标记 -----|                     |
   |                         |--- 14. 发送结束标记 ->|                     |
   |<-- 15. 收到最终结果 ------|                      |                     |
   |                         |--- 16. 完成日志 ----->|                     |
   |                         |                      |                     |-- 17. 关闭文件
   |                         |                      |                     |
八、配置示例
yaml
复制
# agents.yaml - SubAgent配置
agents:
  - id: "opencode-coder"
    name: "代码助手"
    endpoint: "http://opencode-coder:8080"
    capabilities:
      - skillId: "code_gen"
        name: "代码生成"
        keywords: ["写代码", "编程", "function", "bug"]
        priority: 10
      - skillId: "code_review"
        name: "代码审查"
        keywords: ["review", "检查代码", "优化"]
        priority: 8
    
    executionConfig:
      resultMarkers:
        finalPrefix: "[FINAL_RESULT]"
        thinkingPrefix: "[THINKING_PROCESS]"
        chunkSeparator: "\n\n"
      streaming:
        enabled: true
        bufferSize: 1024
        flushInterval: 100  # ms

  - id: "opencode-analyst"
    name: "数据分析师"
    endpoint: "http://opencode-analyst:8080"
    capabilities:
      - skillId: "data_analysis"
        name: "数据分析"
        keywords: ["分析", "统计", "图表", "数据"]
        priority: 10

# default-configs.yaml
defaultUserConfig:
  streamMode: "realtime"
  enableThinking: true
  thinkingDepth: "light"
  returnStrategy:
    type: "chunked"
    finalResultMarker: "[FINAL_RESULT]"
    chunkSize: 500
    debounceMs: 50
  groupMessagePolicy:
    autoRoute: true
    requireExplicitMention: false
    maxResponseTime: 30000
九、关键技术决策说明
表格
决策点	方案	理由
结果标记方式	双向约定	AgentCore发送finalResultMarker给OpenCode，OpenCode在最终输出前添加该标记；同时AgentCore通过SSE解析检测标记
日志写入策略	内存队列+批量追加	避免频繁IO，保证顺序，同时实时落盘
群聊路由	技能匹配+语义相似度+负载均衡	避免单点过载，保证专业问题由专业人士回答
流式控制	用户级配置+消息级覆盖	灵活满足不同场景需求
会话隔离	UserId+SessionId双重隔离	支持同一用户多会话并行
这个设计完整实现了你的需求：WebSocket接入、智能路由、OpenCode后端适配、用户维度日志存储、实时/批量返回模式、群聊角色识别等。需要我针对某个模块进一步细化吗？
目前只有用户单聊，群聊的决策也要做出来，然后我是用golang语言做。你帮我把逻辑用文字写的详细一点。我会交给AI去编写。
编辑
复制
分享
我来为你详细设计这个 AgentCore 系统的逻辑，专门针对 Go语言实现 和 单聊+群聊 场景。我会把每个模块的逻辑拆解得非常细致，方便你交给 AI 编写代码。
一、系统整体架构（Go语言版本）
plain
复制
┌─────────────────────────────────────────────────────────────┐
│                    WebSocket Gateway (Gorilla/Go)            │
│  - 连接管理(map[connId]*Connection)                         │
│  - 心跳检测(time.Ticker)                                    │
│  - 消息解析(json.Unmarshal)                               │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                      AgentCore (主控模块)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │ MessageRouter │  │SessionManager│  │ ConfigManager    │   │
│  │   消息路由     │  │   会话管理    │  │   配置管理        │   │
│  └──────────────┘  └──────────────┘  └──────────────────┘   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │SingleChat    │  │ GroupRouter   │  │ LogPersistence   │   │
│  │   单聊处理器  │  │   群聊路由器   │  │   日志持久化      │   │
│  └──────────────┘  └──────────────┘  └──────────────────┘   │
│  ┌──────────────┐  ┌──────────────┐                          │
│  │SubAgentPool  │  │ StreamManager│                          │
│  │  SubAgent池  │  │   流管理器    │                          │
│  └──────────────┘  └──────────────┘                          │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
┌──────────────┐ ┌──────────┐ ┌──────────────┐
│ OpenCode-1   │ │OpenCode-2│ │  OpenCode-N  │
│ (代码能力)    │ │(分析能力) │ │  (其他能力)   │
└──────────────┘ └──────────┘ └──────────────┘
二、核心数据结构定义（Go Struct）
2.1 用户配置结构
go
复制
// UserProfile 用户全局配置（存储在Redis或数据库）
type UserProfile struct {
    UserID          string                `json:"user_id"`
    CurrentSessionID string               `json:"current_session_id"` // 当前活跃会话
    Config          UserInteractionConfig  `json:"config"`
    CreatedAt       int64                 `json:"created_at"` // Unix时间戳
    UpdatedAt       int64                 `json:"updated_at"`
}

type UserInteractionConfig struct {
    // 流式模式：realtime(实时) / batch(批量) / adaptive(自适应)
    StreamMode string `json:"stream_mode"`
    
    // 是否显示思考过程
    EnableThinking bool `json:"enable_thinking"`
    
    // 返回策略
    ReturnStrategy ReturnStrategy `json:"return_strategy"`
    
    // 超时配置（毫秒）
    TimeoutMs   int `json:"timeout_ms"`
    MaxRetries  int `json:"max_retries"`
}

type ReturnStrategy struct {
    Type              string `json:"type"`               // immediate/final_only/chunked
    FinalResultMarker string `json:"final_result_marker"` // 如 "[FINAL]"
    ChunkSize         int    `json:"chunk_size"`          // 分块字符数
    DebounceMs        int    `json:"debounce_ms"`         // 防抖时间
}
2.2 会话结构
go
复制
// Session 会话信息（存储在Redis，TTL 24小时）
type Session struct {
    SessionID   string    `json:"session_id"`   // UUID生成
    UserID      string    `json:"user_id"`      // 创建者
    Type        string    `json:"type"`         // "single" 或 "group"
    Status      string    `json:"status"`       // "active", "paused", "closed"
    
    // 群聊特有字段
    GroupInfo   *GroupInfo `json:"group_info,omitempty"`
    
    // 执行上下文（最近使用的Agent，用于连续性）
    LastAgentID string    `json:"last_agent_id"`
    
    CreatedAt   int64     `json:"created_at"`
    UpdatedAt   int64     `json:"updated_at"`
    ExpiresAt   int64     `json:"expires_at"`   // 过期时间
}

type GroupInfo struct {
    GroupID     string        `json:"group_id"`
    GroupName   string        `json:"group_name"`
    Members     []GroupMember `json:"members"`      // 群成员列表
    CreatorID   string        `json:"creator_id"`
}

type GroupMember struct {
    UserID      string   `json:"user_id"`
    Role        string   `json:"role"`        // "admin", "member", "bot"
    Skills      []Skill  `json:"skills"`      // 成员技能
    Description string   `json:"description"` // 角色描述
    Priority    int      `json:"priority"`    // 处理优先级(1-10)
    IsActive    bool     `json:"is_active"`   // 是否在线
    Persona     string   `json:"persona"`     // 角色人设
}

type Skill struct {
    SkillID     string   `json:"skill_id"`
    SkillName   string   `json:"skill_name"`
    Keywords    []string `json:"keywords"`    // 触发关键词
    Description string   `json:"description"`
}
2.3 消息结构
go
复制
// Message 消息结构（WebSocket接收）
type Message struct {
    MessageID   string          `json:"message_id"`   // 客户端生成UUID
    SessionID   string          `json:"session_id"`   // 可为空，空则新建
    UserID      string          `json:"user_id"`      // 发送者
    Type        string          `json:"type"`         // "text", "command", "group"
    
    Content     MessageContent  `json:"content"`
    Metadata    MessageMetadata `json:"metadata"`
    Timestamp   int64           `json:"timestamp"`
}

type MessageContent struct {
    Text        string       `json:"text"`
    Attachments []Attachment `json:"attachments,omitempty"`
    Mentions    []string     `json:"mentions,omitempty"` // @的用户ID列表
}

type MessageMetadata struct {
    // 客户端可临时覆盖配置
    OverrideConfig *UserInteractionConfig `json:"override_config,omitempty"`
    
    // 群聊路由辅助信息
    RequiredSkills []string `json:"required_skills,omitempty"`
    TargetAgentID  string   `json:"target_agent_id,omitempty"`
}

// ProcessingStatus 消息处理状态（用于追踪）
type ProcessingStatus struct {
    State       string `json:"state"`        // pending/routing/assigned/processing/completed/failed
    AssignedAgent string `json:"assigned_agent,omitempty"`
    StartedAt   int64  `json:"started_at,omitempty"`
    CompletedAt int64  `json:"completed_at,omitempty"`
    Error       string `json:"error,omitempty"`
}
2.4 SubAgent (OpenCode) 配置
go
复制
// SubAgentConfig OpenCode配置（从配置文件加载）
type SubAgentConfig struct {
    AgentID    string            `json:"agent_id"`
    Name       string            `json:"name"`
    Endpoint   string            `json:"endpoint"`      // http://opencode-xxx:8080
    Capabilities []AgentCapability `json:"capabilities"`
    
    // 执行配置
    ExecutionConfig AgentExecutionConfig `json:"execution_config"`
    
    // 运行时状态
    Health HealthStatus `json:"-"` // 不序列化，运行时维护
}

type AgentCapability struct {
    SkillID      string   `json:"skill_id"`
    SkillName    string   `json:"skill_name"`
    Description  string   `json:"description"`
    Keywords     []string `json:"keywords"`      // 匹配用户请求的关键词
    InputSchema  map[string]interface{} `json:"input_schema"`
    Priority     int      `json:"priority"`      // 技能优先级
}

type AgentExecutionConfig struct {
    // 结果标记配置（关键：用于识别OpenCode返回的数据类型）
    ResultMarkers struct {
        FinalPrefix    string `json:"final_prefix"`     // 如 "[FINAL]"
        ThinkingPrefix string `json:"thinking_prefix"`  // 如 "[THINKING]"
        ErrorPrefix    string `json:"error_prefix"`     // 如 "[ERROR]"
    } `json:"result_markers"`
    
    // 流式配置
    Streaming struct {
        Enabled       bool `json:"enabled"`
        BufferSize    int  `json:"buffer_size"`
        FlushInterval int  `json:"flush_interval"` // 毫秒
    } `json:"streaming"`
    
    // 请求配置
    TimeoutMs          int `json:"timeout_ms"`
    MaxConcurrentReqs  int `json:"max_concurrent_reqs"`
}

type HealthStatus struct {
    CurrentLoad     int     `json:"current_load"`      // 当前并发请求数
    MaxCapacity     int     `json:"max_capacity"`
    AvgLatencyMs    float64 `json:"avg_latency_ms"`
    LastHealthCheck int64   `json:"last_health_check"`
    IsHealthy       bool    `json:"is_healthy"`
}
2.5 执行上下文（单次请求）
go
复制
// ExecutionContext 单次请求的执行上下文
type ExecutionContext struct {
    // 基础信息
    Message       Message           // 原始消息
    Session       Session           // 会话信息
    UserProfile   UserProfile       // 用户配置
    
    // 目标Agent
    TargetAgent   SubAgentConfig    // 选中的SubAgent
    
    // 连接信息
    ConnectionID  string            // WebSocket连接ID
    
    // 群聊特有（单聊时RoleIdentity为nil）
    RoleIdentity  *RoleIdentity     `json:"role_identity,omitempty"`
    
    // 系统提示词（群聊时注入角色信息）
    SystemPrompt  string            `json:"system_prompt"`
    
    // 实际使用的配置（合并默认+用户+消息覆盖）
    EffectiveConfig UserInteractionConfig
}

// RoleIdentity 群聊时的角色身份
type RoleIdentity struct {
    UserID          string   `json:"user_id"`
    RoleDescription string   `json:"role_description"`
    Skills          []Skill  `json:"skills"`
    Persona         string   `json:"persona"`
}
三、核心模块逻辑详解
3.1 WebSocket Gateway 模块
职责：管理所有WebSocket连接，解析消息，转发给AgentCore
核心逻辑：
plain
复制
1. 连接管理
   - 使用 map[string]*Connection 存储所有连接（sync.RWMutex保护）
   - Connection结构包含：connId, userId, wsConn, lastPingTime, sendChan
   
2. 消息读取循环（每个连接一个goroutine）
   - 循环读取 ws.ReadMessage()
   - 解析JSON到 Message 结构
   - 验证必填字段：message_id, user_id, content.text
   - 如果 session_id 为空，生成新的UUID作为session_id
   
3. 消息投递
   - 调用 AgentCore.HandleMessage(connId, message)
   - 使用goroutine异步处理，不阻塞读取循环
   
4. 心跳检测
   - 启动定时器（每30秒）
   - 遍历所有连接，检查 lastPingTime
   - 超过60秒未收到ping，关闭连接
   
5. 发送消息给客户端
   - 每个连接有缓冲sendChan（大小100）
   - 单独goroutine执行 writeLoop()
   - 从sendChan读取，调用 ws.WriteJSON()
关键代码逻辑描述：
go
复制
// Connection 结构
type Connection struct {
    ID          string
    UserID      string
    Conn        *websocket.Conn
    SendChan    chan ServerMessage  // 待发送消息队列
    LastPing    int64               // 最后心跳时间
    IsClosed    bool
    mu          sync.RWMutex
}

// HandleConnection 处理新连接
func (gw *Gateway) HandleConnection(ws *websocket.Conn, userId string) {
    connId := generateUUID()
    conn := &Connection{
        ID:       connId,
        UserID:   userId,
        Conn:     ws,
        SendChan: make(chan ServerMessage, 100),
        LastPing: time.Now().Unix(),
    }
    
    // 存储连接
    gw.connections.Store(connId, conn)
    
    // 启动两个goroutine
    go conn.readLoop(gw.agentCore)   // 读取消息
    go conn.writeLoop()              // 发送消息
}

// readLoop 读取循环
func (c *Connection) readLoop(core *AgentCore) {
    defer func() {
        c.Close()
        core.HandleDisconnect(c.ID)
    }()
    
    for {
        _, data, err := c.Conn.ReadMessage()
        if err != nil {
            return // 连接断开
        }
        
        var msg Message
        if err := json.Unmarshal(data, &msg); err != nil {
            c.SendError("invalid_json", err.Error())
            continue
        }
        
        // 填充connId和生成sessionId（如果为空）
        if msg.SessionID == "" {
            msg.SessionID = generateUUID()
        }
        
        // 异步处理（不阻塞读取）
        go core.HandleMessage(c.ID, msg)
    }
}
3.2 AgentCore 主控模块
职责：中央调度器，协调所有子模块
核心逻辑流程：
plain
复制
HandleMessage(connId, message) 主入口
    │
    ├── 1. 获取或创建会话
    │   └── 调用 SessionManager.GetOrCreateSession(userId, sessionId, msgType)
    │       - 如果sessionId存在，从Redis获取
    │       - 如果不存在，新建Session（判断是单聊还是群聊）
    │       - 单聊：Type="single"
    │       - 群聊：Type="group"，需要初始化GroupInfo
    │
    ├── 2. 获取用户配置
    │   └── 调用 ConfigManager.GetUserConfig(userId)
    │       - 从Redis获取UserProfile
    │       - 如果不存在，使用默认配置创建
    │
    ├── 3. 合并配置（默认 -> 用户配置 -> 消息覆盖）
    │   └── effectiveConfig = mergeConfig(default, userProfile.Config, message.Metadata.OverrideConfig)
    │
    ├── 4. 判断会话类型，分发处理
    │   ├── 如果是 Session.Type == "single"
    │   │   └── 调用 SingleChatHandler.Handle(context)
    │   └── 如果是 Session.Type == "group"
    │       └── 调用 GroupRouter.Handle(context)
    │
    └── 5. 错误处理（defer）
        └── 如果panic或error，发送错误消息给客户端
关键逻辑细节：
go
复制
// HandleMessage 主处理入口
func (ac *AgentCore) HandleMessage(connId string, msg Message) {
    defer func() {
        if r := recover(); r != nil {
            ac.sendError(connId, msg.SessionID, "internal_error", fmt.Sprintf("%v", r))
        }
    }()
    
    // 1. 获取会话
    session, err := ac.sessionMgr.GetOrCreateSession(msg.UserID, msg.SessionID, detectMsgType(msg))
    if err != nil {
        ac.sendError(connId, msg.SessionID, "session_error", err.Error())
        return
    }
    
    // 2. 获取用户配置
    userProfile, err := ac.configMgr.GetUserProfile(msg.UserID)
    if err != nil {
        // 使用默认配置
        userProfile = ac.configMgr.GetDefaultProfile(msg.UserID)
    }
    
    // 3. 更新当前会话（记录活跃会话）
    if userProfile.CurrentSessionID != session.SessionID {
        userProfile.CurrentSessionID = session.SessionID
        ac.configMgr.UpdateUserProfile(userProfile)
    }
    
    // 4. 构建执行上下文
    ctx := ExecutionContext{
        Message:         msg,
        Session:         *session,
        UserProfile:     *userProfile,
        ConnectionID:    connId,
        EffectiveConfig: ac.mergeConfig(userProfile.Config, msg.Metadata.OverrideConfig),
    }
    
    // 5. 路由分发
    if session.Type == "single" {
        ac.singleHandler.Handle(ctx)
    } else {
        ac.groupRouter.Handle(ctx)
    }
}

// detectMsgType 判断消息类型（单聊/群聊）
func detectMsgType(msg Message) string {
    // 如果消息包含mentions多人，或session已标记为group，则为群聊
    if len(msg.Content.Mentions) > 1 {
        return "group"
    }
    // 实际应根据session历史判断，这里简化
    return "single" // 默认单聊，创建session时确定
}
3.3 SessionManager 会话管理模块
职责：管理会话生命周期，存储在Redis
核心逻辑：
plain
复制
GetOrCreateSession(userId, sessionId, sessionType)
    │
    ├── 1. 构造Redis Key: "session:{sessionId}"
    │
    ├── 2. 尝试从Redis获取
    │   └── GET session:{sessionId}
    │       ├── 如果存在：检查userId是否匹配，返回会话
    │       └── 如果不存在：创建新会话
    │
    └── 3. 创建新会话（如果不存在）
        ├── 生成新的SessionID（如果传入为空）
        ├── 填充字段：
        │   SessionID: uuid
        │   UserID: userId
        │   Type: sessionType ("single"或"group")
        │   Status: "active"
        │   CreatedAt: now
        │   ExpiresAt: now + 24h
        └── 存入Redis：SET session:{sessionId} {json} EX 86400
        
UpdateSession(session)
    └── 更新Redis，重置TTL
    
GetSessionHistory(sessionId, limit)
    └── 从日志文件读取历史消息（用于上下文构建）
    
CloseSession(sessionId)
    └── 更新Status为"closed"，保留1小时后删除
Redis数据结构：
plain
复制
Key: session:{sessionId}
Value: JSON字符串
TTL: 86400秒（24小时）

Key: user_profile:{userId}
Value: JSON字符串
TTL: 永久（或7天）
3.4 SingleChatHandler 单聊处理模块
职责：处理单聊消息，选择SubAgent，执行请求
核心逻辑流程：
plain
复制
Handle(context)
    │
    ├── 1. 确定目标SubAgent
    │   └── 调用 selectAgent(context)
    │       ├── 如果 message.Metadata.TargetAgentID 不为空
    │       │   └── 直接使用该Agent
    │       └── 否则调用智能选择逻辑：
    │           ├── 分析消息内容（关键词提取）
    │           ├── 遍历所有SubAgent的Capabilities
    │           ├── 计算匹配分数：关键词命中数 * 技能优先级
    │           ├── 考虑会话连续性（优先使用LastAgentID）
    │           └── 选择分数最高的Agent
    │
    ├── 2. 更新会话LastAgentID（保持连续性）
    │   └── session.LastAgentID = selectedAgent.AgentID
    │       调用 SessionManager.UpdateSession(session)
    │
    ├── 3. 构建执行上下文
    │   └── context.TargetAgent = selectedAgent
    │
    ├── 4. 调用 SubAgentExecutor.Execute(context)
    │   └── 发送给OpenCode执行
    │
    └── 5. 处理执行结果
        └── 由Executor回调或返回结果
Agent选择算法详细逻辑：
go
复制
// selectAgent 选择最合适的SubAgent
func (sh *SingleChatHandler) selectAgent(ctx ExecutionContext) (SubAgentConfig, error) {
    // 1. 检查是否指定了Agent
    if ctx.Message.Metadata.TargetAgentID != "" {
        agent, ok := sh.agentPool.Get(ctx.Message.Metadata.TargetAgentID)
        if ok {
            return agent, nil
        }
    }
    
    // 2. 提取消息关键词（简单分词或正则）
    keywords := extractKeywords(ctx.Message.Content.Text)
    
    // 3. 计算每个Agent的匹配分数
    type scoredAgent struct {
        agent SubAgentConfig
        score float64
    }
    var scored []scoredAgent
    
    for _, agent := range sh.agentPool.GetAll() {
        score := 0.0
        
        // 3.1 技能匹配分数（40%权重）
        for _, cap := range agent.Capabilities {
            for _, keyword := range cap.Keywords {
                if contains(keywords, keyword) {
                    score += float64(cap.Priority) * 0.4
                }
            }
        }
        
        // 3.2 连续性分数（30%权重）
        if ctx.Session.LastAgentID == agent.AgentID {
            score += 10 * 0.3 // 高优先级保持连续
        }
        
        // 3.3 负载分数（20%权重）
        loadScore := 1.0 - (float64(agent.Health.CurrentLoad) / float64(agent.Health.MaxCapacity))
        score += loadScore * 10 * 0.2
        
        // 3.4 延迟分数（10%权重）
        latencyScore := 1.0 / (1.0 + agent.Health.AvgLatencyMs/1000.0)
        score += latencyScore * 10 * 0.1
        
        scored = append(scored, scoredAgent{agent, score})
    }
    
    // 4. 排序选择最高分
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].score > scored[j].score
    })
    
    if len(scored) == 0 {
        return SubAgentConfig{}, errors.New("no available agent")
    }
    
    return scored[0].agent, nil
}
3.5 GroupRouter 群聊路由模块（重点）
职责：解析群消息，智能选择处理成员，构建角色提示
核心逻辑流程：
plain
复制
Handle(context)
    │
    ├── 1. 解析群聊消息
    │   ├── 提取 @mentions 列表（从content.text中解析@userId）
    │   ├── 识别消息类型：@特定人 / @所有人 / 普通消息
    │   └── 分析消息意图（需要哪些技能）
    │
    ├── 2. 确定处理者（核心决策逻辑）
    │   └── 调用 determineHandler(context)
    │       ├── 情况A：消息包含明确的@mentions
    │       │   └── 优先路由给被@的成员（选择第一个被@的）
    │       │
    │       ├── 情况B：消息包含@所有人 或 无明确@但包含关键词
    │       │   └── 调用智能路由算法 selectBestMember()
    │       │       ├── 步骤1：提取所需技能（关键词匹配）
    │       │       ├── 步骤2：筛选有匹配技能的成员
    │       │       ├── 步骤3：多维度评分
    │       │       │   ├── 技能匹配度（40%）：技能keywords命中数
    │       │       │   ├── 角色相关性（30%）：描述与消息语义相似度
    │       │       │   ├── 负载均衡（20%）：最近处理消息数少优先
    │       │       │   └── 活跃度（10%）：在线状态、最近发言
    │       │       └── 步骤4：选择最高分成员
    │       │
    │       └── 情况C：无匹配成员
    │           └── 使用默认Agent（或回复"无合适处理人"）
    │
    ├── 3. 构建角色化执行上下文
    │   └── 调用 buildRoleContext(selectedMember, context)
    │       ├── 设置 ctx.RoleIdentity = {
    │       │   UserID: selectedMember.UserID,
    │       │   RoleDescription: selectedMember.Description,
    │       │   Skills: selectedMember.Skills,
    │       │   Persona: selectedMember.Persona
    │       │ }
    │       ├── 构建 SystemPrompt：
    │       │   "你是群成员'{UserID}'，角色：{Description}
    │       │    你的技能：{Skills列表}
    │       │    请基于你的角色处理以下消息：{原始消息}"
    │       └── 设置 ctx.SystemPrompt = 上述提示词
    │
    ├── 4. 确定使用的SubAgent
    │   └── 根据 selectedMember.Skills 选择对应能力的OpenCode
    │       （每个Skill可映射到特定SubAgentID）
    │
    ├── 5. 执行
    │   └── 调用 SubAgentExecutor.Execute(ctx)
    │       （此时ctx包含RoleIdentity，Executor会注入SystemPrompt）
    │
    └── 6. 记录群聊处理日志（谁处理了哪条消息）
        └── 用于后续负载均衡和统计
智能路由算法详细逻辑：
go
复制
// determineHandler 群聊处理者决策
func (gr *GroupRouter) determineHandler(ctx ExecutionContext) (*GroupMember, error) {
    groupInfo := ctx.Session.GroupInfo
    if groupInfo == nil {
        return nil, errors.New("not a group session")
    }
    
    // 1. 检查明确@mentions
    mentions := ctx.Message.Content.Mentions
    if len(mentions) > 0 {
        // 找到第一个被@且在线的成员
        for _, mentionId := range mentions {
            for _, member := range groupInfo.Members {
                if member.UserID == mentionId && member.IsActive {
                    return &member, nil // 明确指定，直接返回
                }
            }
        }
    }
    
    // 2. 智能路由：分析消息提取所需技能
    requiredSkills := gr.analyzeRequiredSkills(ctx.Message.Content.Text)
    
    // 3. 筛选候选成员（有任意匹配技能且在线）
    candidates := gr.findQualifiedMembers(groupInfo.Members, requiredSkills)
    
    if len(candidates) == 0 {
        // 无匹配，使用默认处理者（群主或第一个bot）
        return gr.getDefaultHandler(groupInfo.Members), nil
    }
    
    // 4. 多维度评分
    return gr.selectBestCandidate(candidates, ctx, requiredSkills), nil
}

// analyzeRequiredSkills 分析消息需要哪些技能
func (gr *GroupRouter) analyzeRequiredSkills(text string) []string {
    skills := []string{}
    
    // 关键词映射表（可配置化）
    skillKeywords := map[string][]string{
        "code":     {"代码", "编程", "bug", "函数", "写个", "实现"},
        "analysis": {"分析", "统计", "数据", "图表", "计算"},
        "writing":  {"写作", "写文章", "文案", "总结", "报告"},
        "image":    {"图片", "图像", "生成图", "画个"},
    }
    
    for skill, keywords := range skillKeywords {
        for _, kw := range keywords {
            if strings.Contains(text, kw) {
                skills = append(skills, skill)
                break
            }
        }
    }
    
    return skills
}

// selectBestCandidate 多维度评分选择最佳候选
func (gr *GroupRouter) selectBestCandidate(
    candidates []GroupMember, 
    ctx ExecutionContext,
    requiredSkills []string,
) *GroupMember {
    type scored struct {
        member *GroupMember
        score  float64
    }
    var scoredList []scored
    
    for _, member := range candidates {
        score := 0.0
        
        // 4.1 技能匹配度（40%）
        skillScore := 0.0
        for _, memberSkill := range member.Skills {
            for _, reqSkill := range requiredSkills {
                if memberSkill.SkillID == reqSkill {
                    skillScore += float64(memberSkill.Priority)
                }
            }
        }
        score += skillScore * 0.4
        
        // 4.2 角色描述语义相似度（30%）- 简单版本：关键词重叠
        descWords := strings.Fields(member.Description)
        msgWords := strings.Fields(ctx.Message.Content.Text)
        overlap := intersection(descWords, msgWords)
        semanticScore := float64(len(overlap)) / float64(len(descWords)+1)
        score += semanticScore * 10 * 0.3
        
        // 4.3 负载均衡（20%）- 获取成员最近处理消息数
        recentCount := gr.getRecentTaskCount(member.UserID, ctx.Session.SessionID)
        loadScore := 1.0 - (float64(recentCount) / 10.0) // 假设最近10条
        if loadScore < 0 {
            loadScore = 0
        }
        score += loadScore * 10 * 0.2
        
        // 4.4 优先级权重（10%）
        priorityScore := float64(member.Priority) / 10.0
        score += priorityScore * 10 * 0.1
        
        scoredList = append(scoredList, scored{&member, score})
    }
    
    // 排序返回最高分
    sort.Slice(scoredList, func(i, j int) bool {
        return scoredList[i].score > scoredList[j].score
    })
    
    return scoredList[0].member
}

// buildRoleContext 构建带角色信息的上下文
func (gr *GroupRouter) buildRoleContext(member *GroupMember, ctx ExecutionContext) ExecutionContext {
    // 复制上下文
    newCtx := ctx
    
    // 设置角色身份
    newCtx.RoleIdentity = &RoleIdentity{
        UserID:          member.UserID,
        RoleDescription: member.Description,
        Skills:          member.Skills,
        Persona:         member.Persona,
    }
    
    // 构建系统提示词（关键：让OpenCode扮演该角色）
    var skillDescs []string
    for _, s := range member.Skills {
        skillDescs = append(skillDescs, fmt.Sprintf("- %s: %s", s.SkillName, s.Description))
    }
    
    newCtx.SystemPrompt = fmt.Sprintf(`
# 角色设定
你是群聊成员 "%s"
角色描述：%s
人设风格：%s

# 你的专业技能
%s

# 当前任务
群ID：%s
触发消息："%s"

请基于你的角色设定，以第一人称处理上述消息。回复时要体现你的专业背景和角色特点。
`,
        member.UserID,
        member.Description,
        member.Persona,
        strings.Join(skillDescs, "\n"),
        ctx.Session.SessionID,
        ctx.Message.Content.Text,
    )
    
    return newCtx
}
3.6 SubAgentExecutor 执行模块
职责：与OpenCode通信，处理流式/批量返回，管理回调
核心逻辑流程：
plain
复制
Execute(context)
    │
    ├── 1. 准备请求
    │   └── buildOpenCodeRequest(context)
    │       ├── 构造messages数组：
    │       │   - 如果有SystemPrompt，插入system角色消息
    │       │   - 插入user角色消息（原始内容）
    │       ├── 读取历史消息（从日志文件）作为context
    │       └── 设置参数：temperature, max_tokens等
    │
    ├── 2. 判断执行模式（根据context.EffectiveConfig.StreamMode）
    │   ├── "realtime"  -> 执行流式模式 executeStream()
    │   ├── "batch"     -> 执行批量模式 executeBatch()
    │   └── "adaptive"  -> 根据消息长度自动选择（长消息用batch）
    │
    ├── 3. 流式模式 executeStream()
    │   ├── 3.1 创建StreamHandler（管理SSE流）
    │   ├── 3.2 发送HTTP POST到OpenCode /v1/chat/completions
    │   │   Header: Accept: text/event-stream
    │   │   Body: {stream: true, ...}
    │   ├── 3.3 读取SSE流（逐行解析）
    │   │   └── 对于每个data:开头的行：
    │   │       ├── 解析JSON，提取content字段
    │   │       ├── 判断类型（根据ResultMarkers前缀）：
    │   │       │   ├── 包含ThinkingPrefix -> isThinking=true
    │   │       │   ├── 包含FinalPrefix -> isFinal=true
    │   │       │   └── 其他 -> 普通chunk
    │   │       ├── 调用 onChunkReceived(chunk, context)
    │   │       │   ├── 发送给WebSocket（实时返回）
    │   │       │   └── 追加写入日志文件
    │   │       └── 如果isFinal，调用 onFinalReceived()
    │   └── 3.4 流结束处理
    │       ├── 发送stream_end标记给WebSocket
    │       └── 完成日志写入（添加结束标记）
    │
    ├── 4. 批量模式 executeBatch()
    │   ├── 4.1 发送HTTP POST（stream: false）
    │   ├── 4.2 等待完整响应
    │   ├── 4.3 一次性写入日志
    │   └── 4.4 发送完整结果给WebSocket
    │
    └── 5. 错误处理
        ├── 超时错误：重试（最多MaxRetries次）
        ├── 连接错误：标记Agent不健康，切换Agent重试
        └── 业务错误：记录日志，返回错误信息
关键逻辑细节：
go
复制
// StreamHandler 流处理器
type StreamHandler struct {
    context         ExecutionContext
    buffer          strings.Builder     // 内容缓冲区
    isThinking      bool                // 当前是否在思考阶段
    logFile         *os.File            // 日志文件句柄
    wsGateway       *WebSocketGateway   // 用于发送消息
}

// executeStream 流式执行
func (ex *SubAgentExecutor) executeStream(ctx ExecutionContext) error {
    // 1. 打开日志文件（追加模式）
    logPath := ex.logPersistence.GetLogFilePath(ctx.UserProfile.UserID, ctx.Session.SessionID)
    logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer logFile.Close()
    
    // 2. 构建请求
    reqBody := ex.buildRequestBody(ctx, true) // stream=true
    
    // 3. 创建HTTP请求
    httpReq, err := http.NewRequest("POST", 
        ctx.TargetAgent.Endpoint+"/v1/chat/completions", 
        bytes.NewReader(reqBody))
    if err != nil {
        return err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Accept", "text/event-stream")
    
    // 4. 发送请求
    client := &http.Client{Timeout: time.Duration(ctx.EffectiveConfig.TimeoutMs) * time.Millisecond}
    resp, err := client.Do(httpReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // 5. 读取SSE流
    reader := bufio.NewReader(resp.Body)
    handler := &StreamHandler{
        context:    ctx,
        logFile:    logFile,
        wsGateway:  ex.wsGateway,
    }
    
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        line = strings.TrimSpace(line)
        if !strings.HasPrefix(line, "data:") {
            continue
        }
        
        data := strings.TrimPrefix(line, "data:")
        if data == "[DONE]" {
            break
        }
        
        // 解析SSE数据
        var streamResp OpenCodeStreamResponse
        if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
            continue // 忽略解析错误
        }
        
        content := streamResp.Choices[0].Delta.Content
        if content == "" {
            continue
        }
        
        // 处理内容块
        handler.processChunk(content)
    }
    
    // 6. 发送结束标记
    handler.sendFinal()
    
    return nil
}

// processChunk 处理单个内容块
func (h *StreamHandler) processChunk(content string) {
    markers := h.context.TargetAgent.ExecutionConfig.ResultMarkers
    
    // 判断类型
    isThinking := strings.HasPrefix(content, markers.ThinkingPrefix)
    isFinal := strings.HasPrefix(content, markers.FinalPrefix)
    
    // 清理前缀
    cleanContent := content
    if isThinking {
        cleanContent = strings.TrimPrefix(content, markers.ThinkingPrefix)
        h.isThinking = true
    } else if isFinal {
        cleanContent = strings.TrimPrefix(content, markers.FinalPrefix)
    }
    
    // 1. 实时发送给WebSocket（根据用户配置）
    if h.context.EffectiveConfig.StreamMode == "realtime" {
        h.wsGateway.SendToConnection(h.context.ConnectionID, ServerMessage{
            Type:        "stream_chunk",
            SessionID:   h.context.Session.SessionID,
            MessageID:   h.context.Message.MessageID,
            Content:     cleanContent,
            IsThinking:  isThinking,
            IsFinal:     isFinal,
            Timestamp:   time.Now().UnixMilli(),
        })
    }
    
    // 2. 写入日志（带时间戳和类型标记）
    logEntry := fmt.Sprintf("[%s] [%s] %s\n", 
        time.Now().Format("2006-01-02T15:04:05"),
        h.getLogType(isThinking, isFinal),
        cleanContent)
    h.logFile.WriteString(logEntry)
    
    // 3. 缓冲内容
    h.buffer.WriteString(cleanContent)
    
    // 4. 如果是最终结果，标记
    if isFinal {
        h.isFinalReceived = true
    }
}

// sendFinal 发送最终标记
func (h *StreamHandler) sendFinal() {
    // 发送stream_end消息
    h.wsGateway.SendToConnection(h.context.ConnectionID, ServerMessage{
        Type:      "stream_end",
        SessionID: h.context.Session.SessionID,
        MessageID: h.context.Message.MessageID,
        Content:   h.buffer.String(), // 完整内容
        IsFinal:   true,
        Timestamp: time.Now().UnixMilli(),
    })
    
    // 写入结束标记到日志
    h.logFile.WriteString(fmt.Sprintf(
        "\n=== SESSION_END %s ===\nTotal Length: %d\n",
        time.Now().Format("2006-01-02T15:04:05"),
        h.buffer.Len(),
    ))
}
3.7 LogPersistence 日志持久化模块
职责：管理用户维度日志存储，支持追加写入和读取
核心逻辑：
plain
复制
目录结构：
/logs/
  └── {userId}/
      ├── {sessionId}.log          # 会话执行日志
      ├── {sessionId}_meta.json    # 会话元数据（可选）
      └── index.json               # 用户会话索引

文件格式（{sessionId}.log）：
[2024-01-15T10:30:00] [CHUNK] 这是第一段输出
[2024-01-15T10:30:01] [THINKING] 正在思考...
[2024-01-15T10:30:02] [CHUNK] 这是第二段
[2024-01-15T10:30:05] [FINAL] 这是最终结果
=== SESSION_END 2024-01-15T10:30:05 ===
Total Length: 156

核心方法：

GetLogFilePath(userId, sessionId)
    └── 返回：/logs/{userId}/{sessionId}.log
        （安全处理：替换特殊字符，防止目录遍历）

AppendLog(userId, sessionId, chunk)
    ├── 1. 获取文件路径
    ├── 2. 确保目录存在：os.MkdirAll(/logs/{userId}, 0755)
    ├── 3. 打开文件（追加模式）：os.OpenFile(path, O_APPEND|O_CREATE|O_WRONLY, 0644)
    ├── 4. 格式化：fmt.Sprintf("[%s] [%s] %s\n", timestamp, type, content)
    ├── 5. 写入：file.WriteString()
    └── 6. 注意：不关闭文件，保持句柄复用（使用sync.Pool或map缓存句柄）

FinalizeLog(userId, sessionId, finalResult)
    ├── 1. 获取文件句柄
    ├── 2. 写入结束标记和统计信息
    ├── 3. 关闭文件句柄
    └── 4. 从缓存中移除句柄

WriteFullLog(userId, sessionId, data) // 批量模式用
    ├── 1. 一次性写入完整内容
    └── 2. 包含请求和响应的完整记录

ReadSessionHistory(userId, sessionId, limit, offset)
    ├── 1. 打开日志文件
    ├── 2. 按行读取
    ├── 3. 解析每行：[timestamp] [type] content
    ├── 4. 返回LogEntry数组（用于构建OpenCode上下文）
    └── 5. 支持limit限制条数（从后往前读，取最近N条）

GetUserSessions(userId)
    └── 读取 /logs/{userId}/index.json 或扫描目录
文件句柄管理（性能优化）：
go
复制
type LogPersistence struct {
    basePath    string
    fileHandles map[string]*os.File  // path -> file handle
    mu          sync.RWMutex          // 保护map
    closeTimer  *time.Timer           // 定时清理空闲句柄
}

// 获取或创建文件句柄（带锁）
func (lp *LogPersistence) getFileHandle(path string) (*os.File, error) {
    lp.mu.RLock()
    if f, ok := lp.fileHandles[path]; ok {
        lp.mu.RUnlock()
        return f, nil
    }
    lp.mu.RUnlock()
    
    lp.mu.Lock()
    defer lp.mu.Unlock()
    
    // 双重检查
    if f, ok := lp.fileHandles[path]; ok {
        return f, nil
    }
    
    // 创建目录
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }
    
    // 打开文件
    f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    lp.fileHandles[path] = f
    return f, nil
}

// 定期清理空闲句柄（每5分钟）
func (lp *LogPersistence) startCleanupRoutine() {
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        for range ticker.C {
            lp.cleanupIdleHandles()
        }
    }()
}
3.8 ConfigManager 配置管理模块
职责：管理用户配置，支持默认配置和动态更新
核心逻辑：
plain
复制
配置文件结构：
/config/
  ├── default.yaml          # 系统默认配置
  └── agents.yaml           # SubAgent配置

Redis存储：
  Key: user_config:{userId} -> UserProfile JSON

核心方法：

LoadDefaultConfig()
    └── 从default.yaml加载到内存

LoadAgentConfigs()
    └── 从agents.yaml加载SubAgent列表
    └── 初始化SubAgentPool

GetUserProfile(userId)
    ├── 1. 尝试从Redis获取：GET user_config:{userId}
    ├── 2. 如果存在，反序列化返回
    ├── 3. 如果不存在：
    │   ├── 使用default.yaml创建新UserProfile
    │   ├── 设置UserID
    │   ├── 存入Redis
    │   └── 返回
    └── 4. 如果Redis失败，使用内存默认配置

UpdateUserProfile(profile)
    └── 序列化后存入Redis，更新UpdatedAt

MergeConfig(base, override1, override2)
    └── 递归合并配置（override优先级高）
    └── 用于：default <- user_config <- message_override

GetDefaultProfile(userId)
    └── 基于default.yaml创建，设置UserID
四、消息协议（WebSocket JSON格式）
4.1 客户端 -> 服务端
JSON
复制
{
  "message_id": "uuid-generated-by-client",
  "session_id": "existing-session-id-or-empty",
  "user_id": "user-123",
  "type": "text",
  "content": {
    "text": "帮我写个排序算法",
    "mentions": ["user-456", "user-789"],
    "attachments": []
  },
  "metadata": {
    "target_agent_id": "opencode-coder",
    "override_config": {
      "stream_mode": "realtime",
      "enable_thinking": true
    }
  },
  "timestamp": 1705312800000
}
4.2 服务端 -> 客户端
JSON
复制
// 流式块
{
  "type": "stream_chunk",
  "message_id": "uuid-from-request",
  "session_id": "session-id",
  "payload": {
    "content": "这是一个代码片段...",
    "is_thinking": false,
    "is_final": false,
    "progress": 45,
    "metadata": {
      "agent_id": "opencode-coder",
      "processing_time_ms": 1200
    }
  },
  "timestamp": 1705312800500
}

// 流结束
{
  "type": "stream_end",
  "message_id": "uuid-from-request",
  "session_id": "session-id",
  "payload": {
    "content": "完整结果内容...",
    "is_final": true,
    "usage": {
      "prompt_tokens": 150,
      "completion_tokens": 500,
      "total_tokens": 650
    }
  },
  "timestamp": 1705312801000
}

// 错误
{
  "type": "error",
  "session_id": "session-id",
  "payload": {
    "code": "agent_timeout",
    "message": "SubAgent响应超时",
    "retryable": true
  }
}
五、关键流程时序（文字描述）
5.1 单聊完整流程
plain
复制
1. 用户A打开WebSocket连接，发送消息
   WebSocketGateway接收 -> 解析JSON -> 调用AgentCore.HandleMessage(connId, msg)
   
2. AgentCore处理：
   - 检查session_id，发现为空，生成新的session-001
   - 创建Session{Type: "single", UserID: "user-a", ...}
   - 存入Redis，TTL=24h
   - 获取UserProfile（不存在则创建默认）
   - 合并配置：default + user + message.override

3. 路由到SingleChatHandler：
   - 分析消息"帮我写个排序算法"
   - 提取关键词：["排序", "算法", "写"]
   - 遍历SubAgent：
     * opencode-coder: 匹配"code"技能（关键词"写"命中），分数=8
     * opencode-analyst: 无匹配，分数=0
   - 选择opencode-coder（最高分）
   - 更新session.LastAgentID = "opencode-coder"

4. 构建ExecutionContext：
   - TargetAgent = opencode-coder配置
   - SystemPrompt为空（单聊不注入角色）
   - EffectiveConfig.StreamMode = "realtime"

5. SubAgentExecutor执行：
   - 判断StreamMode=="realtime"，走流式逻辑
   - 打开日志文件：/logs/user-a/session-001.log
   - 发送HTTP POST到opencode-coder:8080/v1/chat/completions
   - 设置stream=true，包含messages（user: "帮我写个排序算法"）

6. OpenCode返回SSE流：
   - data: {"choices":[{"delta":{"content":"[THINKING]分析需求..."}}]}
   - Executor识别ThinkingPrefix，标记isThinking=true
   - 发送WebSocket消息给user-a（type=stream_chunk, is_thinking=true）
   - 追加日志：[2024-01-15T10:00:00] [THINKING] 分析需求...

7. 继续接收流：
   - data: {"choices":[{"delta":{"content":"```go\nfunc QuickSort..."}}]}
   - 识别为普通CHUNK
   - 发送WebSocket（type=stream_chunk, content="```go\nfunc QuickSort..."）
   - 追加日志：[2024-01-15T10:00:01] [CHUNK] ```go\nfunc QuickSort...

8. 接收最终结果：
   - data: {"choices":[{"delta":{"content":"[FINAL]代码完成"}}]}
   - 识别FinalPrefix
   - 发送WebSocket（type=stream_chunk, is_final=true）
   - 追加日志：[2024-01-15T10:00:05] [FINAL] 代码完成

9. 流结束：
   - 发送WebSocket（type=stream_end）
   - 写入日志结束标记：=== SESSION_END ... ===
   - 关闭日志文件句柄

10. 用户A收到：
    - 实时看到思考过程（如果enable_thinking=true）
    - 实时看到代码生成过程
    - 最后收到stream_end确认完成
5.2 群聊完整流程
plain
复制
场景：群group-001有成员[user-a(群主), user-b(代码专家), user-c(文案专家)]
消息：user-a发送"@user-b 帮我review这段代码"

1. WebSocketGateway接收消息，解析：
   - mentions: ["user-b"]
   - text: "@user-b 帮我review这段代码"

2. AgentCore识别为群聊（mentions非空或session类型）：
   - 获取session（不存在则创建Type="group"）
   - GroupInfo.Members包含user-b信息：
     {UserID: "user-b", Skills: [{SkillID: "code_review", ...}], Description: "资深工程师..."}

3. 路由到GroupRouter：
   - 检测到mentions包含"user-b"
   - 查找成员：找到user-b，且IsActive=true
   - 直接选择user-b作为处理者（明确@优先）

4. 构建角色上下文：
   - RoleIdentity = {
       UserID: "user-b",
       RoleDescription: "资深工程师，擅长代码审查",
       Skills: [...],
       Persona: "专业、严谨、建设性反馈"
     }
   - SystemPrompt = "你是群成员"user-b"...（角色设定）..."

5. 选择SubAgent：
   - 根据user-b的Skills["code_review"]映射到opencode-coder
   - TargetAgent = opencode-coder

6. Executor执行（带角色）：
   - 在OpenCode请求中注入SystemPrompt作为system消息
   - 发送请求，包含：
     messages: [
       {role: "system", content: "你是群成员user-b..."},
       {role: "user", content: "帮我review这段代码"}
     ]

7. OpenCode以user-b的身份回复（体现专业工程师风格）

8. 结果返回：
   - 实时/批量返回给群聊（通过WebSocket广播给所有在线成员）
   - 日志保存到user-a的目录（因为user-a是session创建者）
     /logs/user-a/{sessionId}.log
   - 记录处理者信息：由user-b处理

场景2：user-a发送"谁来帮我写个文案？"（无明确@）

1. GroupRouter分析：
   - mentions为空
   - 提取关键词：["文案", "写"]
   - 所需技能：["writing"]

2. 筛选候选：
   - user-b：技能["code", "review"]，不匹配
   - user-c：技能["writing", "creative"]，匹配writing

3. 评分：
   - user-c技能匹配度=10，描述包含"文案策划"，负载=0
   - 选择user-c

4. 后续同流程，SystemPrompt注入user-c的角色信息
六、Go模块划分建议
plain
复制
deep-agent-core/
├── cmd/
│   └── server/
│       └── main.go                 # 入口：初始化所有模块，启动WebSocket服务
├── internal/
│   ├── core/
│   │   ├── agent_core.go           # AgentCore主结构
│   │   ├── message_router.go       # 消息路由（单聊/群聊分发）
│   │   └── session_manager.go      # 会话管理（Redis操作）
│   ├── handlers/
│   │   ├── single_chat.go          # SingleChatHandler
│   │   └── group_router.go         # GroupRouter（含智能路由算法）
│   ├── agents/
│   │   ├── subagent_pool.go        # SubAgentPool管理
│   │   ├── executor.go             # SubAgentExecutor
│   │   └── opencode_adapter.go     # OpenCode HTTP/SSE适配
│   ├── persistence/
│   │   ├── log_store.go            # LogPersistence（文件操作）
│   │   └── redis_store.go          # Redis操作封装
│   ├── config/
│   │   ├── manager.go              # ConfigManager
│   │   └── loader.go               # YAML配置文件加载
│   ├── websocket/
│   │   ├── gateway.go              # WebSocketGateway
│   │   ├── connection.go           # Connection结构
│   │   └── protocol.go             # 消息结构定义
│   └── types/
│       └── models.go               # 所有Struct定义（UserProfile, Session等）
├── configs/
│   ├── default.yaml                # 默认用户配置
│   └── agents.yaml                 # SubAgent配置
├── pkg/
│   └── utils/
│       ├── uuid.go                 # UUID生成
│       ├── tokenizer.go            # 简单分词
│       └── safety.go               # 路径安全检查
├── go.mod
└── Dockerfile
这个设计涵盖了所有你提到的需求点，并且针对Go语言的并发特性（goroutine、channel、sync包）做了适配。每个模块的职责、输入输出、核心算法都描述得很详细，可以直接交给AI编写具体代码。需要我对某个模块进一步细化吗？