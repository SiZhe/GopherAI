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

## 上传文件
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

## 流式会话模式
大语言模型不再一股脑的输出 **整段消息**。

![流式输出](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E6%B5%81%E5%BC%8F%E8%BE%93%E5%87%BA.png)

## 设备管理
你可以查看你登录的设备；

![设备管理](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E8%AE%BE%E5%A4%87%E7%AE%A1%E7%90%86.png)

并可以控制其登录状态。

![设备管理列表](https://github.com/SiZhe/readmeImage/blob/main/GopherAI/%E8%AE%BE%E5%A4%87%E7%AE%A1%E7%90%86%E5%88%97%E8%A1%A8.png)

## 技术细节

### 1.JWT与设备管理

1.jwt结构：base64(header + payload) + 签名
          token = base64(header + payload) + HS256(base64(header + payload)
2.jwt优点：无状态；服务端不用存储信息;
  jwt缺点：一旦发出，无法主动作废
3.实现jwt的安全传输：
    (1)防止监听：将jwt的token保存在cookie中，并通过“secure=true"保证只通过https加密传输,"httpOnly=true”不让前端看到，避免泄漏;
    (2)防止冒充：在payload中保存device_info信息(device_info = ip + user-agent),每次访问都验证当前 ip + user-agent 是否等于 token的device_info;
    (3)兜底：Access Token + Refresh Token;AT有效期15分钟，保证安全；RT有效期7天，保证无感续期。
    (4)黑名单：如果发现token被盗，将token加入黑名单，该token不允访问
AUTH：1.从cookie中获取AT;
     2.检查是否在黑名单中;
     3.解析AT是否有效; 若AT过期，前端用RT更新AT;
     4.判断device_info是否一致;
设备管理: 1.登录后，将username,device_info,at,rt存入数据库;
         2.下线，从数据库中删除对应设备，并将at，rt加入jwt黑名单以实现禁止访问;