
<script>
    import Header from '$lib/components/layout/Header.svelte';
    import { sessionId } from '$lib/stores/auth.js';
    import { username, setUsername } from '$lib/stores/user.js';  

    let response;
    let currentSessionId;
    let pageUsername;
    $: currentSessionId = $sessionId;
     
    $: {
        if (currentSessionId) {
            testRequest();
        }
    }


    $: pageUsername = $username;

    async function testRequest() {
        try {
            const res = await fetch('http://localhost:9009/', {
                method: 'GET',
                headers: {
                    'X-Session-ID': currentSessionId, // Add the session ID to the headers
                }
            });

            if (res.ok) {
                // Get the username from the response headers
                pageUsername = res.headers.get('X-User-Name');
                setUsername(pageUsername)
                const data = await res.json();  // Assuming the response is JSON
                response = data;
            } else {
                console.error('Failed to fetch data:', res.statusText);
            }
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }
</script>

<div>
  <Header />
  {#if response}
    <p>{response.message} {pageUsername}</p>
  {/if}
</div>

