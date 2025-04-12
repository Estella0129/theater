import { defineStore } from 'pinia';
import { ref } from 'vue';

export const usePeopleStore = defineStore('people', () => {
    const Peoples = ref([]);
    
    async function fetchPeoples(page = 1, pageSize = 20, searchQuery = '') {
        try {
            const url = `/api/v1/admin/people?page=${page}&page_size=${pageSize}` + (searchQuery ? `&search=${encodeURIComponent(searchQuery)}` : '');
            const response = await fetch(url);
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

    async function createPeople(peopleData) {
        try {
            const response = await fetch('/api/v1/admin/people', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(peopleData)
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('创建人员失败:', error);
            throw error;
        }
    }

    async function updatePeople(id, peopleData) {
        try {
            const response = await fetch(`/api/v1/admin/people/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(peopleData)
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            return data;
        } catch (error) {
            console.error('更新人员失败:', error);
            throw error;
        }
    }

    return { Peoples, fetchPeoples, createPeople, updatePeople };
})