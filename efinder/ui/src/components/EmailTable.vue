<template>
    <table class="border border-collapse table-auto border-slate-300">
        <thead class="h-auto">
            <tr class="text-xl border bg-slate-500">
                <th>Subject</th>
                <th>From</th>
                <th>To</th>
            </tr>
        </thead>
        <tbody v-for="email in emails" :key="email.messageID">
            <tr
                class="text-lg border cursor-pointer hover:bg-slate-600 active:bg-yellow-300"
                @click="selectEmail(email)"
            >
                <td v-if="email.subject.length < 40">
                    {{ email.subject }}
                </td>
                <td v-else>{{ email.subject.slice(0, 37) + "..." }}</td>
                <td>{{ email.from }}</td>
                <td>{{ email.to[0] }}</td>
            </tr>
        </tbody>
    </table>
</template>
<script lang="ts">
import type { Email } from "@/types/email";
import { defineComponent, type PropType } from "vue";

export default defineComponent({
    methods: {
        selectEmail(email: Email) {
            this.$emit("selectEmail", email);
        },
    },
    props: {
        emails: Array as PropType<Array<Email>>,
    },
});
</script>
