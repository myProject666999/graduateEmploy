import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Tabs, Typography } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { login, register } from '../services/auth'
import useUserStore from '../store/userStore'

const { Title } = Typography

const Login = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const { login: setLogin } = useUserStore()

  const handleLogin = async (values) => {
    setLoading(true)
    try {
      const response = await login(values)
      const { token, user } = response.data
      setLogin(token, user)
      message.success('登录成功')
      navigate('/')
    } catch (error) {
      console.error('登录失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (values) => {
    if (values.password !== values.confirmPassword) {
      message.error('两次输入的密码不一致')
      return
    }
    setLoading(true)
    try {
      const { password, confirmPassword, ...registerData } = values
      const response = await register({ ...registerData, password })
      const { token, user } = response.data
      setLogin(token, user)
      message.success('注册成功')
      navigate('/')
    } catch (error) {
      console.error('注册失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const loginForm = (
    <Form
      name="login"
      onFinish={handleLogin}
      autoComplete="off"
      size="large"
    >
      <Form.Item
        name="username"
        rules={[{ required: true, message: '请输入用户名' }]}
      >
        <Input prefix={<UserOutlined />} placeholder="用户名" />
      </Form.Item>

      <Form.Item
        name="password"
        rules={[{ required: true, message: '请输入密码' }]}
      >
        <Input.Password prefix={<LockOutlined />} placeholder="密码" />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading} block>
          登录
        </Button>
      </Form.Item>
    </Form>
  )

  const registerForm = (
    <Form
      name="register"
      onFinish={handleRegister}
      autoComplete="off"
      size="large"
    >
      <Form.Item
        name="username"
        rules={[{ required: true, message: '请输入用户名' }]}
      >
        <Input prefix={<UserOutlined />} placeholder="用户名" />
      </Form.Item>

      <Form.Item
        name="name"
        rules={[{ required: true, message: '请输入姓名' }]}
      >
        <Input prefix={<UserOutlined />} placeholder="姓名" />
      </Form.Item>

      <Form.Item
        name="email"
        rules={[{ type: 'email', message: '请输入有效的邮箱' }]}
      >
        <Input prefix={<MailOutlined />} placeholder="邮箱（选填）" />
      </Form.Item>

      <Form.Item
        name="phone"
      >
        <Input prefix={<PhoneOutlined />} placeholder="手机号（选填）" />
      </Form.Item>

      <Form.Item
        name="password"
        rules={[{ required: true, message: '请输入密码' }]}
      >
        <Input.Password prefix={<LockOutlined />} placeholder="密码" />
      </Form.Item>

      <Form.Item
        name="confirmPassword"
        rules={[{ required: true, message: '请确认密码' }]}
      >
        <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading} block>
          注册
        </Button>
      </Form.Item>
    </Form>
  )

  const items = [
    {
      key: 'login',
      label: '登录',
      children: loginForm,
    },
    {
      key: 'register',
      label: '注册',
      children: registerForm,
    },
  ]

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '20px',
      }}
    >
      <Card
        style={{
          width: 400,
          boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={3} style={{ margin: 0 }}>
            高校就业管理系统
          </Title>
          <p style={{ color: '#666', marginTop: 8 }}>
            欢迎使用高校就业管理系统
          </p>
        </div>

        <Tabs defaultActiveKey="login" items={items} centered />

        <div style={{ marginTop: 16, textAlign: 'center', fontSize: 12, color: '#999' }}>
          <p>测试账户：</p>
          <p>管理员: admin / admin123</p>
          <p>教师: teacher1 / teacher123</p>
          <p>学生: student1 / student123</p>
        </div>
      </Card>
    </div>
  )
}

export default Login
