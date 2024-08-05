<script>
    import { username } from '$lib/stores/user.js';
    export let messages = [];

    // Reactive variable to store the current username
    let currentUser;
    $: currentUser = $username;
</script>

<div class="message-screen">
    {#if messages.length > 0}
        {#each messages as message, i}
            <div class="message {message.username === currentUser ? 'message-own' : ''}">
                <div class="message-content">
                    <p class="message-text">{message.text}</p>
                    <p class="message-username">{message.username}</p>
                </div>
            </div>
        {/each}
    {:else}
        <p class="no-messages">No messages yet.</p>
    {/if}
</div>

<style>
    .message-screen {
        background-color: #1E1E1E; /* Dark background */
        border-radius: 8px;
        padding: 15px;
        height: 300px;
        overflow-y: auto;
        box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.3);
        margin-bottom: 15px;
    }

    .message {
        margin-bottom: 10px;
        display: flex;
        align-items: center;
    }

    .message-own {
        justify-content: flex-end;
    }

    .message-content {
        max-width: 75%;
        background-color: #3A3A3A; /* Dark gray for message bubble */
        padding: 10px 15px;
        border-radius: 15px;
        color: #E0E0E0; /* Light text color */
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        position: relative;
    }

    .message-own .message-content {
        background-color: #82AAFF; /* Blue for own messages */
        color: #1E1E1E; /* Dark text color */
    }

    .message-text {
        margin: 0;
        font-size: 14px;
    }

    .message-username {
        margin-top: 5px;
        font-size: 12px;
        color: #888888; /* Muted color for username */
        text-align: right;
    }

    .message-own .message-username {
        color: #1E1E1E; /* Dark color for username in own messages */
    }

    .no-messages {
        color: #888888; /* Muted text color for "No messages" */
        text-align: center;
        margin-top: 20px;
    }
</style>

