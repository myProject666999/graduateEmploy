import React, { useState, useEffect } from 'react'
import { Card, Row, Col, Statistic, Typography, Tag, Spin, Empty } from 'antd'
import {
  UserOutlined,
  ShopOutlined,
  FileTextOutlined,
  TeamOutlined,
  BarChartOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  HeartOutlined,
} from '@ant-design/icons'
import useUserStore from '../store/userStore'
import request from '../utils/request'

const { Title, Paragraph } = Typography

const Dashboard = () => {
  const { user } = useUserStore()
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState(null)

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

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    setLoading(true)
    try {
      let endpoint = ''
      if (user?.role === 'admin' || user?.role === 'teacher') {
        endpoint = '/teacher/statistics'
      }
      
      if (endpoint) {
        const response = await request.get(endpoint)
        setStats(response.data)
      }
    } catch (error) {
      console.error('获取统计数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const renderAdminStats = () => (
    <Row gutter={[16, 16]}>
      <Col xs={24} sm={12} md={6}>
        <Card>
          <Statistic
            title="学生总数"
            value={stats?.student_count || 0}
            prefix={<TeamOutlined />}
            valueStyle={{ color: '#3f8600' }}
          />
        </Card>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Card>
          <Statistic
            title="岗位总数"
            value={stats?.job_count || 0}
            prefix={<ShopOutlined />}
            valueStyle={{ color: '#1890ff' }}
          />
        </Card>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Card>
          <Statistic
            title="已就业人数"
            value={stats?.employed_count || 0}
            prefix={<CheckCircleOutlined />}
            valueStyle={{ color: '#52c41a' }}
          />
        </Card>
      </Col>
      <Col xs={24} sm={12} md={6}>
        <Card>
          <Statistic
            title="待处理应聘"
            value={stats?.pending_applications || 0}
            prefix={<ClockCircleOutlined />}
            valueStyle={{ color: '#faad14' }}
          />
        </Card>
      </Col>
    </Row>
  )

  const renderStudentWelcome = () => (
    <div>
      <Card>
        <Title level={4}>欢迎使用高校就业管理系统</Title>
        <Paragraph>
          亲爱的 <Tag color="blue">{user?.name || user?.username}</Tag>，欢迎您使用高校就业管理系统！
        </Paragraph>
        <Paragraph>
          作为学生用户，您可以：
        </Paragraph>
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={8}>
            <Card
              size="small"
              hoverable
              icon={<ShopOutlined />}
              title="浏览岗位"
            >
              浏览和搜索招聘岗位信息
            </Card>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Card
              size="small"
              hoverable
              icon={<FileTextOutlined />}
              title="投递简历"
            >
              在线投递应聘申请
            </Card>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Card
              size="small"
              hoverable
              icon={<HeartOutlined />}
              title="收藏岗位"
            >
              收藏感兴趣的岗位
            </Card>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Card
              size="small"
              hoverable
              icon={<FileTextOutlined />}
              title="就业政策"
            >
              查看就业政策信息
            </Card>
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Card
              size="small"
              hoverable
              icon={<UserOutlined />}
              title="就业信息"
            >
              提交和管理就业信息
            </Card>
          </Col>
        </Row>
      </Card>
    </div>
  )

  return (
    <div>
      <Title level={3}>
        欢迎回来，{user?.name || user?.username} <Tag color="blue">[{getRoleName(user?.role)}]</Tag>
      </Title>

      <Spin spinning={loading}>
        {(user?.role === 'admin' || user?.role === 'teacher') && stats ? (
          renderAdminStats()
        ) : (
          renderStudentWelcome()
        )}

        <Card style={{ marginTop: 16 }}>
          <Title level={5}>使用提示</Title>
          <Paragraph>
            1. 点击左侧菜单可以访问各个功能模块
          </Paragraph>
          <Paragraph>
            2. 点击右上角头像可以进入个人中心和修改密码
          </Paragraph>
          <Paragraph>
            3. 如果您在使用过程中有任何问题，请联系系统管理员
          </Paragraph>
        </Card>
      </Spin>
    </div>
  )
}

export default Dashboard
