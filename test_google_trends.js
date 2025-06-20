const axios = require('axios');

const API_BASE = 'http://localhost:8080';
const TEST_URL = 'https://trends.google.com/trends/explore?q=';

async function createTestLink() {
    console.log('🔗 创建Google Trends测试链接...');
    
    try {
        const response = await axios.post(`${API_BASE}/api/external-links`, {
            url: TEST_URL,
            category: '测试网站',
            description: 'Google Trends 探索页面',
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
    console.log('\n🔍 测试Google Trends链接检测...');
    console.log('⏰ 注意：Google Trends可能需要较长时间响应...');
    const startTime = Date.now();
    
    try {
        const response = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkId],
            all: false
        }, {
            timeout: 120000, // 2分钟超时，适应Google Trends的响应时间
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
                console.log('✅ 检测结果: Google Trends链接可用');
                console.log('📝 消息:', result.message);
            } else {
                console.log('❌ 检测结果: Google Trends链接不可用');
                console.log('🚫 错误信息:', result.error_message || result.message);
                
                // 分析可能的原因
                if (result.error_message && result.error_message.includes('timeout')) {
                    console.log('💡 提示: Google Trends对自动化请求有严格限制，超时是常见现象');
                } else if (result.error_message && result.error_message.includes('403')) {
                    console.log('💡 提示: Google Trends可能拒绝了机器人请求');
                } else if (result.error_message && result.error_message.includes('429')) {
                    console.log('💡 提示: 请求过于频繁，被Google限流');
                }
            }
        }
        
    } catch (error) {
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.error(`❌ 检测失败！耗时: ${duration}秒`);
        console.error('错误详情:', error.response?.data || error.message);
        
        if (error.code === 'ECONNABORTED') {
            console.error('💡 前端请求超时，可能需要更长的等待时间');
        }
    }
}

async function testDirectAccess() {
    console.log('\n🌐 直接访问Google Trends测试（作为对比）...');
    
    try {
        const response = await axios.get(TEST_URL, {
            timeout: 30000, // 30秒超时
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
        
        if (error.code === 'ECONNABORTED') {
            console.error('💡 直接访问也超时，说明Google Trends确实响应较慢');
        } else if (error.response?.status === 403) {
            console.error('💡 Google Trends拒绝了请求，可能是反爬虫机制');
        }
    }
}

async function analyzeIssue() {
    console.log('\n🔬 Google Trends检测问题分析:');
    console.log('1. ⏰ 响应时间: Google Trends通常需要15-30秒响应');
    console.log('2. 🤖 反爬虫: Google有严格的机器人检测机制');
    console.log('3. 🌍 地理位置: 可能根据IP地址返回不同内容');
    console.log('4. 🔑 认证: 某些功能可能需要登录');
    console.log('5. 📊 动态内容: 页面可能通过JavaScript动态加载');
    
    console.log('\n💡 建议的解决方案:');
    console.log('- 🕐 增加超时时间到30-60秒');
    console.log('- 🎭 使用更完整的浏览器User-Agent');
    console.log('- 🔄 对于特定域名使用不同的检测策略');
    console.log('- ⚠️ 考虑将某些大型网站标记为"特殊处理"');
}

async function main() {
    console.log('🧪 Google Trends 链接检测专项测试');
    console.log('🔗 测试URL:', TEST_URL);
    console.log('=' .repeat(60));
    
    // 分析问题
    await analyzeIssue();
    
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
    console.log('📝 注意: Google Trends的检测结果可能不稳定，这是正常现象');
}

// 运行测试
main().catch(console.error); 