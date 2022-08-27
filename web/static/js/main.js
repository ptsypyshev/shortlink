Vue.createApp({
    delimiters: ['{%', '%}'],
    data: () => ({
        shortlink: {
            longLink: "",
            shortURL: "",
        },
        showResult: false,
        isURLValid: false,
    }),
    methods: {
        change:function(e){
            const url = e.target.value
            this.showResult = false
            this.isURLValid = validator.isURL(this.shortlink.longLink, {require_protocol: true});
        },
        shortenLink() {
            const requestOptions = {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ "long_link": this.shortlink.longLink })
            };
            fetch('/api/links/', requestOptions)
                .then(async response => {
                    const data = await response.json();
                    // check for error response
                    if (!response.ok) {
                        // get error message from body or default to response status
                        const error = (data && data.message) || response.status;
                        return Promise.reject(error);
                    }
                    var newUrl = window.location.protocol + "//" + window.location.host + "/" + data.created;
                    this.shortlink.shortURL = newUrl;
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
            this.showResult = !this.showResult;
        }
    },
}).mount('.shortener');