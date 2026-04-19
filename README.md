# GopherAI
这是一个AI应用服务平台。你可以使用多种大语言模型进行对话，并可上传文件建立专属于对话的知识库。

## 用户注册
1. 填写个人邮箱后，点击 **“发送验证码”**，系统会将验证码发至个人邮箱(验证码有效时间**5分钟**);
2. 注册成功后，系统会将 **专属id** 发至个人邮箱，你可以使用 **个人邮箱** 或 **专属id** 进行登录。

![用户注册](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E7%94%A8%E6%88%B7%E6%B3%A8%E5%86%8C.png)

## AI会话
1. 首先，点击左侧会话列表下的 **“+新对话”** 创建对话;
2. 接下来，选择会话模型，包括 **豆包_Seed_2.0** 以及 **DeepSeek_v3.2**，即可开始会话;
3. 如果某次会话的历史消息未加载完全，可点击 **“同步历史数据”**。

![创建聊天](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E5%88%9B%E5%BB%BA%E8%81%8A%E5%A4%A9.png)

> **会话上下文记忆** （模型知道用户发送多次“你好”）
>
![上下文记忆](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E4%B8%8A%E4%B8%8B%E6%96%87%E8%AE%B0%E5%BF%86.png)

## 流式会话模式
大语言模型不再一股脑的输出 **整段消息**。

![流式输出](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E6%B5%81%E5%BC%8F%E8%BE%93%E5%87%BA.png)

## 上传文件构建知识库-检索增强生成(RAG)
你可以上传专属于该对话的文件(目前仅限于.md)以构建知识库。
> **未构建知识库**
> 
![未构建知识库](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E6%9C%AA%E6%9E%84%E5%BB%BA%E7%9F%A5%E8%AF%86%E5%BA%93.png)

> **上传文件**
> 
![上传文件](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E4%B8%8A%E4%BC%A0%E6%96%87%E4%BB%B6.png)

> **已构建知识库**
> 
![已构建知识库](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E5%B7%B2%E6%9E%84%E5%BB%BA%E7%9F%A5%E8%AF%86%E5%BA%93.png)

## ReAct Agent for RAG
通过ReAct Agent智能管理RAG。
> **未构建 ReAct Agent**: 每次提问都会使用RAG
>
![NullRagReActAgent](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/NullRagReActAgent.png)

> **构建 ReAct Agent for RAG**: Agent根据用户问题智能判断 **本次问题是否需要RAG**
>
![RagReActAgent](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/RagReActAgent.png)

> **ReAct Agent 可使用以下3个Tool:**    
1.检查用户在当前会话中是否有上传的文档    
2.使用 RAG 检索文档    
3.使用小型 LLM 判断是否需要使用 RAG 技术

![ReActTools](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/ReActTools.png)


## 设备管理
你可以查看你登录的设备；

![设备管理](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E8%AE%BE%E5%A4%87%E7%AE%A1%E7%90%86.png)

并可以控制其登录状态。

![设备管理列表](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E8%AE%BE%E5%A4%87%E7%AE%A1%E7%90%86%E5%88%97%E8%A1%A8.png)

# 技术细节

## 1.JWT与设备管理

