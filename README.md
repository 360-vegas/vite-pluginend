# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).

# 插件开发与自动化生成

## 插件生成器

1. 进入"插件生成器"页面（如 /plugin-generator）。
2. 输入插件名、主导航、次导航、描述，点击"生成插件"。
3. 系统会自动在根目录下生成 plugin-xxx 文件夹，包含 navigation.json、README.md、pages/index.vue 等标准文件。
4. 重启或热更新后，主导航会自动集成新插件。

## 自动化流程

- 插件目录结构、导航、页面、文档均自动生成，极大提升开发效率。
- 支持多插件共存，主系统自动扫描所有 plugin-*/navigation.json。
- 详细规范见 .cursor/rules/ 相关文档。

## 插件自动生成器（新版）

- 在"插件生成器"页面一次性填写插件所有元信息（中文名、英文简称、作者、描述、图标、分类、标签、页面等）。
- 支持动态添加页面，页面路径/文件名可自定义。
- 生成后自动创建 plugin-xxx 目录，生成 meta.ts、pages、assets 等。
- meta.ts 作为插件唯一元信息入口，所有导航、页面、路由均自动集成。
- 支持批量扫描所有 plugin-*/meta.ts，实现主导航、次导航、路由的自动注册。
