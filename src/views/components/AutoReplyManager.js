export default {
    delimiters: ['[[', ']]'],
    data() {
        return {
            replies: [],
            newReply: {
                keyword: '',
                response: '',
                match_type: 'Exact'
            },
            isLoading: false,
            isSubmitting: false
        }
    },
    mounted() {
        this.fetchReplies();
    },
    methods: {
        async fetchReplies() {
            this.isLoading = true;
            try {
                const response = await window.http.get('/api/v1/autoreply');
                this.replies = response.data.results;
            } catch (error) {
                console.error(error);
                window.showErrorInfo("Failed to fetch auto replies");
            } finally {
                this.isLoading = false;
            }
        },
        async addReply() {
            if (!this.newReply.keyword || !this.newReply.response) {
                window.showErrorInfo("Keyword and Response are required");
                return;
            }
            this.isSubmitting = true;
            try {
                await window.http.post('/api/v1/autoreply', this.newReply);
                window.showSuccessInfo("Auto reply added");
                this.newReply = { keyword: '', response: '', match_type: 'Exact' };
                this.fetchReplies();
            } catch (error) {
                console.error(error);
                window.showErrorInfo("Failed to add auto reply");
            } finally {
                this.isSubmitting = false;
            }
        },
        async deleteReply(id) {
            if (!confirm("Are you sure you want to delete this rule?")) return;
            try {
                await window.http.delete(`/api/v1/autoreply/${id}`);
                window.showSuccessInfo("Auto reply deleted");
                this.fetchReplies();
            } catch (error) {
                console.error(error);
                window.showErrorInfo("Failed to delete auto reply");
            }
        }
    },
    template: `
    <div class="ui card fluid">
        <div class="content">
            <div class="header">Auto Reply Manager</div>
            <div class="meta">Manage your automatic responses</div>
            
            <div class="ui form" style="margin-top: 15px;">
                <div class="three fields">
                    <div class="field">
                        <label>Keyword</label>
                        <input type="text" v-model="newReply.keyword" placeholder="Enter keyword...">
                    </div>
                    <div class="field">
                        <label>Match Type</label>
                        <select v-model="newReply.match_type" class="ui dropdown">
                            <option value="Exact">Exact Match</option>
                            <option value="Contains">Contains</option>
                        </select>
                    </div>
                    <div class="field">
                        <label>Response</label>
                        <input type="text" v-model="newReply.response" placeholder="Enter response message...">
                    </div>
                </div>
                <button class="ui primary button" :class="{loading: isSubmitting}" @click="addReply">
                    <i class="plus icon"></i> Add Rule
                </button>
            </div>

            <div class="ui divider"></div>

            <div v-if="isLoading" class="ui active centered inline loader"></div>
            
            <table v-else class="ui celled table">
                <thead>
                    <tr>
                        <th>Keyword</th>
                        <th>Match Type</th>
                        <th>Response</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="reply in replies" :key="reply.id">
                        <td>[[ reply.keyword ]]</td>
                        <td>
                            <div class="ui label" :class="reply.match_type === 'Exact' ? 'blue' : 'teal'">
                                [[ reply.match_type ]]
                            </div>
                        </td>
                        <td>[[ reply.response ]]</td>
                        <td>
                            <button class="ui icon red button mini" @click="deleteReply(reply.id)">
                                <i class="trash icon"></i>
                            </button>
                        </td>
                    </tr>
                    <tr v-if="replies.length === 0">
                        <td colspan="4" class="center aligned">No auto reply rules found</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    `
}
