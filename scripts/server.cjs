const express = require('express');
const bodyParser = require('body-parser');
const { spawn } = require('child_process');
const path = require('path');
const cors = require('cors');
const selfCheckRouter = require('./self-check-api.cjs');

const app = express();
const PORT = 3001;
app.use(cors());
app.use(bodyParser.json());

// 添加自检API路由
app.use(selfCheckRouter);

app.post('/api/create-plugin', (req, res) => {
  const optionsStr = JSON.stringify(req.body);
  const scriptPath = path.join(__dirname, 'create-plugin.cjs');
  const child = spawn('node', [scriptPath, optionsStr]);
  let stdout = '';
  let stderr = '';
  child.stdout.on('data', (data) => {
    stdout += data.toString();
  });
  child.stderr.on('data', (data) => {
    stderr += data.toString();
  });
  child.on('exit', code => {
    if (code === 0) {
      res.json({ success: true, output: stdout.trim() });
    } else {
      res.status(500).json({ success: false, error: stderr || '插件生成失败，可能已存在' });
    }
  });
});

app.post('/api/pack-plugin', (req, res) => {
  const { key } = req.body;
  if (!key) {
    res.status(400).json({ success: false, error: '缺少插件标识符' });
    return;
  }

  const pluginDir = path.join(__dirname, '../src/plugins', `plugin-${key}`);
  if (!fs.existsSync(pluginDir)) {
    res.status(404).json({ success: false, error: '插件不存在' });
    return;
  }

  // 创建临时打包目录
  const tempDir = path.join(__dirname, '../temp', `plugin-${key}-${Date.now()}`);
  fs.mkdirSync(tempDir, { recursive: true });

  try {
    // 复制插件文件到临时目录
    fs.cpSync(pluginDir, path.join(tempDir, `plugin-${key}`), { recursive: true });

    // 创建 package.json
    const packageJson = {
      name: `plugin-${key}`,
      version: '1.0.0',
      private: true,
      type: 'module',
      dependencies: {
        'vue': '^3.4.0',
        'element-plus': '^2.5.0'
      }
    };
    fs.writeFileSync(path.join(tempDir, 'package.json'), JSON.stringify(packageJson, null, 2));

    // 创建安装脚本
    const installScript = `#!/bin/bash
echo "正在安装插件 ${key}..."
npm install
echo "插件安装完成！"
`;
    fs.writeFileSync(path.join(tempDir, 'install.sh'), installScript);
    fs.chmodSync(path.join(tempDir, 'install.sh'), '755');

    // 打包
    const archiver = require('archiver');
    const output = fs.createWriteStream(path.join(__dirname, '../dist', `plugin-${key}.zip`));
    const archive = archiver('zip', { zlib: { level: 9 } });

    output.on('close', () => {
      // 清理临时目录
      fs.rmSync(tempDir, { recursive: true, force: true });
      res.json({ 
        success: true, 
        message: '插件打包成功',
        file: `plugin-${key}.zip`
      });
    });

    archive.on('error', (err) => {
      throw err;
    });

    archive.pipe(output);
    archive.directory(tempDir, false);
    archive.finalize();

  } catch (error) {
    // 清理临时目录
    fs.rmSync(tempDir, { recursive: true, force: true });
    res.status(500).json({ 
      success: false, 
      error: `打包失败: ${error.message}` 
    });
  }
});

const fs = require('fs');
const archiver = require('archiver');

app.get('/pack-plugin', (req, res) => {
  const { key } = req.query;
  if (!key) return res.status(400).send('缺少插件 key');
  const pluginDir = path.join(__dirname, `../plugin-${key}`);
  if (!fs.existsSync(pluginDir)) return res.status(404).send('插件不存在');

  res.setHeader('Content-Type', 'application/zip');
  res.setHeader('Content-Disposition', `attachment; filename=plugin-${key}.zip`);

  const archive = archiver('zip');
  archive.directory(pluginDir, false);
  archive.finalize();
  archive.pipe(res);
});

app.listen(PORT, () => {
  console.log(`插件生成服务已启动: http://localhost:${PORT}`);
});