1. **jwt结构**：base64(header + payload) + 签名;  
    token = base64(header + payload) + HS256(base64(header + payload)  
2. **jwt优点**：无状态；服务端不用存储信息;  
   **jwt缺点**：一旦发出，无法主动作废  
3. **实现jwt的安全传输**：  
    (1)防止监听：将jwt的token保存在cookie中，并通过“secure=true"保证只通过https加密传输,"httpOnly=true”不让前端看到，避免泄漏;  
    (2)防止冒充：在payload中保存device_info信息(device_info = ip + user-agent),每次访问都验证当前 ip + user-agent 是否等于 token的device_info;  
    (3)兜底：Access Token + Refresh Token;AT有效期15分钟，保证安全；RT有效期7天，保证无感续期。  
    (4)黑名单：如果发现token被盗，将token加入黑名单，该token不允访问  
4. **AUTH**:  
    1.从cookie中获取AT;  
    2.检查是否在黑名单中;  
    3.解析AT是否有效; 若AT过期，前端用RT更新AT;  
    4.判断device_info是否一致;  
5. **设备管理**:  
    1.登录后，将username,device_info,at,rt存入数据库;  
    2.下线，从数据库中删除对应设备，并将at，rt加入jwt黑名单以实现禁止访问;  

## 2.流式输出

SSE (Server-Sent Events) :一种基于 HTTP 协议的 **服务器单向推送** 技术，允许服务器主动、持续地将实时数据发送到客户端  
每条消息以 \n\n 分隔

**实现**  
c.Header("Content-Type", "text/event-stream")  
c.Header("Cache-Control", "no-cache")  
c.Header("Connection", "keep-alive")  
c.Header("Access-Control-Allow-Origin", "*")

c.writer.Write([]byte("data: " + msg + "\n\n"))

## 3.检索增强生成(RAG)

**准备**  
    (1)获取embedding model -> 将文字转化为多维向量  
    (2)向量数据库: milvus
1. 将上传文件保存至服务端，再将文件通过eino框架中 **splitter.Transformer** 进行分割得到 **[]schema.Document**
2. 通过 **indexer.Store** 将 **[]schema.Document** 保存至向量数据库
3. 判断判断该对话中用户是否上传文件判断是否要进行RAG
4. 通过 **retriever.Retrieve** 对用户的最后一条信息进行检索
5. 构建提示词并加入到对话上下文中


## 4.什么是ReAct？

ReAct 编排是一种让大型语言模型（LLM）更智能地执行任务的策略，它的核心思想是让 LLM 像人一样**思考（Reason）和行动（Act）**。想象一下，你正在解决一个复杂的问题：你不会一下子就得出答案，而是会先思考问题，然后采取一些行动（比如查找信息、执行计算），根据行动的结果再进一步思考，如此循环，直到问题解决。ReAct 就是让 LLM 模拟这个过程。

### ReAct 的核心原理

ReAct 的名字来源于两个关键部分：

- **Reason（思考）**: LLM 会生成**思维链**（Chain-of-Thought），这是一种详细的推理过程，包括它对问题的理解、如何拆解问题、以及接下来打算做什么。这就像一个人在心里默念自己的思考过程。
- **Act（行动）**: LLM 会根据它的思考，调用外部工具或执行特定操作。这些“工具”可以是：
    - **搜索引擎**：查找最新信息或特定事实。
    - **计算器**：执行数学运算。
    - **代码解释器**：运行代码来处理数据或验证逻辑。
    - **API 接口**：与外部系统交互，比如查询天气、发送邮件等。

### ReAct 的工作流程

ReAct 的工作流程可以概括为以下步骤：

1. **接收任务**：用户给 LLM 一个任务或问题。
2. **思考（Reason）**：LLM 首先会思考这个问题，生成一个计划或下一步的行动方针。这个思考过程是可见的，比如它会说“我需要先查找XX信息”或者“这个问题可以分解为A和B”。
3. **行动（Act）**：根据思考结果，LLM 会选择一个合适的工具并执行操作。例如，如果它需要查找信息，它会调用搜索引擎并生成搜索查询。
4. **观察（Observe）**：LLM 会接收到工具执行后的结果（例如，搜索结果、计算结果等）。
5. **循环**：LLM 会根据观察到的结果，再次回到**思考**步骤。它会分析结果，决定下一步是继续使用工具、修改之前的计划，还是已经得到了最终答案并准备输出。
6. **输出答案**：当 LLM 认为任务已经完成时，它会生成最终的答案。

![ReAct](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/ReAct.png)

### 为什么 ReAct 如此有效？

ReAct 编排之所以强大，主要有以下几个原因：

- **提高准确性**：通过调用外部工具获取最新、准确的信息，LLM 可以避免“幻觉”（即生成不真实的信息）。
- **处理复杂任务**：将大任务分解为小步骤，并通过循环的“思考-行动-观察”过程逐步解决，使 LLM 能够处理更复杂、需要多步推理和外部协作的任务。
- **可解释性**：LLM 生成的思维链使得它的决策过程更加透明，我们可以看到它是如何思考和解决问题的，这有助于调试和理解。
- **适应动态环境**：LLM 可以根据外部工具的反馈动态调整其行为和策略。

