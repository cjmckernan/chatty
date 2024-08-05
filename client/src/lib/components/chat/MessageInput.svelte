
<script>
    export let onSendMessage = () => {};
    let newMessage = '';

    function sendMessage() {
        if (newMessage.trim() !== '') {
            onSendMessage(newMessage); 
            newMessage = ''; 
        }
    }

    function handleKeyPress(event) {
        if (event.key === 'Enter') {
            sendMessage();
        }
    }
</script>

<div class="message-input">
    <input 
        type="text" 
        bind:value={newMessage} 
        placeholder="Enter your message..." 
        on:keypress={handleKeyPress} 
     />
    <button on:click={sendMessage}>Send</button>
</div>


<style>
    .message-input {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-top: 10px;
        background-color: #2C2C2C; 
        padding: 10px;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    }

    input {
        flex: 1;
        padding: 10px 12px;
        background-color: #3A3A3A; 
        border: none;
        border-radius: 4px;
        color: #E0E0E0; 
        font-size: 14px;
        outline: none;
        transition: background-color 0.2s ease;
    }

    input::placeholder {
        color: #888888; 
    }

    input:focus {
        background-color: #4A4A4A;
    }

    button {
        padding: 10px 16px;
        background-color: #82AAFF; 
        border: none;
        border-radius: 4px;
        color: white;
        font-weight: bold;
        cursor: pointer;
        transition: background-color 0.2s ease, transform 0.2s ease;
    }

    button:hover {
        background-color: #699CFF; 
    }

    button:active {
        transform: scale(0.95); 
    }
</style>

