import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import useUserStore from './store/userStore'

import Login from './pages/Login'
import MainLayout from './layouts/MainLayout'
import Dashboard from './pages/Dashboard'
import Profile from './pages/Profile'
import ChangePassword from './pages/ChangePassword'

const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useUserStore()
  return isAuthenticated ? children : <Navigate to="/login" replace />
}

const App = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="profile" element={<Profile />} />
            <Route path="change-password" element={<ChangePassword />} />
            
            <Route path="jobs" element={<div>岗位信息页面（开发中）</div>} />
            <Route path="my-applications" element={<div>我的应聘页面（开发中）</div>} />
            <Route path="favorites" element={<div>我的收藏页面（开发中）</div>} />
            <Route path="policies" element={<div>就业政策页面（开发中）</div>} />
            <Route path="announcements" element={<div>系统公告页面（开发中）</div>} />
            <Route path="my-employment" element={<div>就业信息页面（开发中）</div>} />
            
            <Route path="users" element={<div>用户管理页面（开发中）</div>} />
            <Route path="categories" element={<div>岗位分类管理页面（开发中）</div>} />
            <Route path="applications" element={<div>学生应聘管理页面（开发中）</div>} />
            <Route path="employments" element={<div>就业信息管理页面（开发中）</div>} />
            <Route path="comments" element={<div>岗位评论管理页面（开发中）</div>} />
            <Route path="carousels" element={<div>轮播图管理页面（开发中）</div>} />
            <Route path="statistics" element={<div>就业统计页面（开发中）</div>} />
            <Route path="students" element={<div>学生管理页面（开发中）</div>} />
            
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
