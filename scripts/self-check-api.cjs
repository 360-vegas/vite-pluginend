console.log('PATH:', process.env.PATH);
console.log('__dirname:', __dirname);

const express = require('express');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const os = require('os');
const { MongoClient } = require('mongodb');
const router = express.Router();

router.get('/api/self-check', async (req, res) => {
  const logs = [];
  function log(msg) { logs.push(msg); }

  // 1. 插件目录写权限检测
  let pluginsDir = path.join(__dirname, '../src/plugins');
  let writePermission = false;
  try {
    if (!fs.existsSync(pluginsDir)) {
      fs.mkdirSync(pluginsDir, { recursive: true });
      log(`[info] 自动创建插件目录: ${pluginsDir}`);
    }
    const testFile = path.join(pluginsDir, 'self-check-test.txt');
    fs.writeFileSync(testFile, 'test');
    fs.unlinkSync(testFile);
    writePermission = true;
    log('[ok] 插件目录可写');
  } catch (e) {
    log(`[error] 插件目录不可写: ${e.message}`);
    writePermission = false;
  }

  // 2. Node/npm 检测（Windows 优化）
  function checkCommand(cmd) {
    try {
      const result = execSync(cmd, { encoding: 'utf8' });
      // 只要有一行有效路径就算通过
      if (result && result.split(/\r?\n/).some(line => line.trim().length > 0)) {
        log(`[ok] 命令可用: ${result}`);
        return true;
      }
    } catch (e) {
      log(`[error] 命令不可用: ${cmd} (${e.message})`);
      return false;
    }
  }
  const node = checkCommand('where node');
  const npm = checkCommand('where npm');

  // 3. 依赖检测
  let dependencies = false;
  try {
    const nodeModules = path.join(__dirname, '../node_modules');
    if (fs.existsSync(nodeModules)) {
      // 检查关键依赖
      require.resolve('express');
      require.resolve('vue');
      dependencies = true;
      log('[ok] 依赖已安装');
    } else {
      log('[error] node_modules 目录不存在, 请运行 "npm install"');
    }
  } catch (e) {
    log(`[error] 依赖缺失: ${e.message}`);
    dependencies = false;
  }

  // 4. Vite 脚本检测
  let viteScript = false;
  try {
    const pkg = require('../package.json');
    const scripts = pkg.scripts || {};
    viteScript = Object.values(scripts).some(s => s && s.includes('vite'));
    log(viteScript ? '[ok] Vite 脚本存在' : '[error] Vite 脚本缺失 (package.json)');
  } catch (e) {
    log(`[error] 读取 package.json 失败: ${e.message}`);
    viteScript = false;
  }

  // 5. MongoDB 连接测试
  let mongodb = false;
  try {
    const client = new MongoClient('mongodb://localhost:27017', {
      serverSelectionTimeoutMS: 2000, // 2秒超时
      connectTimeoutMS: 2000
    });
    await client.connect();
    await client.db('admin').command({ ping: 1 });
    await client.close();
    mongodb = true;
    log('[ok] MongoDB 连接成功');
  } catch (e) {
    log(`[warn] MongoDB 连接失败: ${e.message}`);
    mongodb = false;
  }

  // 6. MySQL 端口检测（Windows 优化）
  function checkPort(port) {
    try {
      let cmd;
      if (os.platform() === 'win32') {
        cmd = `netstat -ano | findstr :${port}`;
        const result = execSync(cmd, { encoding: 'utf8' });
        if (result.includes('LISTENING')) {
          log(`[ok] 端口 ${port} 已被监听`);
          return true;
        }
      } else {
        cmd = `lsof -i :${port} -sTCP:LISTEN || netstat -anp | grep LISTEN | grep :${port}`;
        execSync(cmd, { stdio: 'ignore' });
        log(`[ok] 端口 ${port} 已被监听`);
        return true;
      }
    } catch (e) {
      log(`[warn] 端口 ${port} 未被监听 (服务可能未运行)`);
      return false;
    }
  }
  const mysql = checkPort(3306);

  // 7. create-plugin.cjs 检测
  let createPluginScript = false;
  try {
    const scriptPath = path.join(__dirname, 'create-plugin.cjs');
    createPluginScript = fs.existsSync(scriptPath);
    log(createPluginScript ? '[ok] create-plugin.cjs 存在' : '[warn] create-plugin.cjs 缺失');
  } catch (e) {
    log(`[error] 检查 create-plugin.cjs 失败: ${e.message}`);
    createPluginScript = false;
  }

  const result = {
    writePermission,
    node,
    npm,
    dependencies,
    viteScript,
    mongodb,
    mysql,
    createPluginScript,
    logs: logs.join('\n')
  };
  res.json(result);
});

module.exports = router;