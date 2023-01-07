<template>
    <TopBar />
    <SearchBar @searchText="searchText" />
    <div class="grid mx-20 mt-4 2xl:grid-cols-2">
        <EmailTable :emails="emails" @selectEmail="selectEmail" />
        <EmailContent :selectedEmail="selectedEmail" :textInput="textInput" />
    </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
import TopBar from "../components/TopBar.vue";
import SearchBar from "../components/SearchBar.vue";
import EmailTable from "../components/EmailTable.vue";
import EmailContent from "../components/EmailContent.vue";
import axios from "axios";
import type { Email } from "../types/email";
export default defineComponent({
    data() {
        return {
            textInput: "" as string,
            emails: [] as Email[],
            selectedEmail: {
                messageID: "",
                content: "",
                date: "",
                from: "",
                subject: "",
                to: [],
            } as Email,
        };
    },
    async mounted() {
        const resp = await axios.get(
            "http://localhost:8080/search?text=" + this.textInput
        );
        if (resp.status === 200) {
            this.emails = resp.data;
            return;
        }
        this.emails = [];
    },
    methods: {
        selectEmail(email: Email) {
            this.selectedEmail = email;
        },
        async searchText(searchTextInput: string) {
            this.textInput = searchTextInput;
            if (searchTextInput !== "") {
                const resp = await axios.get(
                    "http://localhost:8080/search?text=" + searchTextInput
                );
                if (resp.status === 200) {
                    this.emails = resp.data;
                    return;
                }
            }
            this.selectedEmail = {
                messageID: "",
                content: "",
                date: "",
                from: "",
                subject: "",
                to: [],
            };
            this.emails = [];
        },
    },
    components: {
        TopBar,
        SearchBar,
        EmailTable,
        EmailContent,
    },
});
</script>
