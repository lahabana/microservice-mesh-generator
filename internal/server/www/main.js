// Example starter JavaScript for disabling form submissions if there are invalid fields
(function () {
    'use strict'

    const randomParamsWithDefaults = {
        "seed": new Date().getTime(),
        "numServices": 3,
        "minReplicas": 2,
        "maxReplicas": 2,
        "yaml": null,
    }

    document.onreadystatechange = () => {
        if (document.readyState === "complete") {
            const randomResponseContainer = document.querySelector("#random-response-container");
            randomResponseContainer.classList.add('d-none');
            const randomForm = document.querySelector("#form-random")
            const randomFormAlert = randomForm.querySelector(".failure-alert")

            // Watch changes to url to call the api
            let previousUrl;
            let observer = new MutationObserver(async function (mutations) {
                if (location.href !== previousUrl) {
                    // Set form values to whatever is in the url or default
                    let url = new URL(document.location);
                    let apiURL = new URL(document.location);
                    for (const key in randomParamsWithDefaults) {
                        let elt = randomForm.querySelector(`input[name='${key}']`);
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
                        randomResponseContainer.removeChild(randomResponseContainer.firstChild);
                        randomResponseContainer.classList.add('d-none');
                    }
                    if (newUrl.hash === '#random') {
                        //
                        const asYaml = newUrl.searchParams.has('yaml');
                        if (asYaml) {
                            apiURL.pathname = '/api/random.yaml';
                        } else {
                            apiURL.pathname = '/api/random.mmd';
                        }
                        let response = await fetch(apiURL)
                        if (!response.ok) {
                            randomFormAlert.removeChild(randomFormAlert.firstChild)
                            let pre = document.createElement('pre');
                            pre.textContent = JSON.stringify(await response.json(), null, 4);
                            randomFormAlert.appendChild(pre);
                            randomFormAlert.classList.remove('d-none');
                        } else {
                            let pre = document.createElement("pre");
                            if (asYaml) {
                                pre.classList.add("yaml");
                            } else {
                                pre.classList.add("mermaid");
                            }
                            pre.textContent = await response.text();
                            randomResponseContainer.appendChild(pre);
                            randomResponseContainer.classList.remove('d-none');
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
