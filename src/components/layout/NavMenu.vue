<template>
  <div class="nav-menu">
    <el-menu
      :default-active="activeMenu"
      class="el-menu-vertical"
      :collapse="isCollapse"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
    >
      <MenuTree :items="allMenuItems" />
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMenuStore } from '@/stores/menu'
import { Plus, Collection } from '@element-plus/icons-vue'
import { resolveDynamicComponent } from 'vue'

const router = useRouter()
const route = useRoute()
const menuStore = useMenuStore()

const isCollapse = ref(false)
const generatorMenu = [
  {
    path: '/plugin-generator',
    name: 'PluginGenerator',
    meta: {
      title: '插件生成器',
      icon: 'Plus',
      order: 1
    }
  },
  {
    path: '/app-market',
    name: 'AppMarket',
    meta: {
      title: '应用市场',
      icon: 'Collection',
      order: 2
    }
  }
]
const pluginMenu = [
  {
    path: '/plugin-demo',
    name: 'PluginDemo',
    meta: {
      title: '演示插件',
      icon: 'StarFilled',
      order: 1000
    },
    children: [
      {
        path: '/plugin-demo/feature',
        name: 'PluginDemoFeature',
        meta: {
          title: '功能页',
          icon: 'Document',
          order: 1001
        }
      }
    ]
  }
]

const allMenuItems = computed(() => {
  const items = [...generatorMenu, ...menuStore.menuItems, ...pluginMenu]
  console.log('allMenuItems:', items)
  return items
})
const activeMenu = computed(() => route.path)

const handleMenuClick = (item: any) => {
  router.push(item.path)
}

// 递归渲染菜单组件
const MenuTree = {
  name: 'MenuTree',
  props: {
    items: {
      type: Array,
      required: true
    }
  },
  setup(props: any) {
    return () => props.items.map((item: any) => {
      if (item.children && item.children.length > 0) {
        return (
          <el-submenu index={item.path} key={item.path}>
            {{
              title: () => (
                <>
                  {item.meta.icon && <el-icon><component is={resolveDynamicComponent(item.meta.icon)} /></el-icon>}
                  <span>{item.meta.title}</span>
                </>
              ),
              default: () => <MenuTree items={item.children} />
            }}
          </el-submenu>
        )
      } else {
        return (
          <el-menu-item index={item.path} key={item.path} onClick={() => handleMenuClick(item)}>
            {item.meta.icon && <el-icon><component is={resolveDynamicComponent(item.meta.icon)} /></el-icon>}
            <span>{item.meta.title}</span>
          </el-menu-item>
        )
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.nav-menu {
  height: 100%;
  background-color: #304156;

  .el-menu-vertical {
    height: 100%;
    border-right: none;
  }
}
</style> 