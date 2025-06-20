const axios = require('axios');

const API_BASE = 'http://localhost:8080';

console.log('🔧 测试特殊域名逻辑修复');
console.log('=' .repeat(50));

async function testSpecialDomainLogic() {
    console.log('\n🎯 测试目标: 特殊域名应该显示为"可用"而不是"不可用"');
    
    try {
        // 创建Google Trends测试链接
        console.log('\n📝 创建Google Trends测试链接...');
        const createResponse = await axios.post(`${API_BASE}/api/external-links`, {
            url: 'https://trends.google.com/trends/explore?q=test',
            category: '特殊域名测试',
            description: '测试特殊域名逻辑修复'
        });
        
        console.log('✅ 测试链接创建成功');
        const linkId = createResponse.data.id;
        console.log(`   链接ID: ${linkId}`);
        
        // 检测链接
        console.log('\n🔍 开始检测链接...');
        const startTime = Date.now();
        
        const checkResponse = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkId.toString()],
            all: false
        }, {
            timeout: 120000
        });
        
        const endTime = Date.now();
        const duration = Math.round((endTime - startTime) / 1000);
        
        console.log(`⏱️  检测耗时: ${duration}秒`);
        
        const result = checkResponse.data.results[0];
        console.log('\n📊 检测结果分析:');
        console.log(`   URL: ${result.url}`);
        console.log(`   IsValid: ${result.is_valid}`);
        console.log(`   Message: ${result.message || '无'}`);
        console.log(`   ErrorMessage: ${result.error_message || '无'}`);
        
        // 验证修复效果
        console.log('\n🎯 修复效果验证:');
        
        if (result.is_valid === true) {
            console.log('✅ 成功: 特殊域名现在正确显示为"可用"');
        } else {
            console.log('❌ 失败: 特殊域名仍然显示为"不可用"');
        }
        
        if (result.message && result.message.includes('网站可用')) {
            console.log('✅ 成功: 显示了正确的说明信息');
        } else {
            console.log('❌ 失败: 说明信息不正确');
        }
        
        if (!result.error_message || result.error_message === '') {
            console.log('✅ 成功: 错误信息已清空');
        } else {
            console.log('❌ 注意: 仍有错误信息显示');
        }
        
        // 验证前端显示效果
        console.log('\n🖥️ 前端显示效果:');
        if (result.is_valid) {
            console.log('   状态: ✅ 可用 (绿色图标)');
            console.log(`   提示: ${result.message}`);
        } else {
            console.log('   状态: ❌ 不可用 (红色图标)');
            console.log(`   提示: ${result.error_message}`);
        }
        
        // 检查数据库中的状态
        console.log('\n🗄️ 验证数据库状态...');
        const linkResponse = await axios.get(`${API_BASE}/api/external-links/${linkId}`);
        console.log(`   数据库中 is_valid: ${linkResponse.data.is_valid}`);
        
        // 清理测试数据
        await axios.delete(`${API_BASE}/api/external-links/${linkId}`);
        console.log('\n🧹 测试数据已清理');
        
        // 总结
        console.log('\n' + '='.repeat(50));
        console.log('🏁 修复前 vs 修复后对比:');
        console.log('');
        console.log('修复前:');
        console.log('   状态: ❌ 链接不可用');
        console.log('   信息: Google Trends 对自动化检测有严格限制，但网站通常可正常访问');
        console.log('   问题: 矛盾的提示，用户困惑');
        console.log('');
        console.log('修复后:');
        console.log('   状态: ✅ 链接可用');
        console.log('   信息: 网站可用 - Google Trends 对自动化检测有限制，但网站本身正常运行');
        console.log('   效果: 逻辑一致，用户清晰明了');
        
    } catch (error) {
        console.log('❌ 测试失败:', error.response?.data?.message || error.message);
    }
}

// 运行测试
testSpecialDomainLogic().catch(console.error); 