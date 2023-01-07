export interface Email {
    messageID: string;
    date: string;
    from: string;
    to: string[];
    content: string;
    subject: string;
}
