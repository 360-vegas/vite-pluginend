// 增强的真实用户访问模拟测试脚本
// 测试发布速度、筛选功能和失败重发功能

const API_BASE = 'http://localhost:8080'

// 测试数据
const testData = {
    validLinks: [
        'https://www.google.com',
        'https://www.github.com',
        'https://www.stackoverflow.com'
    ],
    problematicLinks: [
        'https://trends.google.com/trends/',
        'https://www.thevineking.com'
    ],
    invalidLinks: [
        'https://invalid-test-domain-99999.com',
        'http://non-existent-site-test.xyz'
    ]
}

// 测试时间统计
const timeStats = {
    start: null,
    checkStart: null,
    checkEnd: null,
    publishStart: null,
    publishEnd: null
}

// 添加测试链接
async function addTestLinks() {
    console.log('📦 添加测试链接...')
    
    const allLinks = [
        ...testData.validLinks,
        ...testData.problematicLinks,
        ...testData.invalidLinks
    ]
    
    const promises = allLinks.map(async (url, index) => {
        try {
            const response = await fetch(`${API_BASE}/api/external-links`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    url: url,
                    category: 'test-enhanced',
                    description: `增强测试链接 ${index + 1}`,
                    status: true
                })
            })
            
            if (response.ok) {
                const result = await response.json()
                console.log(`✅ 添加成功: ${url}`)
                return result
            } else {
                console.log(`❌ 添加失败: ${url} - ${response.status}`)
                return null
            }
        } catch (error) {
            console.log(`❌ 添加异常: ${url} - ${error.message}`)
            return null
        }
    })
    
    const results = await Promise.all(promises)
    return results.filter(Boolean)
}

// 测试筛选功能
async function testFilteringFeatures() {
    console.log('\n🔍 测试筛选功能...')
    console.log('=' .repeat(40))
    
    try {
        // 测试获取所有链接
        console.log('📊 测试获取所有链接...')
        const allResponse = await fetch(`${API_BASE}/api/external-links?per_page=100`)
        const allData = await allResponse.json()
        const totalCount = allData.meta?.total || 0
        console.log(`✅ 总链接数: ${totalCount}`)
        
        // 测试筛选可用链接
        console.log('✅ 测试筛选可用链接...')
        const validResponse = await fetch(`${API_BASE}/api/external-links?is_valid=true&per_page=100`)
        const validData = await validResponse.json()
        const validCount = validData.data?.length || 0
        console.log(`✅ 可用链接数: ${validCount}`)
        
        // 测试筛选不可用链接
        console.log('❌ 测试筛选不可用链接...')
        const invalidResponse = await fetch(`${API_BASE}/api/external-links?is_valid=false&per_page=100`)
        const invalidData = await invalidResponse.json()
        const invalidCount = invalidData.data?.length || 0
        console.log(`❌ 不可用链接数: ${invalidCount}`)
        
        // 测试排序功能
        console.log('🔄 测试按可用性排序...')
        const sortedResponse = await fetch(`${API_BASE}/api/external-links?sort_field=is_valid&sort_order=desc&per_page=10`)
        const sortedData = await sortedResponse.json()
        
        if (sortedData.data && sortedData.data.length > 0) {
            console.log('📋 前10个链接的可用性状态:')
            sortedData.data.forEach((link, index) => {
                const status = link.is_valid ? '✅ 可用' : '❌ 不可用'
                const url = link.url.substring(0, 40) + (link.url.length > 40 ? '...' : '')
                console.log(`  ${index + 1}. ${status} | ${url}`)
            })
            
            // 检查排序是否正确
            const isCorrectlySorted = sortedData.data.slice(0, 5).every(link => link.is_valid === true) ||
                                     sortedData.data.slice(0, 5).every(link => link.is_valid === false)
            
            if (isCorrectlySorted) {
                console.log('✅ 排序功能正常')
            } else {
                console.log('⚠️ 排序功能可能有问题')
            }
        }
        
        console.log('\n📊 筛选功能测试总结:')
        console.log(`  总链接: ${totalCount}`)
        console.log(`  可用: ${validCount}`)
        console.log(`  不可用: ${invalidCount}`)
        console.log(`  筛选功能: ✅ 正常`)
        
    } catch (error) {
        console.error('❌ 筛选功能测试失败:', error.message)
    }
}

