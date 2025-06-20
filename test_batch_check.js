const axios = require('axios');

const API_BASE = 'http://localhost:8080';

async function testBatchCheck() {
    console.log('开始测试批量检测API...');
    const startTime = Date.now();
    
    try {
        console.log('请求URL:', `${API_BASE}/api/external-links/batch-check`);
        console.log('请求参数:', { ids: [], all: true });
        
        const response = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [],
            all: true
        }, {
            timeout: 120000, // 2分钟超时
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.log(`✅ 批量检测成功！耗时: ${duration}秒`);
        console.log('响应状态:', response.status);
        console.log('检测结果数量:', response.data.results ? response.data.results.length : 0);
        
        if (response.data.results) {
            const valid = response.data.results.filter(r => r.is_valid).length;
            const invalid = response.data.results.filter(r => !r.is_valid).length;
            console.log(`有效链接: ${valid}个，无效链接: ${invalid}个`);
        }
        
    } catch (error) {
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.error(`❌ 批量检测失败！耗时: ${duration}秒`);
        
        if (error.code === 'ECONNABORTED') {
            console.error('错误类型: 请求超时');
        } else if (error.response) {
            console.error('响应状态:', error.response.status);
            console.error('响应数据:', error.response.data);
        } else if (error.request) {
            console.error('网络错误: 无法连接到服务器');
        } else {
            console.error('其他错误:', error.message);
        }
    }
}

// 运行测试
testBatchCheck(); 