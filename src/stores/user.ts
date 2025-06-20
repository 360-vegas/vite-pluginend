import { defineStore } from 'pinia'
import { ref } from 'vue'

interface UserInfo {
  id: number
  username: string
  avatar: string
  email: string
  role: string
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo>({
    id: 0,
    username: 'admin',
    avatar: 'https://avatars.githubusercontent.com/u/1?v=4',
    email: 'admin@example.com',
    role: 'admin'
  })

  const token = ref('')

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
  }

  const setToken = (newToken: string) => {
    token.value = newToken
  }

  const logout = async () => {
    // 这里可以添加登出API调用
    userInfo.value = {
      id: 0,
      username: '',
      avatar: '',
      email: '',
      role: ''
    }
    token.value = ''
  }

  return {
    userInfo,
    token,
    setUserInfo,
    setToken,
    logout
  }
}) 