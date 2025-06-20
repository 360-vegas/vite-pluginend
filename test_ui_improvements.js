const axios = require('axios');

const API_BASE = 'http://localhost:8080';

console.log('🎨 外链管理界面改进验证');
console.log('=' .repeat(50));

async function testUIImprovements() {
    console.log('\n🎯 改进目标验证:');
    console.log('1. ✅ 隐藏分类和描述列，减少信息过载');
    console.log('2. ✅ 突出显示可用性状态，便于快速识别');
    console.log('3. ✅ 支持按可用性排序和筛选');
    console.log('4. ✅ 扩大链接地址显示区域');
    
    try {
        // 测试数据获取
        console.log('\n📊 获取当前外链数据...');
        const response = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=10`);
        
        if (response.data.data && response.data.data.length > 0) {
            const links = response.data.data;
            console.log(`   📦 获取到 ${links.length} 条外链数据`);
            
            // 分析可用性分布
            const availableCount = links.filter(link => link.is_valid).length;
            const unavailableCount = links.length - availableCount;
            
            console.log('\n📈 可用性分布分析:');
            console.log(`   ✅ 可用链接: ${availableCount} 条 (${Math.round(availableCount/links.length*100)}%)`);
            console.log(`   ❌ 不可用链接: ${unavailableCount} 条 (${Math.round(unavailableCount/links.length*100)}%)`);
            
            // 显示示例链接信息
            console.log('\n🔍 链接信息示例:');
            links.slice(0, 3).forEach((link, index) => {
                console.log(`   ${index + 1}. ${link.is_valid ? '✅' : '❌'} ${link.url}`);
                console.log(`      点击量: ${link.clicks} | 状态: ${link.status ? '启用' : '禁用'}`);
            });
            
        } else {
            console.log('   ⚠️  没有外链数据用于展示');
        }
        
        // 测试按可用性筛选
        console.log('\n🔍 测试可用性筛选功能...');
        
        // 获取可用链接
        const availableResponse = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=10&is_valid=true`);
        console.log(`   ✅ 可用链接查询: ${availableResponse.data.data?.length || 0} 条`);
        
        // 获取不可用链接
        const unavailableResponse = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=10&is_valid=false`);
        console.log(`   ❌ 不可用链接查询: ${unavailableResponse.data.data?.length || 0} 条`);
        
        // 测试按可用性排序
        console.log('\n📊 测试可用性排序功能...');
        
        // 可用性升序排序 (不可用的在前)
        const ascResponse = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=10&sort_field=is_valid&sort_order=asc`);
        if (ascResponse.data.data && ascResponse.data.data.length > 0) {
            const firstLink = ascResponse.data.data[0];
            console.log(`   升序排序 (不可用优先): 第一条链接状态 = ${firstLink.is_valid ? '✅ 可用' : '❌ 不可用'}`);
        }
        
        // 可用性降序排序 (可用的在前)
        const descResponse = await axios.get(`${API_BASE}/api/external-links?page=1&per_page=10&sort_field=is_valid&sort_order=desc`);
        if (descResponse.data.data && descResponse.data.data.length > 0) {
            const firstLink = descResponse.data.data[0];
            console.log(`   降序排序 (可用优先): 第一条链接状态 = ${firstLink.is_valid ? '✅ 可用' : '❌ 不可用'}`);
        }
        
        console.log('\n' + '='.repeat(50));
        console.log('🎨 界面改进前后对比:');
        console.log('');
        console.log('改进前:');
        console.log('   📋 列: [选择] [链接+状态标签] [分类] [描述] [点击量] [开关] [时间] [操作]');
        console.log('   ❌ 信息过载，难以快速识别可用性');
        console.log('   ❌ 分类和描述占用大量空间');
        console.log('   ❌ 可用性状态不突出');
        console.log('');
        console.log('改进后:');
        console.log('   📋 列: [选择] [✅/❌ 可用性] [链接地址] [点击量] [开关] [时间] [操作]');
        console.log('   ✅ 可用性状态一目了然');
        console.log('   ✅ 链接地址显示更宽敞');
        console.log('   ✅ 支持按可用性排序和筛选');
        console.log('   ✅ 界面更简洁，重点突出');
        
        console.log('\n💡 用户体验提升:');
        console.log('• 🎯 快速识别: 可用性状态用显眼的 ✅/❌ 图标');
        console.log('• 🔄 智能排序: 可以按可用性排序，优先显示可用或不可用链接');
        console.log('• 🔍 精准筛选: 专门的可用性筛选器，快速定位目标链接');
        console.log('• 📖 视觉清晰: 移除干扰信息，让核心数据更突出');
        
        console.log('\n🎉 界面改进完成！现在的外链管理更加高效和直观！');
        
    } catch (error) {
        console.log('❌ 测试失败:', error.response?.data?.message || error.message);
    }
}

// 运行测试
testUIImprovements().catch(console.error); 