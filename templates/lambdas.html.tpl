{{ define "content" }}
    <script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/11.6.0/highlight.min.js"></script>
    <div>
        <p>
            You can now setup lambda functions. Just write or paste lua code into the area below 
            and save a functrion with a given name. You can overwrite functions to fix issues. 
        </p>
        <input id="name-input" class="flex-1 w-32 p-2 border-2 border-raisinblack rounded" type="text" placeholder="Name" />
        <div class="relative h-72">
            <textarea id="editing" class="border-0 outline-0"></textarea>
            <pre id="highlighting" aria-hidden="true"><code class="language-lua max-h-60" id="highlighting-content">
            </code></pre>
        </div>
        <button data-namespace="{{ .Namespace }}" class="px-4 py-2 rounded-full bg-rosybrown" type="button" id="save-btn">Save</button>
    </div>
    <script>
        hljs.highlightAll();

        const codeElement = document.querySelector('#highlighting-content');
        const textAreaElement = document.querySelector('#editing');
        textAreaElement.addEventListener('input', () => {
            codeElement.innerHTML = textAreaElement.value.replace(new RegExp("&", "g"), "&").replace(new RegExp("<", "g"), "<");
            hljs.highlightElement(codeElement);
        });

        const nameInput = document.querySelector('#name-input');
        const saveBtn = document.querySelector('#save-btn');
        saveBtn.addEventListener('click', () => {
            fetch(`/api/lambda/${saveBtn.dataset.namespace}/${nameInput.value}`, {
                headers: {
                    'Content-Type': 'text/plain'
                },
                method: 'POST',
                body: textAreaElement.value,
            }).then(async (resp) => {
                textAreaElement.value = "\n";
                codeElement.innerHTML = "\n";
                hljs.highlightElement(codeElement);
            });
        });
    </script>
{{ end }}