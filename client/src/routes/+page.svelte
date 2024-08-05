<script>
    import Header from '$lib/components/layout/Header.svelte';
    import Messages from '$lib/components/chat/Messages.svelte';
    import Topics from '$lib/components/chat/Topics.svelte';
    import MessageInput from '$lib/components/chat/MessageInput.svelte';

    import { sessionId } from '$lib/stores/auth.js';
    import { username, setUsername } from '$lib/stores/user.js';  

    let response;
    let currentSessionId;
    let pageUsername;
    let socket;
    let messages = []; 
    let selectedTopic = ''; 
    $: currentSessionId = $sessionId;
     
    $: {
        if (currentSessionId) {
            testRequest();
            connectWebSocket();
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
                pageUsername = res.headers.get('X-User-Name');
                setUsername(pageUsername);
                const data = await res.json();  // Assuming the response is JSON
                response = data;
            } else {
                console.error('Failed to fetch data:', res.statusText);
            }
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }

    function connectWebSocket() {
        if (!currentSessionId) return;

        socket = new WebSocket(`ws://localhost:9009/ws?sessionId=${currentSessionId}`);

        socket.onopen = () => {
            console.log('WebSocket connection established');
        };

        socket.onmessage = (event) => {
            console.log('Message from server:', event.data);
            const message = JSON.parse(event.data);
            messages = [...messages, message]; // Update messages array
        };

        socket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    async function fetchMessagesByTopic(topic) {
        try {
            const res = await fetch(`http://localhost:9009/messages/${topic}`, {
                method: 'GET',
                headers: {
                    'X-Session-ID': currentSessionId,
                }
            });

            if (res.ok) {
                const data = await res.json();
                messages = data.messages; // Update the messages list with fetched messages
            } else {
                console.error('Failed to fetch messages:', res.statusText);
            }
        } catch (error) {
            console.error('Error fetching messages:', error);
        }
    }

    function onSelectTopic(topic) {
        selectedTopic = topic;
        messages = []; // Clear messages when topic changes
        fetchMessagesByTopic(topic);
    }

    function sendMessageToServer(message) {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ topic: selectedTopic, username: pageUsername, text: message }));
        }
    }
</script>

<div>
  <Header />
  {#if currentSessionId}
    <div class="container">
      <Topics bind:selectedTopic onSelectTopic={onSelectTopic} />
      <Messages {messages} />
      <MessageInput onSendMessage={sendMessageToServer} /> 
    </div>
  {/if}
</div>

<style>
  .container {
    padding: 20px;
  }
</style>

