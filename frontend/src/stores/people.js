import { defineStore } from 'pinia';
import { ref } from 'vue';

export const usePeopleStore = defineStore('people', () => {
    const Peoples = ref([]);
    
    async function fetchPeoples(page = 1, pageSize = 20) {
        try {
            const response = await fetch(`/api/v1/frontend/peoples?page=${page}&page_size=${pageSize}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            Peoples.value = data.results;
            return data;
        } catch (error) {
            console.error('获取人物列表失败:', error);
            throw error;
        }
    }

    return { Peoples, fetchPeoples };
})