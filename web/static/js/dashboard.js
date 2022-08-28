Vue.createApp({
    delimiters: ['{%', '%}'],
    data: () => ({
        user: "",
        links: [],
        showResult: false,
        isURLValid: false,
    }),
    methods: {
        getLinks() {
            const requestOptions = {
                method: 'GET'
            };
            let id = document.querySelector('meta[name="userid"]').content;
            fetch('/api/users/'+id+'/links/', requestOptions)
                .then(async response => {
                    const data = await response.json();
                    // check for error response
                    if (!response.ok) {
                        // get error message from body or default to response status
                        const error = (data && data.message) || response.status;
                        return Promise.reject(error);
                    }
                    this.links = data.found;
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
            this.showResult = true;
        }
    },
    beforeMount(){
        this.getLinks()
    },
}).mount('.linklist');