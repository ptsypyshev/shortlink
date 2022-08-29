Vue.createApp({
    delimiters: ['{%', '%}'],
    data: () => ({
        shortlink: {
            longLink: "",
            shortURL: "",
        },
        showShortLink: false,
        isURLValid: false,
        user: "",
        links: [],
        showLinks: false,
        askInitDB: false,
    }),
    watch: {
        showShortLink() {
            this.getLinksOnDashboard();
        }
    },
    methods: {
        change: function () {
            // const url = e.target.value
            this.showShortLink = false
            this.isURLValid = validator.isURL(this.shortlink.longLink, {require_protocol: true});
        },
        shortenLink() {
            const requestOptions = {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({"long_link": this.shortlink.longLink, "is_active": true})
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
                    this.shortlink.shortURL = window.location.protocol + "//" + window.location.host + "/" + data.created;
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
            this.showShortLink = true;
        },
        getLinks() {
            const requestOptions = {
                method: 'GET'
            };
            let id = document.querySelector('meta[name="userid"]').content;
            fetch('/api/users/' + id + '/links/', requestOptions)
                .then(async response => {
                    const data = await response.json();
                    // check for error response
                    if (!response.ok) {
                        // get error message from body or default to response status
                        const error = (data && data.message) || response.status;
                        return Promise.reject(error);
                    }
                    let links = data.found;
                    let linksLen = links.length;
                    for (let i = 0; i < linksLen; i++) {
                        links[i].short_link = window.location.protocol + "//" + window.location.host + "/" + links[i].short_link;
                    }
                    this.links = links;
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
            this.showLinks = true;
        },
        initDB() {
            this.getRequest('/dbinit/', "Do you really want to Init DB? All data will be lost...");
            this.getLinks();
        },
        addDemoData() {
            this.getRequest('/demodb/', "Do you really want add demo data?");
            this.getLinks();
        },
        getRequest(path, msg) {
            let answer = confirm(msg);
            if (answer) {
                const requestOptions = {
                    method: 'GET'
                };
                fetch(path, requestOptions)
                    .then(async response => {
                        await response.json();
                        // check for error response
                        if (!response.ok) {
                            // get error message from body or default to response status
                            const error = (data && data.message) || response.status;
                            alert("Error during Init DB: " + response.statusText)
                            return Promise.reject(error);
                        }
                        alert("Successful Init DB: " + response.statusText)
                    })
                    .catch(error => {
                        this.errorMessage = error;
                        console.error('There was an error!', error);
                    });
            }
        },
        getLinksOnDashboard() {
            let page_template = document.querySelector('meta[name="page_template"]').content;
            if (page_template == "dashboard") {
                this.getLinks();
            }
        }
    },
    beforeMount(){
        this.getLinksOnDashboard();
    },
}).mount('.vueapp');