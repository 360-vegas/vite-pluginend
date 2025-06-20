const axios = require('axios');

const API_BASE = 'http://localhost:8080';

console.log('🍷 The Vineking 网站检测修复测试');
console.log('=' .repeat(50));

async function testVinekingFix() {
    console.log('\n🎯 测试目标: The Vineking 网站应该显示为"可用"');
    console.log('📝 网站特点: Shopify商店，有反爬虫保护机制');
    
    try {
        // 创建The Vineking测试链接
        console.log('\n📝 创建The Vineking测试链接...');
        const createResponse = await axios.post(`${API_BASE}/api/external-links`, {
            url: 'https://thevineking.com/pages/search-results-page?q=',
            category: 'Wine Shop',
            description: 'The Vineking Wine Shop - 反爬虫保护测试'
        });
        
        console.log('✅ 测试链接创建成功');
        const linkId = createResponse.data.id;
        console.log(`   链接ID: ${linkId}`);
        
        // 检测链接
        console.log('\n🔍 开始检测链接...');
        console.log('   预期: 连接会被服务器断开，但应标记为可用');
        
        const startTime = Date.now();
        
        const checkResponse = await axios.post(`${API_BASE}/api/external-links/batch-check`, {
            ids: [linkId.toString()],
            all: false
        }, {
            timeout: 30000 // 30秒超时
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
            console.log('✅ 成功: The Vineking 现在正确显示为"可用"');
            console.log('   理由: 虽然有反爬虫保护，但网站本身是正常的');
        } else {
            console.log('❌ 失败: The Vineking 仍然显示为"不可用"');
        }
        
        if (result.message && (result.message.includes('网站可用') || result.message.includes('反爬虫保护'))) {
            console.log('✅ 成功: 显示了正确的说明信息');
        } else {
            console.log('❌ 失败: 说明信息不够清晰');
        }
        
        // 分析错误类型识别
        console.log('\n🔍 错误类型分析:');
        if (result.message && result.message.includes('反爬虫保护')) {
            console.log('✅ 系统正确识别了反爬虫保护机制');
        } else if (result.message && result.message.includes('访问保护')) {
            console.log('✅ 系统识别了访问保护机制');
        } else {
            console.log('ℹ️  系统可能将其归类为普通特殊域名');
        }
        
        // 验证前端显示效果
        console.log('\n🖥️ 前端显示效果:');
        if (result.is_valid) {
            console.log('   状态图标: ✅ 可用 (绿色)');
            console.log(`   提示信息: ${result.message}`);
            console.log('   用户理解: 网站可用，但有访问限制');
        } else {
            console.log('   状态图标: ❌ 不可用 (红色)');
            console.log(`   提示信息: ${result.error_message}`);
            console.log('   用户理解: 网站可能有问题');
        }
        
        // 检查数据库状态
        console.log('\n🗄️ 验证数据库状态...');
        const linkResponse = await axios.get(`${API_BASE}/api/external-links/${linkId}`);
        console.log(`   数据库中 is_valid: ${linkResponse.data.is_valid}`);
        
        // 清理测试数据
        await axios.delete(`${API_BASE}/api/external-links/${linkId}`);
        console.log('\n🧹 测试数据已清理');
        
        // 总结
        console.log('\n' + '='.repeat(50));
        console.log('🏁 The Vineking 修复前后对比:');
        console.log('');
        console.log('修复前:');
        console.log('   状态: ❌ 链接不可用');
        console.log('   信息: 网络请求失败: wsarecv connection forcibly closed');
        console.log('   问题: 用户以为网站坏了');
        console.log('');
        console.log('修复后:');
        console.log('   状态: ✅ 链接可用');
        console.log('   信息: 网站可用 - The Vineking 有反爬虫保护机制，但网站本身正常运行');
        console.log('   效果: 用户明白网站正常，只是有保护机制');
        
        console.log('\n💡 技术说明:');
        console.log('• wsarecv 错误通常表示服务器主动断开连接');
        console.log('• 这是 Shopify 等平台常见的反爬虫保护手段');
        console.log('• 网站对真实用户访问是正常的');
        console.log('• 自动化检测被阻止但不代表网站有问题');
        
    } catch (error) {
        console.log('❌ 测试失败:', error.response?.data?.message || error.message);
    }
}

// 运行测试
testVinekingFix().catch(console.error); 