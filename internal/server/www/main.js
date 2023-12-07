// Example starter JavaScript for disabling form submissions if there are invalid fields
(function () {
    'use strict'

    function getRandomInt() {
        return Math.floor(Math.random() * 100000000000);
    }

    const formParamsWithDefaults = {
        "seed": getRandomInt(),
        "numServices": 3,
        "minReplicas": 2,
        "maxReplicas": 2,
        "percentEdge": 50,
        "yaml": null,
        "k8sNamespace": "microservice-mesh",
        "k8sApp": "api-play",
        "defineContent": JSON.stringify({services: [{"replicas": 2, "edges": [1]}, {"replicas": 2}]})
    }

    async function sendRequestAndPopulate(url, body, alertContainer, responseContainer, responseClasses) {
        let params = {}
        if (body) {
            params = {
                method: 'POST',
                headers: {
                    "Content-Type": "application/json",
                },
                body: body,
            }
        }
        let response = await fetch(url, params)
        if (!response.ok) {
            alertContainer.removeChild(alertContainer.firstChild)
            let pre = document.createElement('pre');
            pre.textContent = JSON.stringify(await response.json(), null, 4);
            alertContainer.appendChild(pre);
            alertContainer.classList.remove('d-none');
            return false
        }
        let pre = document.createElement("pre");
        if (responseClasses) {
            pre.classList.add(responseClasses)
        }
        pre.textContent = await response.text();
        responseContainer.querySelector('.highlight').appendChild(pre);
        responseContainer.classList.remove('d-none');
        return true
    }

    document.onreadystatechange = () => {
        if (document.readyState === "complete") {
            // Setup randomForm
            const randomResponseContainer = document.querySelector("#random-response-container");
            const randomK8sResponseContainer = document.querySelector("#random-k8s-response-container");
            const defineResponseContainer = document.querySelector("#define-response-container");
            const defineK8sResponseContainer = document.querySelector("#define-k8s-response-container");

            document.querySelectorAll(".nav-link").forEach((elt) => elt.addEventListener('click', function(e) {
                window.location.href = e.target.href
            }))
            if (location.hash === "#define") {
                document.querySelector("#v-pills-define").classList.add("active", "show")
                document.querySelector("#v-pills-define-tab").classList.add("active")

                document.querySelector("#v-pills-random").classList.remove("active", "show")
                document.querySelector("#v-pills-random-tab").classList.remove("active")
            } else {
                document.querySelector("#v-pills-random").classList.add("active", "show")
                document.querySelector("#v-pills-random-tab").classList.add("active")

                document.querySelector("#v-pills-define-tab").classList.remove("active")
                document.querySelector("#v-pills-define").classList.remove("active", "show")
            }

            const randomForm = document.querySelector('#form-random')
            const randomFormAlert = randomForm.querySelector(".failure-alert")
            const defineForm = document.querySelector('#form-define')
            const defineFormAlert = defineForm.querySelector(".failure-alert")
            randomForm.querySelector('button.seed-refresh').addEventListener('click', function (event) {
                event.preventDefault();
                randomForm.querySelector("input[name='seed']").value = getRandomInt()
            });
            const copyToClipboardList = document.querySelectorAll('button.btn-clipboard')
            copyToClipboardList.forEach(tooltipTriggerEl => {
                tooltipTriggerEl.addEventListener('click', (event) => {
                    navigator.clipboard.writeText(tooltipTriggerEl.parentElement.parentElement.querySelector('.highlight').innerText);
                })
                const tt = new bootstrap.Tooltip(tooltipTriggerEl);
                tt.setContent({'.tooltip-inner': 'Copy to clipboard'});
                return tt;
            });

            // Watch changes to url to call the api
            let previousUrl;
            let handlePageChange = async () => {
                if (location.href !== previousUrl) {
                    // Set form values to whatever is in the url or default
                    let url = new URL(document.location);
                    let apiURL = new URL(document.location);
                    let form = url.hash === '#define' ? defineForm : randomForm;
                    for (const key in  formParamsWithDefaults) {
                        let elt = form.querySelector(`input[name='${key}'], select[name='${key}'], textarea[name='${key}']`);
                        if (!elt) {
                            continue
                        }
                        if (elt.type === "checkbox") {
                            elt.checked = url.searchParams.has(key);
                            elt.value = true;
                        } else {
                            elt.value = url.searchParams.has(key) ? url.searchParams.get(key) : formParamsWithDefaults[key]
                        }
                        if (elt.value !== '') {
                            apiURL.searchParams.set(key, elt.value);
                        }
                    }
                    let newUrl = new URL(location.href);
                    let oldUrl = previousUrl && new URL(previousUrl);
                    previousUrl = location.href;
                    if (oldUrl?.hash === '#random') {
                        let rrhl = randomResponseContainer.querySelector('.highlight');
                        if (rrhl.firstChild) {
                            rrhl.removeChild(rrhl.firstChild);
                        }

                        let k8srrhl = randomK8sResponseContainer.querySelector('.highlight');
                        if (k8srrhl.firstChild) {
                            k8srrhl.removeChild(k8srrhl.firstChild);
                        }
                        randomResponseContainer.classList.add('d-none');
                        randomK8sResponseContainer.classList.add('d-none');
                    } else if (oldUrl?.hash === '#define') {
                        let rrhl = defineResponseContainer.querySelector('.highlight');
                        if (rrhl.firstChild) {
                            rrhl.removeChild(rrhl.firstChild);
                        }

                        let k8srrhl = defineK8sResponseContainer.querySelector('.highlight');
                        if (k8srrhl.firstChild) {
                            k8srrhl.removeChild(k8srrhl.firstChild);
                        }
                        defineResponseContainer.classList.add('d-none');
                        defineK8sResponseContainer.classList.add('d-none');
                    }
                    if (newUrl.hash === '#random') {
                        const asYaml = newUrl.searchParams.has('yaml');
                        if (asYaml) {
                            apiURL.pathname = '/api/random.yaml';
                        } else {
                            apiURL.pathname = '/api/random.mmd';
                        }
                        let success = await sendRequestAndPopulate(apiURL, undefined, randomFormAlert, randomResponseContainer, asYaml ? ["yaml"] : ["mermaid"])
                        if (success) {
                            let k8sUrl = new URL(apiURL);
                            k8sUrl.pathname = '/api/random.yaml';
                            k8sUrl.searchParams.set("k8s", 'true');
                            await sendRequestAndPopulate(k8sUrl, undefined, randomFormAlert, randomK8sResponseContainer)
                        }
                    } else if (newUrl.hash === "#define") {
                        const asYaml = newUrl.searchParams.has('yaml');
                        if (asYaml) {
                            apiURL.pathname = '/api/define.yaml';
                        } else {
                            apiURL.pathname = '/api/define.mmd';
                        }
                        let payload = apiURL.searchParams.get('defineContent')
                        let success = await sendRequestAndPopulate(apiURL, payload, defineFormAlert, defineResponseContainer, asYaml ? ["yaml"] : ["mermaid"])
                        if (success) {
                            let k8sUrl = new URL(apiURL);
                            k8sUrl.pathname = '/api/define.yaml';
                            k8sUrl.searchParams.set("k8s", 'true');
                            await sendRequestAndPopulate(k8sUrl, payload, defineFormAlert, defineK8sResponseContainer)
                        }
                    }
                }
            };
            window.addEventListener('hashchange', handlePageChange);
            window.addEventListener('popstate', handlePageChange);
            handlePageChange()
            document.querySelector(".random-params-group").addEventListener("change", (event) => {
                // Sync strongly min and max to avoid invalid states
                let minReplicas = event.currentTarget.querySelector("input[name='minReplicas']");
                let maxReplicas = event.currentTarget.querySelector("input[name='maxReplicas']");
                let invalid = parseInt(minReplicas.value || '0', 10) > parseInt(maxReplicas.value || '0', 10);
                if (event.target === minReplicas && invalid) {
                    maxReplicas.value = minReplicas.value;
                }
                if (event.target === maxReplicas && invalid) {
                    minReplicas.value = maxReplicas.value;
                }
            })

            // Loop over them and prevent submission
            Array.prototype.slice.call(document.getElementsByTagName('form'))
                .forEach(function (form) {
                    form.addEventListener('submit', function (event) {
                        if (event.target.classList.contains('needs-validation') && !event.target.checkValidity()) {
                            event.preventDefault()
                            event.stopPropagation()
                        }
                        form.classList.add('was-validated')
                    }, false)
                })
        }
    };
})()
