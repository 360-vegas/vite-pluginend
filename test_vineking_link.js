const axios = require('axios');

const API_BASE = 'http://localhost:8080';
const TEST_URL = 'https://thevineking.com/pages/search-results-page?q=';

async function createTestLink() {
    console.log('🔗 创建测试链接...');
    
    try {
        const response = await axios.post(`${API_BASE}/api/external-links`, {
            url: TEST_URL,
            category: '测试网站',
            description: 'The Vineking 搜索页面',
            status: true
        }, {
            timeout: 30000,
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        console.log('✅ 链接创建成功:', response.data);
        return response.data.id || response.data._id;
    } catch (error) {
        console.error('❌ 创建链接失败:', error.response?.data || error.message);
        return null;
    }
}

async function testSingleLink(linkId) {
    console.log('\n🔍 测试单个链接检测...');
    const startTime = Date.now();
    
    try {
        const response = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkId],
            all: false
        }, {
            timeout: 60000,
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.log(`⏱️ 检测耗时: ${duration}秒`);
        console.log('📊 检测结果:', JSON.stringify(response.data, null, 2));
        
        if (response.data.results && response.data.results.length > 0) {
            const result = response.data.results[0];
            if (result.is_valid) {
                console.log('✅ 检测结果: 链接可用');
                console.log('📝 消息:', result.message);
            } else {
                console.log('❌ 检测结果: 链接不可用');
                console.log('🚫 错误信息:', result.error_message || result.message);
            }
        }
        
    } catch (error) {
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.error(`❌ 检测失败！耗时: ${duration}秒`);
        console.error('错误详情:', error.response?.data || error.message);
    }
}

async function testDirectAccess() {
    console.log('\n🌐 直接访问测试（作为对比）...');
    
    try {
        const response = await axios.get(TEST_URL, {
            timeout: 10000,
            headers: {
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
            }
        });
        
        console.log('✅ 直接访问成功');
        console.log('📊 状态码:', response.status);
        console.log('📄 内容长度:', response.data.length);
        console.log('🏷️ 内容类型:', response.headers['content-type']);
        
    } catch (error) {
        console.error('❌ 直接访问失败:', error.response?.status || error.message);
    }
}

async function main() {
    console.log('🧪 The Vineking 链接检测测试');
    console.log('🔗 测试URL:', TEST_URL);
    console.log('=' .repeat(60));
    
    // 首先直接访问测试
    await testDirectAccess();
    
    // 创建测试链接
    const linkId = await createTestLink();
    
    if (linkId) {
        // 测试我们的检测系统
        await testSingleLink(linkId);
    }
    
    console.log('\n' + '='.repeat(60));
    console.log('🏁 测试完成');
}

// 运行测试
main().catch(console.error); 