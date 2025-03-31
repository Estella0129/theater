import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMovieStore = defineStore('movie', () => {
  const movies = ref([])
  
  const fetchMovies = async (page = 1, pageSize = 20) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies?page=${page}&page_size=${pageSize}`)
      const data = await response.json()
      if (page === 1) {
        movies.value = data.results;
      } else {
        movies.value = [...movies.value, ...data.results];
      }
      return data;
    } catch (error) {
      console.error('Failed to fetch movies:', error)
      throw error;
    }
  }
  
  const searchMovies = async (query, page = 1, pageSize = 20) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies/search?query=${encodeURIComponent(query)}&page=${page}&page_size=${pageSize}`)
      const data = await response.json()
      movies.value = data.results;
      return data;
    } catch (error) {
      console.error('Failed to search movies:', error)
      throw error;
    }
  }
  
  return { movies, fetchMovies, searchMovies }
})