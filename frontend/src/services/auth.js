import request from '../utils/request'

export const login = (data) => {
  return request.post('/login', data)
}

export const register = (data) => {
  return request.post('/register', data)
}

export const getCurrentUser = () => {
  return request.get('/user/me')
}

export const updateProfile = (data) => {
  return request.put('/user/profile', data)
}

export const changePassword = (data) => {
  return request.post('/user/change-password', data)
}
