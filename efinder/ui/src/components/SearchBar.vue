<template>
    <form
        @submit.prevent="searchText"
        class="flex justify-center mt-4 col-span-full"
    >
        <input
            class="w-1/2 py-2 text-4xl text-center border rounded"
            v-model="searchTextInput"
            type="text"
            placeholder=":)"
        />
    </form>
</template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
    data() {
        return {
            searchTextInput: "" as string,
        };
    },
    methods: {
        searchText() {
            this.$emit("searchText", this.searchTextInput);
        },
    },
});
</script>
