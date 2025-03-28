import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMovieStore = defineStore('movie', () => {
  const movies = ref([])
  
  const fetchMovies = async () => {
    try {
      const response = await fetch('/api/v1/frontend/movies')
      const data = await response.json()
      movies.value = data.results;
    } catch (error) {
      console.error('Failed to fetch movies:', error)
    }
  }
  
  return { movies, fetchMovies }
})