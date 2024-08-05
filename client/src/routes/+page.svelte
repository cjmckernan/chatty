<script>
    import Header from '$lib/components/layout/Header.svelte'
    import { onMount } from 'svelte';
    let response;

    onMount(async () => {
        await testRequest();
    });

    async function testRequest() {
        try {
            const res = await fetch('http://localhost:9009/');
            if (res.ok) {
                const data = await res.json();  // Assuming the response is JSON
                response = data;
            } else {
                console.error('Failed to fetch topics:', res.statusText);
            }
        } catch (error) {
            console.error('Error fetching topics:', error);
        }
    }
</script>

<div>
  <Header />
  {#if response}
    <p>{response.message}</p>
  {/if}
</div>

