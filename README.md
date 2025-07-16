# webapp-tieba

1. 完成时间：2025年5月16日  
2. 优化时间：2025年7月 至 2025年8月  
3. 优化内容：

- [ ] 对 bluebell 后端代码进行全体重构优化  
  包括：

  1. **优化 access token 和 refresh token 机制**

     **核心逻辑：**

     - **客户端正常请求：**  
       客户端（如前端应用）在每次发起需要授权的 API 请求时，在请求头中携带 access_token：
       ```js
       axios.defaults.headers.common['Authorization'] = 'Bearer ' + accessToken;
       ```

     - **服务端 JWT 认证中间件：**  
       在 gin 中使用 JWT 中间件对接口进行保护。当 access_token 过期时，不应直接返回 `401 Unauthorized`，而应返回约定好的错误码：
       ```json
       { "code": 1007, "msg": "access token expired" }
       ```

     - **客户端响应拦截器（核心）：**  
       前端需设置 axios 全局响应拦截器，判断是否返回了 code = 1007 的响应：
       ```js
       axios.interceptors.response.use(
         response => response,
         async error => {
           if (error.response?.data?.code === 1007) {
             // 触发刷新 token 逻辑
           }
           return Promise.reject(error);
         }
       );
       ```

     **自动刷新逻辑：**

     - 暂停其他请求：  
       拦截后将所有请求放入队列，等待刷新完成再重新发起

     - 调用刷新接口：  
       使用本地 refresh_token 请求 `/refreshtoken`

     - 刷新成功：
       1. 获取新的 access_token 和 refresh_token
       2. 更新本地存储和 axios 默认请求头
       3. 重新发送失败的请求
       4. 释放并执行等待队列中的请求

     - 刷新失败：
       1. 清除本地凭证信息
       2. 跳转至登录页面，要求重新登录

  2. **增加按关键词模糊匹配搜索帖子的功能**
  3. **增加侧边栏实时热点新闻浏览功能**