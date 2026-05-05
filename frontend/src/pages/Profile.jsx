import React, { useState, useEffect } from 'react'
import { Form, Input, Button, Card, Typography, message, Spin, Row, Col, Avatar } from 'antd'
import { UserOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons'
import useUserStore from '../store/userStore'
import { updateProfile, getCurrentUser } from '../services/auth'

const { Title } = Typography

const Profile = () => {
  const { user, updateUser } = useUserStore()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    if (user) {
      form.setFieldsValue({
        username: user.username,
        name: user.name,
        email: user.email,
        phone: user.phone,
        department: user.department,
        major: user.major,
        className: user.class_name,
        studentNo: user.student_no,
        teacherNo: user.teacher_no,
      })
    }
  }, [user, form])

  const handleSubmit = async (values) => {
    setLoading(true)
    try {
      const response = await updateProfile(values)
      updateUser(response.data)
      message.success('更新个人信息成功')
    } catch (error) {
      console.error('更新个人信息失败:', error)
    } finally {
      setLoading(false)
    }
  }

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

  return (
    <div>
      <Title level={3}>个人中心</Title>

      <Row gutter={24}>
        <Col xs={24} md={8}>
          <Card style={{ textAlign: 'center' }}>
            <Avatar size={100} icon={<UserOutlined />} style={{ marginBottom: 16 }} />
            <Title level={4}>{user?.name || user?.username}</Title>
            <p style={{ color: '#666' }}>
              用户名：{user?.username}
            </p>
            <p style={{ color: '#666' }}>
              角色：{getRoleName(user?.role)}
            </p>
          </Card>
        </Col>

        <Col xs={24} md={16}>
          <Card title="基本信息">
            <Form
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              initialValues={{
                username: user?.username,
                name: user?.name,
                email: user?.email,
                phone: user?.phone,
                department: user?.department,
                major: user?.major,
                className: user?.class_name,
                studentNo: user?.student_no,
                teacherNo: user?.teacher_no,
              }}
            >
              <Row gutter={16}>
                <Col xs={24} sm={12}>
                  <Form.Item
                    name="username"
                    label="用户名"
                  >
                    <Input disabled prefix={<UserOutlined />} />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item
                    name="name"
                    label="姓名"
                  >
                    <Input prefix={<UserOutlined />} placeholder="请输入姓名" />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16}>
                <Col xs={24} sm={12}>
                  <Form.Item
                    name="email"
                    label="邮箱"
                  >
                    <Input prefix={<MailOutlined />} placeholder="请输入邮箱" />
                  </Form.Item>
                </Col>
                <Col xs={24} sm={12}>
                  <Form.Item
                    name="phone"
                    label="手机号"
                  >
                    <Input prefix={<PhoneOutlined />} placeholder="请输入手机号" />
                  </Form.Item>
                </Col>
              </Row>

              {user?.role === 'student' && (
                <>
                  <Row gutter={16}>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="department"
                        label="院系"
                      >
                        <Input placeholder="请输入院系" />
                      </Form.Item>
                    </Col>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="major"
                        label="专业"
                      >
                        <Input placeholder="请输入专业" />
                      </Form.Item>
                    </Col>
                  </Row>
                  <Row gutter={16}>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="className"
                        label="班级"
                      >
                        <Input placeholder="请输入班级" />
                      </Form.Item>
                    </Col>
                    <Col xs={24} sm={12}>
                      <Form.Item
                        name="studentNo"
                        label="学号"
                      >
                        <Input placeholder="请输入学号" />
                      </Form.Item>
                    </Col>
                  </Row>
                </>
              )}

              {user?.role === 'teacher' && (
                <Row gutter={16}>
                  <Col xs={24} sm={12}>
                    <Form.Item
                      name="department"
                      label="院系"
                    >
                      <Input placeholder="请输入院系" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} sm={12}>
                    <Form.Item
                      name="teacherNo"
                      label="工号"
                    >
                      <Input placeholder="请输入工号" />
                    </Form.Item>
                  </Col>
                </Row>
              )}

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存修改
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Profile
