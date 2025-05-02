import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useGenreStore = defineStore('genre', () => {
  const genres = ref([])

  const fetchGenres = async () => {
    try {
      const response = await fetch(`/api/v1/admin/genres`)
      const data = await response.json()
      genres.value = data
      console.log(data)
      return data
      
    } catch (error) {
      console.error('Failed to fetch genre:', error)
      throw error
    }
  }

  const getGenreById = async (id) => {
    try {
      const response = await fetch(`/api/v1/admin/genres/${id}`)
      const data = await response.json()
      return data
    } catch (error) {
      console.error('Failed to fetch genre:', error)
      throw error
    }
  }

  const addGenre = async (genre) => {
    try {
      const response = await fetch('/api/v1/admin/genres', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(genre)
      })
      const data = await response.json()
      return data
    } catch (error) {
      console.error('Failed to add genre:', error)
      throw error
    }
  }

  const updateGenre = async (genre) => {
    try {
      // 确保genre对象是有效的JSON数据
      const jsonData = JSON.parse(JSON.stringify(genre));
      const response = await fetch(`/api/v1/admin/genres/${genre.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(jsonData)
      })
      const data = await response.json()
      return data
    } catch (error) {
      console.error('Failed to update genre:', error)
      throw error
    }
  }

  const deleteGenre = async (id) => {
    try {
      const response = await fetch(`/api/v1/admin/genres/${id}`, {
        method: 'DELETE'
      })
      return response.ok
    } catch (error) {
      console.error('Failed to delete genre:', error)
      throw error
    }
  }

  return {
    genres,
    fetchGenres,
    getGenreById,
    addGenre,
    updateGenre,
    deleteGenre
  }
})