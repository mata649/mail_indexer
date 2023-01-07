<template>
    <div class="mx-10 text-lg" v-if="selectedEmail !== undefined">
        <div v-if="selectedEmail.messageID.length > 0" class="text-center">
            <span
                v-html="highlight(selectedEmail.messageID)"
                class="font-bold"
            ></span>
        </div>
        <div v-if="selectedEmail.from.length > 0">
            <span class="font-bold">From: </span>
            <span v-html="highlight(selectedEmail.from)"></span>
        </div>
        <div v-if="selectedEmail.to.length > 0">
            <span class="font-bold">To: </span>
            <span
                v-html="
                    highlight(
                        selectedEmail.to.length > 6 && !hideTo
                            ? selectedEmail.to.slice(0, 6).join(', ')
                            : selectedEmail.to.join(', ')
                    )
                "
            ></span>
            <button
                v-if="selectedEmail.to.length > 6"
                class="font-bold"
                @click="hideTo = !hideTo"
            >
                {{ hideTo ? "less" : "..." }}
            </button>
        </div>
        <div v-if="selectedEmail.date.length > 0">
            <span class="font-bold">Date: </span>
            <span v-html="highlight(selectedEmail.date)"></span>
        </div>
        <div v-if="selectedEmail.subject.length > 0">
            <span class="font-bold">Subject: </span>
            <span v-html="highlight(selectedEmail.subject)"></span>
        </div>
        <div
            class="h-screen mt-2 overflow-auto"
            v-if="selectedEmail.content.length > 0"
        >
            <p v-html="highlight(selectedEmail.content)"></p>
        </div>
    </div>
</template>
<script lang="ts">
import type { Email } from "@/types/email";
import { defineComponent, type PropType } from "vue";

export default defineComponent({
    data() {
        return { hideTo: false as boolean };
    },
    methods: {
        highlight(content: string) {
            if (!this.textInput) {
                return content;
            }
            return content.replace(
                new RegExp(this.textInput, "gi"),
                (match) => {
                    return '<span class="bg-yellow-300">' + match + "</span>";
                }
            );
        },
    },
    props: {
        selectedEmail: Object as PropType<Email>,
        textInput: String,
    },
});
</script>
