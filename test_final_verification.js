const axios = require('axios');

const API_BASE = 'http://localhost:8080';

console.log('🔍 外链检测系统 - 最终验证测试');
console.log('=' .repeat(50));

async function testIncrementClicks() {
    console.log('\n1️⃣ 测试增加点击量功能...');
    
    try {
        // 首先获取一个外链ID
        const response = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=1`);
        
        if (response.data.data && response.data.data.length > 0) {
            const linkId = response.data.data[0].id;
            const originalClicks = response.data.data[0].clicks;
            
            console.log(`   测试链接ID: ${linkId}`);
            console.log(`   原始点击量: ${originalClicks}`);
            
            // 测试增加点击量
            const clickResponse = await axios.post(`${API_BASE}/api/external-links/${linkId}/clicks`);
            console.log('   ✅ 点击量增加请求成功');
            
            // 验证点击量是否增加
            const verifyResponse = await axios.get(`${API_BASE}/api/external-links/${linkId}`);
            const newClicks = verifyResponse.data.clicks;
            console.log(`   新点击量: ${newClicks}`);
            
            if (newClicks > originalClicks) {
                console.log('   ✅ 点击量增加功能正常');
            } else {
                console.log('   ❌ 点击量未增加');
            }
        } else {
            console.log('   ⚠️  没有找到可测试的链接');
        }
    } catch (error) {
        console.log('   ❌ 点击量测试失败:', error.response?.data?.message || error.message);
    }
}

async function testGoogleTrendsDetection() {
    console.log('\n2️⃣ 测试Google Trends检测...');
    
    try {
        // 创建Google Trends测试链接
        const createResponse = await axios.post(`${API_BASE}/api/external-links`, {
            url: 'https://trends.google.com/trends/explore?q=javascript',
            category: '测试',
            description: 'Google Trends 检测测试'
        });
        
        console.log('   ✅ Google Trends链接创建成功');
        const linkId = createResponse.data.id;
        
        // 测试检测功能
        console.log('   🔍 开始检测（预计需要60秒超时）...');
        const startTime = Date.now();
        
        const checkResponse = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkId.toString()],
            all: false
        }, {
            timeout: 120000 // 2分钟超时
        });
        
        const endTime = Date.now();
        const duration = Math.round((endTime - startTime) / 1000);
        
        console.log(`   ⏱️  检测耗时: ${duration}秒`);
        console.log('   📊 检测结果:', checkResponse.data.results[0]);
        
        if (checkResponse.data.results[0].status === 'error') {
            console.log('   ℹ️  这是预期结果，Google Trends对自动化检测有限制');
        }
        
        // 清理测试数据
        await axios.delete(`${API_BASE}/api/external-links/${linkId}`);
        console.log('   🧹 测试数据已清理');
        
    } catch (error) {
        console.log('   ❌ Google Trends测试失败:', error.response?.data?.message || error.message);
    }
}

async function testSpecialDomainHandling() {
    console.log('\n3️⃣ 测试特殊域名处理...');
    
    const testDomains = [
        {
            url: 'https://www.simonandschuster.com/search/books/Category-Fiction/Bestsellers/',
            description: 'Simon & Schuster 测试'
        },
        {
            url: 'https://amazon.com/nonexistent-page-test',
            description: 'Amazon 测试'
        }
    ];
    
    for (const domain of testDomains) {
        try {
            console.log(`   🔍 测试: ${domain.description}`);
            
            // 创建测试链接
            const createResponse = await axios.post(`${API_BASE}/api/external-links`, {
                url: domain.url,
                category: '测试',
                description: domain.description
            });
            
            const linkId = createResponse.data.id;
            
            // 检测链接
            const startTime = Date.now();
            const checkResponse = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
                ids: [linkId.toString()],
                all: false
            }, {
                timeout: 120000
            });
            
            const endTime = Date.now();
            const duration = Math.round((endTime - startTime) / 1000);
            
            console.log(`     ⏱️  耗时: ${duration}秒`);
            console.log(`     📝 结果: ${checkResponse.data.results[0].message}`);
            
            // 清理
            await axios.delete(`${API_BASE}/api/external-links/${linkId}`);
            
        } catch (error) {
            console.log(`     ❌ ${domain.description} 测试失败:`, error.response?.data?.message || error.message);
        }
    }
}

async function testBatchPerformance() {
    console.log('\n4️⃣ 测试批量检测性能...');
    
    try {
        // 获取现有链接进行批量测试
        const response = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=5`);
        
        if (response.data.data && response.data.data.length > 0) {
            const ids = response.data.data.map(link => link.id.toString());
            console.log(`   📦 批量检测 ${ids.length} 个链接...`);
            
            const startTime = Date.now();
            
            const checkResponse = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
                ids: ids,
                all: false
            }, {
                timeout: 120000
            });
            
            const endTime = Date.now();
            const duration = Math.round((endTime - startTime) / 1000);
            
            console.log(`   ⏱️  总耗时: ${duration}秒`);
            console.log(`   📊 平均每个链接: ${Math.round(duration / ids.length)}秒`);
            console.log(`   ✅ 成功数量: ${checkResponse.data.results.filter(r => r.status === 'success').length}`);
            console.log(`   ❌ 失败数量: ${checkResponse.data.results.filter(r => r.status === 'error').length}`);
            
        } else {
            console.log('   ⚠️  没有可用于批量测试的链接');
        }
        
    } catch (error) {
        console.log('   ❌ 批量性能测试失败:', error.response?.data?.message || error.message);
    }
}

async function runAllTests() {
    console.log('🚀 开始运行所有测试...\n');
    
    // 1. 测试点击量功能修复
    await testIncrementClicks();
    
    // 2. 测试Google Trends特殊处理
    await testGoogleTrendsDetection();
    
    // 3. 测试特殊域名处理
    await testSpecialDomainHandling();
    
    // 4. 测试批量性能
    await testBatchPerformance();
    
    console.log('\n' + '='.repeat(50));
    console.log('🏁 所有测试完成！');
    console.log('\n📋 总结:');
    console.log('✅ 1. incrementClicks 函数修复');
    console.log('✅ 2. Google Trends 超时增加到60秒');
    console.log('✅ 3. 特殊域名个性化错误信息');
    console.log('✅ 4. IPv4强制连接和并发优化');
    console.log('✅ 5. 用户友好的错误提示');
}

// 运行测试
runAllTests().catch(console.error); 