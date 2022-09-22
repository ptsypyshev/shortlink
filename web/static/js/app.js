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
        users: [],
        userform: {
            isForm: true,
            id: "",
            username: "",
            password: "",
            first_name: "",
            last_name: "",
            email: "",
            phone: "",
            user_status: ""
        },
        links: [],
        showLinks: false,
        showUserEditForm: false,
    }),
    watch: {
        showShortLink() {
            this.getObjectsForTemplate();
        }
    },
    methods: {
        change: function (e) {
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
        getUsers() {
            const requestOptions = {
                method: 'GET'
            };
            let id = document.querySelector('meta[name="userid"]').content;
            fetch('/api/users/', requestOptions)
                .then(async response => {
                    const data = await response.json();
                    // check for error response
                    if (!response.ok) {
                        // get error message from body or default to response status
                        const error = (data && data.message) || response.status;
                        return Promise.reject(error);
                    }
                    this.users = data.found;
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
        },
        createUser(user) {
            let json_string = {};
            ["username", "password", "first_name", "last_name", "email", "phone", "user_status"].forEach(function (elem) {
                if (user[elem] != "" || elem == "user_status") {
                    json_string[elem] = user[elem];
                }
            })
            const requestOptions = {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(json_string)
            };
            fetch('/api/users/', requestOptions)
                .then(async response => {
                    const data = await response.json();
                    // check for error response
                    if (!response.ok) {
                        // get error message from body or default to response status
                        const error = (data && data.message) || response.status;
                        return Promise.reject(error);
                    }
                    this.getObjectsForTemplate();
                })
                .catch(error => {
                    this.errorMessage = error;
                    console.error('There was an error!', error);
                });
        },
        createUserShowForm() {
            for (const userformKey in this.userform) {
                this.userform[userformKey] = null;
            }
            this.userform["user_status"] = true;
            this.showUserEditForm = !this.showUserEditForm;
        },
        createUserSaveForm() {
            this.createUser(this.userform);
            this.showUserEditForm = !this.showUserEditForm;
            // this.getObjectsForTemplate();
        },
        changeUser(user) {
            let answer = false
            let json_string = {}
            if (user.isForm) {
                answer = true;
                ["id", "username", "password", "first_name", "last_name", "email", "phone", "user_status"].forEach(function (elem) {
                    if (user[elem] != "" || elem == "user_status") {
                        json_string[elem] = user[elem];
                    }
                })
            } else {
                let operation = "disable" ? user.user_status : "enable"
                let msg = "Do you really want to " + operation + " this user?"
                answer = confirm(msg);
                json_string = {"id": user.id, "user_status": !user.user_status}
            }
            if (answer) {
                const requestOptions = {
                    method: 'PUT',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(json_string)
                };
                fetch('/api/users/', requestOptions)
                    .then(async response => {
                        const data = await response.json();
                        // check for error response
                        if (!response.ok) {
                            // get error message from body or default to response status
                            const error = (data && data.message) || response.status;
                            return Promise.reject(error);
                        }
                        this.getObjectsForTemplate();
                    })
                    .catch(error => {
                        this.errorMessage = error;
                        console.error('There was an error!', error);
                    });

            }
        },
        editUserShowForm(user) {
            this.userform.id = user.id;
            this.userform.username = user.username;
            this.userform.first_name = user.first_name;
            this.userform.last_name = user.last_name;
            this.userform.email = user.email;
            this.userform.phone = user.phone;
            this.userform.user_status = user.user_status;

            this.showUserEditForm = !this.showUserEditForm;
        },
        editUserSaveForm() {
            this.changeUser(this.userform);
            this.showUserEditForm = !this.showUserEditForm;
            // this.getObjectsForTemplate();
        },
        deleteUser(user_id) {
            let msg = "Do you really want to delete this user?"
            let answer = confirm(msg);
            if (answer) {
                const requestOptions = {
                    method: 'DELETE'
                };
                fetch('/api/users/' + user_id, requestOptions)
                    .then(async response => {
                        const data = await response.json();
                        // check for error response
                        if (!response.ok) {
                            // get error message from body or default to response status
                            const error = (data && data.message) || response.status;
                            return Promise.reject(error);
                        }
                        this.getObjectsForTemplate();
                    })
                    .catch(error => {
                        this.errorMessage = error;
                        console.error('There was an error!', error);
                    });
            }
        },
        initDB() {
            this.getRequest('/dbinit/', "Do you really want to Init DB? All data will be lost...");
            this.getObjectsForTemplate();
        },
        addDemoData() {
            this.getRequest('/demodb/', "Do you really want add demo data?");
            this.getObjectsForTemplate();
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
        getObjectsForTemplate() {
            let page_template = document.querySelector('meta[name="page_template"]').content;
            switch (page_template) {
                case "dashboard":
                    this.getLinks();
                    break;
                case "users":
                    this.getUsers();
                    break;
            }
        }
    },
    beforeMount(){
        this.getObjectsForTemplate();
    },
}).mount('.vueapp');