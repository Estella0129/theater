import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMovieStore = defineStore('movie', () => {
  const movies = ref([])

  const fetchMovies = async (query = { page: 1, pageSize: 20, genre: '' }) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies?page=${query.page}&page_size=${query.pageSize}&genre=${query.genre}`)
      const data = await response.json()
      movies.value = data.results;
      return data;
    } catch (error) {
      console.error('Failed to fetch movies:', error)
      throw error;
    }
  }
  
  const searchResults = ref([])
  const searchTotalPages = ref(1)
  
  const searchMovies = async (query, page = 1, pageSize = 20) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies?query=${encodeURIComponent(query)}&page=${page}&page_size=${pageSize}`)
      const data = await response.json()
      searchResults.value = data.results;
      searchTotalPages.value = data.total_pages;
      return data;
    } catch (error) {
      console.error('Failed to search movies:', error)
      throw error;
    }
  }
  
  const getMovieById = async (id) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies/${id}`)
      const data = await response.json()

      const director = data.Credits.find(c => c.credit_type == "crew" && c.order == 0)
      const cast = data.Credits.find(c => c.credit_type == "cast" && c.order == 0)

      data.director = director && director.People ? director.People.name : ""
      data.cast = cast && cast.People ? cast.People.name : ""

      return data
    } catch (error) {
      console.error('Failed to fetch movie:', error)
      throw error
    }
  }
const genres = ref([])

const fetchGenres = async () => {
  try {
    const response = await fetch(`/api/v1/frontend/genres`)
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
    const response = await fetch(`/api/v1/frontend/genres/${id}`)
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Failed to fetch genre:', error)
    throw error
  }
}

  return { movies, fetchMovies, searchMovies, searchResults, searchTotalPages, getMovieById, genres, fetchGenres,getGenreById  }
})