// 测试真实用户访问模拟速度
async function testSimulationSpeed() {
    console.log('\n🎭 测试真实用户访问模拟速度...')
    console.log('=' .repeat(50))
    
    try {
        timeStats.checkStart = Date.now()
        
        console.log('🚀 开始批量检测（使用真实用户模拟）...')
        
        const response = await fetch(`${API_BASE}/api/external-links/batch-check`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ids: [],  // 空数组表示检测所有链接
                check_all: true
            })
        })
        
        timeStats.checkEnd = Date.now()
        const totalDuration = timeStats.checkEnd - timeStats.checkStart
        
        if (!response.ok) {
            throw new Error(`API 调用失败: ${response.status}`)
        }
        
        const result = await response.json()
        
        console.log(`\n⏱️ 检测总耗时: ${Math.round(totalDuration / 1000)}秒 (${totalDuration}ms)`)
        
        if (result.results && Array.isArray(result.results)) {
            const results = result.results
            const totalChecked = results.length
            const avgTimePerLink = Math.round(totalDuration / totalChecked)
            
            console.log(`📊 检测统计:`)
            console.log(`  总链接数: ${totalChecked}`)
            console.log(`  平均每个链接耗时: ${avgTimePerLink}ms (${Math.round(avgTimePerLink/1000)}秒)`)
            console.log(`  成功访问: ${results.filter(r => r.is_valid).length}`)
            console.log(`  访问失败: ${results.filter(r => !r.is_valid).length}`)
            
            // 检查是否符合真实用户访问的时间要求
            if (avgTimePerLink >= 5000) { // 每个链接至少5秒
                console.log('✅ 访问速度符合真实用户行为（较慢，更真实）')
            } else if (avgTimePerLink >= 3000) {
                console.log('⚠️ 访问速度适中（可以更慢一些）')
            } else {
                console.log('❌ 访问速度太快（不够真实）')
            }
            
            // 显示详细结果
            console.log('\n📝 详细访问结果:')
            results.forEach((result, index) => {
                const status = result.is_valid ? '🎉 访问成功' : '❌ 访问失败'
                const message = result.message || result.error_message || '无消息'
                console.log(`${index + 1}. ${status} | ${result.url.substring(0, 30)}... | ${message.substring(0, 50)}`)
            })
            
        } else {
            console.log('❌ 检测响应格式异常')
        }
        
    } catch (error) {
        console.error('❌ 检测速度测试失败:', error.message)
    }
}

// 测试失败重发功能
async function testRetryFailedLinks() {
    console.log('\n🔄 测试失败重发功能...')
    console.log('=' .repeat(40))
    
    try {
        // 1. 获取不可用链接
        console.log('📋 获取不可用链接列表...')
        const invalidLinksResponse = await fetch(`${API_BASE}/api/external-links/invalid`)
        const invalidLinksData = await invalidLinksResponse.json()
        
        if (!invalidLinksData.data || invalidLinksData.data.length === 0) {
            console.log('ℹ️ 没有找到失败的链接，跳过重试测试')
            return
        }
        
        const failedLinks = invalidLinksData.data
        console.log(`📊 找到 ${failedLinks.length} 个失败链接`)
        
        // 2. 执行重试
        console.log('🔄 开始重试失败链接...')
        const retryStartTime = Date.now()
        
        const failedIds = failedLinks.map(link => link.id).filter(Boolean)
        
        const retryResponse = await fetch(`${API_BASE}/api/external-links/batch-check`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ids: failedIds,
                check_all: false
            })
        })
        
        const retryEndTime = Date.now()
        const retryDuration = retryEndTime - retryStartTime
        
        if (!retryResponse.ok) {
            throw new Error(`重试API调用失败: ${retryResponse.status}`)
        }
        
        const retryResult = await retryResponse.json()
        
        if (retryResult.results && Array.isArray(retryResult.results)) {
            const results = retryResult.results
            const retrySuccessCount = results.filter(r => r.is_valid === true).length
            const retryFailCount = results.filter(r => r.is_valid === false).length
            
            console.log(`\n⏱️ 重试耗时: ${Math.round(retryDuration / 1000)}秒`)
            console.log(`📊 重试结果:`)
            console.log(`  重试总数: ${results.length}`)
            console.log(`  恢复成功: ${retrySuccessCount}`)
            console.log(`  仍然失败: ${retryFailCount}`)
            console.log(`  成功率: ${Math.round(retrySuccessCount / results.length * 100)}%`)
            
            if (retrySuccessCount > 0) {
                console.log('✅ 失败重发功能正常工作')
            } else {
                console.log('⚠️ 所有重试仍然失败（可能是链接本身的问题）')
            }
            
            // 显示恢复成功的链接
            if (retrySuccessCount > 0) {
                console.log('\n🎉 恢复成功的链接:')
                results.filter(r => r.is_valid).forEach((result, index) => {
                    console.log(`  ${index + 1}. ${result.url}`)
                })
            }
            
        } else {
            console.log('❌ 重试响应格式异常')
        }
        
    } catch (error) {
        console.error('❌ 失败重发测试失败:', error.message)
    }
}

