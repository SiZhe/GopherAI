module.exports = {
    devServer: {
        port: 8080,
        // 1. 允许所有 Host 头（解决 Invalid Host header 错误）
        allowedHosts: 'all',
        // 2. 添加响应头，跳过 ngrok 警告页面
        headers: {
            'ngrok-skip-browser-warning': '1'
        },
        // 保留你原有的 proxy 配置（不变）
        proxy: {
            '/api': {
                target: 'http://localhost:8888',
                changeOrigin: true,
                pathRewrite: {
                    '^/api': '/api/v1'
                }
            }
        }
    }
}