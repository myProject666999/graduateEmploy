import React, { useState } from 'react'
import { Layout, Menu, Dropdown, Avatar, theme, Typography, message } from 'antd'
import {
  HomeOutlined,
  UserOutlined,
  LogoutOutlined,
  LockOutlined,
  FileTextOutlined,
  NotificationOutlined,
  HeartOutlined,
  TeamOutlined,
  BarChartOutlined,
  PictureOutlined,
  TagsOutlined,
  AppstoreOutlined,
  ShopOutlined,
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import useUserStore from '../store/userStore'

const { Header, Sider, Content } = Layout
const { Title } = Typography

const roleMenus = {
  admin: [
    {
      key: '/dashboard',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/users',
      icon: <TeamOutlined />,
      label: '用户管理',
    },
    {
      key: '/categories',
      icon: <TagsOutlined />,
      label: '岗位分类管理',
    },
    {
      key: '/jobs',
      icon: <ShopOutlined />,
      label: '招聘信息管理',
    },
    {
      key: '/applications',
      icon: <FileTextOutlined />,
      label: '学生应聘管理',
    },
    {
      key: '/employments',
      icon: <UserOutlined />,
      label: '就业信息管理',
    },
    {
      key: '/comments',
      icon: <NotificationOutlined />,
      label: '岗位评论管理',
    },
    {
      key: '/policies',
      icon: <FileTextOutlined />,
      label: '就业政策管理',
    },
    {
      key: '/announcements',
      icon: <NotificationOutlined />,
      label: '系统公告管理',
    },
    {
      key: '/carousels',
      icon: <PictureOutlined />,
      label: '轮播图管理',
    },
    {
      key: '/statistics',
      icon: <BarChartOutlined />,
      label: '就业统计管理',
    },
  ],
  teacher: [
    {
      key: '/dashboard',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/students',
      icon: <TeamOutlined />,
      label: '学生管理',
    },
    {
      key: '/employments',
      icon: <UserOutlined />,
      label: '就业信息管理',
    },
    {
      key: '/statistics',
      icon: <BarChartOutlined />,
      label: '就业统计',
    },
  ],
  student: [
    {
      key: '/dashboard',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/jobs',
      icon: <ShopOutlined />,
      label: '岗位信息',
    },
    {
      key: '/my-applications',
      icon: <FileTextOutlined />,
      label: '我的应聘',
    },
    {
      key: '/favorites',
      icon: <HeartOutlined />,
      label: '我的收藏',
    },
    {
      key: '/policies',
      icon: <FileTextOutlined />,
      label: '就业政策',
    },
    {
      key: '/announcements',
      icon: <NotificationOutlined />,
      label: '系统公告',
    },
    {
      key: '/my-employment',
      icon: <UserOutlined />,
      label: '就业信息',
    },
  ],
}

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useUserStore()
  const [collapsed, setCollapsed] = useState(false)
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken()

  const role = user?.role || 'student'
  const menuItems = roleMenus[role] || roleMenus.student

  const getRoleName = (role) => {
    switch (role) {
      case 'admin':
        return '管理员'
      case 'teacher':
        return '教师'
      case 'student':
        return '学生'
      default:
        return '用户'
    }
  }

  const handleMenuClick = ({ key }) => {
    navigate(key)
  }

  const handleLogout = () => {
    logout()
    message.success('退出登录成功')
    navigate('/login')
  }

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'password',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: () => navigate('/change-password'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={(value) => setCollapsed(value)}
        theme="dark"
        width={240}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'rgba(255, 255, 255, 0.1)',
          }}
        >
          <Title
            level={4}
            style={{ color: '#fff', margin: 0, whiteSpace: 'nowrap' }}
          >
            {collapsed ? '就业' : '高校就业管理系统'}
          </Title>
        </div>

        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>

      <Layout>
        <Header
          style={{
            padding: '0 24px',
            background: colorBgContainer,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <div />
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
                <Avatar icon={<UserOutlined />} />
                <span>{user?.name || user?.username}</span>
                <span style={{ color: '#999', fontSize: 12 }}>
                  [{getRoleName(user?.role)}]
                </span>
              </div>
            </Dropdown>
          </div>
        </Header>

        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
            minHeight: 280,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