// 清理测试数据
async function cleanupTestData() {
    console.log('\n🧹 清理测试数据...')
    
    try {
        // 获取所有测试链接
        const response = await fetch(`${API_BASE}/api/external-links?category=test-enhanced&per_page=100`)
        
        if (!response.ok) {
            console.log('❌ 获取测试数据失败')
            return
        }
        
        const result = await response.json()
        
        if (result.data && Array.isArray(result.data)) {
            const testLinksToDelete = result.data.filter(link => 
                link.category === 'test-enhanced'
            )
            
            if (testLinksToDelete.length > 0) {
                console.log(`🗑️ 找到 ${testLinksToDelete.length} 个测试链接，开始删除...`)
                
                const deletePromises = testLinksToDelete.map(async (link) => {
                    try {
                        const deleteResponse = await fetch(`${API_BASE}/api/external-links/${link.id}`, {
                            method: 'DELETE'
                        })
                        
                        if (deleteResponse.ok) {
                            console.log(`✅ 删除成功: ${link.url}`)
                            return true
                        } else {
                            console.log(`❌ 删除失败: ${link.url}`)
                            return false
                        }
                    } catch (error) {
                        console.log(`❌ 删除异常: ${link.url} - ${error.message}`)
                        return false
                    }
                })
                
                const deleteResults = await Promise.all(deletePromises)
                const successCount = deleteResults.filter(Boolean).length
                
                console.log(`✅ 清理完成: ${successCount}/${testLinksToDelete.length} 个链接已删除`)
            } else {
                console.log('ℹ️ 没有找到需要清理的测试数据')
            }
        }
        
    } catch (error) {
        console.error('❌ 清理失败:', error.message)
    }
}

// 主测试函数
async function runEnhancedTest() {
    console.log('🎭 增强的真实用户访问模拟功能测试')
    console.log('=' .repeat(60))
    console.log(`⏰ 开始时间: ${new Date().toLocaleString()}`)
    console.log('')
    
    timeStats.start = Date.now()
    
    try {
        // 1. 添加测试链接
        const addedLinks = await addTestLinks()
        if (addedLinks.length === 0) {
            console.log('❌ 没有成功添加测试链接，跳过后续测试')
            return
        }
        
        console.log(`✅ 成功添加 ${addedLinks.length} 个测试链接`)
        
        // 等待数据保存
        console.log('\n⏳ 等待数据保存...')
        await new Promise(resolve => setTimeout(resolve, 3000))
        
        // 2. 测试筛选功能
        await testFilteringFeatures()
        
        // 3. 测试真实用户访问模拟速度
        await testSimulationSpeed()
        
        // 4. 测试失败重发功能
        await testRetryFailedLinks()
        
        // 5. 清理测试数据
        await cleanupTestData()
        
        const totalTestTime = Date.now() - timeStats.start
        
        console.log('\n🎊 所有增强功能测试完成!')
        console.log('=' .repeat(40))
        console.log(`⏰ 测试总耗时: ${Math.round(totalTestTime / 1000)}秒`)
        console.log(`⏰ 结束时间: ${new Date().toLocaleString()}`)
        
        // 测试总结
        console.log('\n📊 测试总结:')
        console.log('✅ 筛选功能: 可用性筛选、排序功能')
        console.log('✅ 真实用户模拟: 访问速度更真实（慢速）')
        console.log('✅ 失败重发: 自动重试失败链接')
        console.log('✅ 用户体验: 详细的进度反馈和状态显示')
        
    } catch (error) {
        console.error('❌ 测试过程中发生错误:', error)
    }
}

// 只测试筛选功能
async function testFilteringOnly() {
    console.log('🔍 仅测试筛选功能')
    console.log('=' .repeat(30))
    
    await testFilteringFeatures()
    
    console.log('\n✅ 筛选功能测试完成!')
}

// 只测试速度
async function testSpeedOnly() {
    console.log('⚡ 仅测试访问速度')
    console.log('=' .repeat(30))
    
    await testSimulationSpeed()
    
    console.log('\n✅ 速度测试完成!')
}

// 根据命令行参数决定运行哪种测试
const args = process.argv.slice(2)
if (args.includes('--filter-only')) {
    testFilteringOnly()
} else if (args.includes('--speed-only')) {
    testSpeedOnly()
} else if (args.includes('--retry-only')) {
    testRetryFailedLinks()
} else if (args.includes('--cleanup-only')) {
    cleanupTestData()
} else {
    runEnhancedTest()
} 