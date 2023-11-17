// Example starter JavaScript for disabling form submissions if there are invalid fields
(function () {
    'use strict'

    function getRandomInt() {
        return Math.floor(Math.random() * 100000000000);
    }

    const randomParamsWithDefaults = {
        "seed": getRandomInt(),
        "numServices": 3,
        "minReplicas": 2,
        "maxReplicas": 2,
        "percentEdge": 50,
        "yaml": null,
        "k8sNamespace": "microservice-mesh",
        "k8sApp": "api-play",
    }

    async function sendRequestAndPopulate(url, alertContainer, responseContainer, responseClasses) {
        let response = await fetch(url)
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

            const randomForm = document.querySelector('#form-random')
            const randomFormAlert = randomForm.querySelector(".failure-alert")
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
            let observer = new MutationObserver(async function (mutations) {
                if (location.href !== previousUrl) {
                    // Set form values to whatever is in the url or default
                    let url = new URL(document.location);
                    let apiURL = new URL(document.location);
                    for (const key in randomParamsWithDefaults) {
                        let elt = randomForm.querySelector(`input[name='${key}'], select[name='${key}']`);
                        if (!elt) {
                            console.error(`No input in form ${key}`);
                            continue
                        }
                        if (elt.type === "checkbox") {
                            elt.checked = url.searchParams.has(key);
                            elt.value = true;
                        } else {
                            elt.value = url.searchParams.has(key) ? url.searchParams.get(key) : randomParamsWithDefaults[key]
                        }
                        if (elt.value !== '') {
                            apiURL.searchParams.set(key, elt.value);
                        }
                    }
                    let newUrl = new URL(location.href);
                    let oldUrl = previousUrl && new URL(previousUrl);
                    previousUrl = location.href;
                    console.log(`URL changed from ${previousUrl} to ${location.href}`);
                    if (oldUrl?.hash === '#random') {
                        let rrhl = randomResponseContainer.querySelector('.highlight');
                        rrhl.removeChild(rrhl.firstChild);
                        let k8srrhl = randomK8sResponseContainer.querySelector('.highlight');
                        k8srrhl.removeChild(k8srrhl.firstChild);
                        randomResponseContainer.classList.add('d-none');
                        randomK8sResponseContainer.classList.add('d-none');
                    }
                    if (newUrl.hash === '#random') {
                        const asYaml = newUrl.searchParams.has('yaml');
                        if (asYaml) {
                            apiURL.pathname = '/api/random.yaml';
                        } else {
                            apiURL.pathname = '/api/random.mmd';
                        }
                        let success = await sendRequestAndPopulate(apiURL, randomFormAlert, randomResponseContainer, asYaml ? ["yaml"] : ["mermaid"])
                        if (success) {
                            let k8sUrl = new URL(apiURL);
                            k8sUrl.pathname = '/api/random.yaml';
                            k8sUrl.searchParams.set("k8s", 'true');
                            await sendRequestAndPopulate(k8sUrl, randomFormAlert, randomK8sResponseContainer)
                        }
                    }
                }
            });
            observer.observe(document.body, {childList: true, subtree: true});
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
