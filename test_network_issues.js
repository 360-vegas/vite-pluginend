const axios = require('axios');

const API_BASE = 'http://localhost:8080';

// 问题链接测试用例
const problemLinks = [
    {
        url: 'https://trends.google.com/trends/explore?q=',
        name: 'Google Trends',
        expectedIssue: 'IPv6连接问题',
        expectedTimeout: 35
    },
    {
        url: 'https://www.simonandschuster.com/search/books/Category-Fiction/Bestsellers/_/N-g1hZi9v/Ntt-',
        name: 'Simon & Schuster',
        expectedIssue: '超时问题',
        expectedTimeout: 25
    },
    {
        url: 'https://www.baidu.com',
        name: '百度',
        expectedIssue: '应该正常',
        expectedTimeout: 18
    }
];

async function createTestLinks() {
    console.log('🔗 创建测试链接...');
    const linkIds = [];
    
    for (const testCase of problemLinks) {
        try {
            const response = await axios.post(`${API_BASE}/api/external-links`, {
                url: testCase.url,
                category: '网络问题测试',
                description: `${testCase.name} - ${testCase.expectedIssue}`,
                status: true
            }, {
                timeout: 30000,
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            
            console.log(`✅ ${testCase.name} 链接创建成功`);
            linkIds.push({
                id: response.data.id || response.data._id,
                ...testCase
            });
        } catch (error) {
            console.error(`❌ ${testCase.name} 链接创建失败:`, error.response?.data || error.message);
        }
    }
    
    return linkIds;
}

async function testIndividualLink(linkInfo) {
    console.log(`\n🔍 测试 ${linkInfo.name}...`);
    console.log(`🔗 URL: ${linkInfo.url}`);
    console.log(`⏰ 预期超时: ${linkInfo.expectedTimeout}秒`);
    console.log(`🎯 预期问题: ${linkInfo.expectedIssue}`);
    
    const startTime = Date.now();
    
    try {
        const response = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkInfo.id],
            all: false
        }, {
            timeout: 120000, // 2分钟超时
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.log(`⏱️ 实际检测时间: ${duration}秒`);
        
        if (response.data.results && response.data.results.length > 0) {
            const result = response.data.results[0];
            
            if (result.is_valid) {
                console.log('✅ 检测结果: 链接可用');
                console.log('📝 消息:', result.message);
            } else {
                console.log('❌ 检测结果: 链接不可用');
                console.log('🚫 错误信息:', result.error_message || result.message);
                
                // 分析错误类型
                const errorMsg = result.error_message || result.message || '';
                if (errorMsg.includes('IPv6')) {
                    console.log('🔧 状态: IPv6问题已被识别');
                } else if (errorMsg.includes('超时')) {
                    console.log('🔧 状态: 超时问题已被识别');
                } else if (errorMsg.includes('网络')) {
                    console.log('🔧 状态: 网络问题已被识别');
                }
            }
        }
        
        // 性能分析
        if (duration > linkInfo.expectedTimeout) {
            console.log(`⚠️ 警告: 检测时间(${duration}s)超过预期(${linkInfo.expectedTimeout}s)`);
        } else {
            console.log(`✅ 性能: 检测时间符合预期`);
        }
        
    } catch (error) {
        const endTime = Date.now();
        const duration = (endTime - startTime) / 1000;
        
        console.error(`❌ 检测失败！耗时: ${duration}秒`);
        console.error('错误详情:', error.response?.data || error.message);
    }
}

async function analyzeImprovements() {
    console.log('\n📊 优化效果分析:');
    console.log('🔧 已实施的优化:');
    console.log('  1. ✅ 强制使用IPv4，避免IPv6连接问题');
    console.log('  2. ✅ 扩展超时域名列表，包含主要大型网站');
    console.log('  3. ✅ 增加默认超时时间从12秒到18秒');
    console.log('  4. ✅ 改进错误信息格式化，提供更清晰的错误描述');
    console.log('  5. ✅ 增加TLS握手超时时间');
    console.log('  6. ✅ 独立数据库上下文，避免网络超时影响数据更新');
    
    console.log('\n🎯 预期改善:');
    console.log('  • Google Trends: IPv6问题解决，超时时间35秒');
    console.log('  • Simon & Schuster: 被识别为慢网站，超时时间25秒');
    console.log('  • 其他网站: 默认超时时间增加到18秒');
    console.log('  • 错误信息: 更清晰的问题描述');
}

async function main() {
    console.log('🧪 网络问题综合测试');
    console.log('🎯 测试目标: 验证IPv6和超时问题的解决效果');
    console.log('=' .repeat(60));
    
    // 分析优化效果
    await analyzeImprovements();
    
    // 创建测试链接
    const linkIds = await createTestLinks();
    
    if (linkIds.length === 0) {
        console.log('❌ 没有成功创建测试链接，无法继续测试');
        return;
    }
    
    // 逐个测试链接
    for (const linkInfo of linkIds) {
        await testIndividualLink(linkInfo);
        // 短暂延迟，避免请求过于频繁
        await new Promise(resolve => setTimeout(resolve, 2000));
    }
    
    console.log('\n' + '='.repeat(60));
    console.log('🏁 测试完成');
    console.log('📝 请观察错误信息的改善和IPv6问题的解决情况');
}

// 运行测试
main().catch(console.error); 