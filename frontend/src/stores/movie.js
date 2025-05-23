import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useMovieStore = defineStore('movie', () => {
  const movies = ref([])

  const fetchMovies = async (query = { page: 1, pageSize: 20, genre: '', sort_by: '' }) => {
    try {
      let url = `/api/v1/frontend/movies?page=${query.page}&page_size=${query.pageSize}`;
      if (query.genre) url += `&genre=${query.genre}`;
      if (query.sort_by) url += `&sort_by=${query.sort_by}`;
      
      const response = await fetch(url);
      const data = await response.json();
      console.log('原始电影数据:', data.results);
      movies.value = data.results.map(movie => {
        return {
          ...movie,
          releaseDate: movie.release_date,
          rating: movie.vote_average / 2,
          director: movie.Director?.People?.name || "暂无导演信息"
        };
      });
      return data;
    } catch (error) {
      console.error('Failed to fetch movies:', error);
      throw error;
    }
  }
  const fetchAdminMovies = async (query = { page: 1, pageSize: 20, genre: '', sort_by: '' }) => {
    try {
      let url = `/api/v1/admin/movies?page=${query.page}&page_size=${query.pageSize}`;
      if (query.genre) url += `&genre=${query.genre}`;
      if (query.sort_by) url += `&sort_by=${query.sort_by}`;
      
      const response = await fetch(url);
      const data = await response.json();
      console.log('原始电影数据:', data.results);
      movies.value = data.results.map(movie => {
        return {
          ...movie,
          releaseDate: movie.release_date,
          rating: movie.vote_average / 2,
          director: movie.Director?.People?.name || "暂无导演信息"
        };
      });
      return data;
    } catch (error) {
      console.error('Failed to fetch movies:', error);
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
  
  const fetchPeople = async (params = {}) => {
    try {
      const query = new URLSearchParams(params).toString();
      const response = await fetch(`/api/v1/admin/people?${query}`);
      const data = await response.json();
      return data.results.map(p => ({ 
        id: p.id, 
        name: p.name,
        profile_path: p.profile_path
      }));
    } catch (error) {
      console.error('Failed to fetch people:', error);
      throw error;
    }
  };

  const updateMovie = async (movieData) => {
    try {
      const response = await fetch(`/api/v1/admin/movies/${movieData.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(movieData)
      })
      const data = await response.json()
      return data;
    } catch (error) {
      console.error('Failed to update movie:', error)
      throw error;
    }
  }


  const createMovie = async (movieData) => {
    try {
      const response = await fetch('/api/v1/admin/movies', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(movieData)
      })
      const data = await response.json()
      return data
    } catch (error) {
      console.error('Failed to create movie:', error)
      throw error
    }
  }

  const getMovieById = async (id) => {
    try {
      const response = await fetch(`/api/v1/frontend/movies/${id}`)
      const data = await response.json()

      // 确保Credits数组存在
      if (!data.Credits) {
        data.Credits = []
      }

      const director = data.Credits.find(c => c.credit_type == "crew" && c.order == 0)
      const cast = data.Credits.find(c => c.credit_type == "cast" && c.order == 0)

      data.director = director && director.People ? director.People.name : "暂无导演信息"
      data.cast = cast && cast.People ? cast.People.name : "暂无主演信息"

      data.rating = data.vote_average/2;

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

  return { movies, fetchMovies,fetchAdminMovies, updateMovie,createMovie, searchMovies, searchResults, searchTotalPages, getMovieById, genres, fetchGenres, getGenreById, fetchPeople }
})