const { execSync } = require('child_process')
const fs = require('fs')
const path = require('path')

// 获取当前版本
function getCurrentVersion() {
  try {
    const version = execSync('git describe --tags --abbrev=0').toString().trim()
    return version.replace('v', '')
  } catch (error) {
    return '1.0.0'
  }
}

// 获取构建时间
function getBuildTime() {
  return new Date().toISOString()
}

// 获取变更日志
function getChangelog() {
  try {
    const lastTag = execSync('git describe --tags --abbrev=0').toString().trim()
    const log = execSync(`git log ${lastTag}..HEAD --pretty=format:"- %s"`).toString()
    return log || '无更新内容'
  } catch (error) {
    return '无更新内容'
  }
}

// 生成版本信息
function generateVersionInfo() {
  const version = getCurrentVersion()
  const buildTime = getBuildTime()
  const changelog = getChangelog()

  const versionInfo = {
    version,
    buildTime,
    changelog,
    files: {
      js: [
        `/static/js/app.${version}.js`,
        `/static/js/vendor.${version}.js`
      ],
      css: [
        `/static/css/app.${version}.css`
      ]
    }
  }

  // 写入版本信息文件
  fs.writeFileSync(
    path.join(__dirname, '../public/version.json'),
    JSON.stringify(versionInfo, null, 2)
  )

  console.log('版本信息已生成:', versionInfo)
}

// 更新版本号
function updateVersion(type = 'patch') {
  const currentVersion = getCurrentVersion()
  const [major, minor, patch] = currentVersion.split('.').map(Number)

  let newVersion
  switch (type) {
    case 'major':
      newVersion = `${major + 1}.0.0`
      break
    case 'minor':
      newVersion = `${major}.${minor + 1}.0`
      break
    case 'patch':
      newVersion = `${major}.${minor}.${patch + 1}`
      break
    default:
      throw new Error('无效的版本更新类型')
  }

  // 创建新的 Git 标签
  execSync(`git tag v${newVersion}`)
  execSync('git push origin --tags')

  console.log(`版本已更新: ${currentVersion} -> ${newVersion}`)
  return newVersion
}

// 构建项目
function buildProject() {
  try {
    execSync('npm run build', { stdio: 'inherit' })
    console.log('项目构建成功')
  } catch (error) {
    console.error('项目构建失败:', error)
    process.exit(1)
  }
}

// 部署项目
function deployProject() {
  try {
    // 这里添加部署命令，例如上传到 CDN
    console.log('项目部署成功')
  } catch (error) {
    console.error('项目部署失败:', error)
    process.exit(1)
  }
}

// 主函数
async function main() {
  const args = process.argv.slice(2)
  const command = args[0]

  switch (command) {
    case 'version':
      const type = args[1] || 'patch'
      updateVersion(type)
      break
    case 'build':
      generateVersionInfo()
      buildProject()
      break
    case 'deploy':
      generateVersionInfo()
      buildProject()
      deployProject()
      break
    default:
      console.log('可用命令:')
      console.log('  version [major|minor|patch] - 更新版本号')
      console.log('  build - 构建项目')
      console.log('  deploy - 部署项目')
  }
}

main() 