import apiClient from "./apiClient";

const STORAGE_KEY = "offline_assignments_queue";

class OfflineQueue {
    constructor() {
        this.queue = this.loadQueue();
        if (typeof window !== "undefined") {
            window.addEventListener("online", () => this.processQueue());
        }
    }

    loadQueue() {
        const saved = localStorage.getItem(STORAGE_KEY);
        return saved ? JSON.parse(saved) : [];
    }

    saveQueue() {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(this.queue));
    }

    async send(url, payload) {
        if (!navigator.onLine) {
            this.addToQueue(url, payload);
            return { status: "offline" };
        }

        try {
            const response = await apiClient.post(url, payload);
            return response;
        } catch (error) {
            if (!error.response) {
                console.warn("Network error, queuing request...");
                this.addToQueue(url, payload);
                return { status: "offline" };
            }

            throw error;
        }
    }

    addToQueue(url, payload) {
        const item = {
            id: crypto.randomUUID(),
            url,
            payload,
            timestamp: Date.now(),
        };
        this.queue.push(item);
        this.saveQueue();
    }

    async processQueue() {
        if (this.queue.length === 0) return;
        if (!navigator.onLine) return;

        console.log(`Processing offline queue: ${this.queue.length} items`);

        const queueClone = [...this.queue];
        const failedItems = [];

        for (const item of queueClone) {
            try {
                await apiClient.post(item.url, item.payload);
                console.log(`Item ${item.id} sent successfully`);
            } catch (error) {
                console.error(`Failed to send item ${item.id}`, error);

                if (!error.response) {
                    failedItems.push(item);
                }
            }
        }

        this.queue = failedItems;
        this.saveQueue();
    }
}

const offlineQueue = new OfflineQueue();

export default offlineQueue;
