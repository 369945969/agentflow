// Users API services
import { API_BASE_URL } from '../config/api'

export interface User {
  id: string
  name: string
  email?: string
  avatar?: string
  created_at?: string
  updated_at?: string
}

export interface FetchUsersParams {
  search?: string
  limit?: number
}

// Fetch all users with optional filtering and pagination
export const fetchUsers = async (params: FetchUsersParams = {}): Promise<User[]> => {
  try {
    const queryParams = new URLSearchParams()
    if (params.search) queryParams.append('search', params.search)
    if (params.limit) queryParams.append('limit', String(params.limit))

    const url = params.search
      ? `${API_BASE_URL}/users/?${queryParams}`
      : `${API_BASE_URL}/users/`

    const response = await fetch(url)

    if (!response.ok) {
      throw new Error(`Failed to fetch users: ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    console.error('Error fetching users:', error)
    return []
  }
}

// Fetch users sorted by ID descending (newest first)
export const fetchUsersDesc = async (): Promise<User[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/users/?order=desc&by=id`)

    if (!response.ok) {
      throw new Error(`Failed to fetch users: ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    console.error('Error fetching users:', error)
    return []
  }
}

// Fetch users with search filter
export const searchUsers = async (query: string): Promise<User[]> => {
  try {
    const response = await fetch(`${API_BASE_URL}/users/?search=${encodeURIComponent(query)}`)

    if (!response.ok) {
      throw new Error(`Failed to search users: ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    console.error('Error searching users:', error)
    return []
  }
}