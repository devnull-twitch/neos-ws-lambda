{{ define "content" }}
    <section>
        <h1 class="font-bold text-lg pt-6 pb-2">Lambda Server</h1>
        <p>
            Welcome. You can start a session by clicking the button below. <br />
            You can then setup multiple lambda functions. Lambda functions are provided 
            written in lua.
        </p>
    </section>
    <section>
        <h3 class="font-bold pt-3 pb-1">Initial persistnt variables</h3>
        <p>
            Your lambda functions may load and save any number of variables for later
            reuse by other lambda functions.<br />
            Here you can define the initial variables and values.
        </p>
        <div class="js-varval-container container w-1/2 py-6">
            <button class="px-4 py-2 mt-2 rounded-full bg-maximumblue" type="button" id="plus-btn">+</button>
        </div>
        <button class="px-4 py-2 rounded-full bg-rosybrown" type="button" id="start-btn">Start Session</button>
    </section>
    <template id="variablerow">
        <div class="flex py-1 js-varrow">
            <input type="text" class="flex-1 w-32 p-2 border-2 border-raisinblack rounded mr-3" placeholder="Variable name" />
            <input type="text" class="flex-1 w-32 p-2 border-2 border-raisinblack rounded" placeholder="Variable value" />
        </div>
    </template>
    <script>
        document.querySelector('#start-btn').addEventListener('click', () => {
            let args = [];
            document.querySelectorAll('.js-varrow').forEach((row) => {
                args.push(`${row.querySelector('input:first-child').value}=${row.querySelector('input:nth-child(2)').value}`);
            });

            fetch('/api/session', {
                headers: {
                    'Content-Type': 'text/plain'
                },
                method: 'POST',
                body: args.join('|'),
            }).then(async (resp) => {
                return resp.text();
            }).then((r) => {
                window.location.href = `/session/${r}`;
            });
        });

        const template = document.querySelector('#variablerow');
        const varContainer = document.querySelector('.js-varval-container');

        function addVarBlock() {
            const newRow = template.content.cloneNode(true);
            varContainer.insertBefore(newRow, varContainer.firstChild);
        }

        addVarBlock();

        document.querySelector('#plus-btn').addEventListener('click', addVarBlock);
    </script>
{{ end }}