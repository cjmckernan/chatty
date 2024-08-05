<script>
    export let topics = [];
    export let selectedTopic = '';
    export let onSelectTopic = () => {};

    function handleSelect(event) {
        onSelectTopic(event.target.value);
    }
</script>

<div class="topics-container">
    <label for="topics">Select Topic:</label>
    <select id="topics" bind:value={selectedTopic} on:change={handleSelect}>
        {#each topics as topic}
            <option value={topic}>{topic}</option>
        {/each}
    </select>
</div>

<style>
    .topics-container {
        margin-bottom: 15px;
    }

    label {
        display: block;
        margin-bottom: 8px;
        color: #E0E0E0; /* Light text color */
        font-weight: bold;
        font-size: 14px;
    }

    select {
        width: 100%;
        padding: 10px 12px;
        background-color: #3A3A3A; /* Dark background for the select box */
        color: #E0E0E0; /* Light text color */
        border: none;
        border-radius: 4px;
        font-size: 14px;
        cursor: pointer;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        transition: background-color 0.2s ease, box-shadow 0.2s ease;
    }

    select:hover {
        background-color: #4A4A4A; /* Slightly lighter background on hover */
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
    }

    select:focus {
        outline: none;
        box-shadow: 0 0 0 3px rgba(130, 170, 255, 0.5); /* Blue focus ring */
    }

    option {
        background-color: #3A3A3A;
        color: #E0E0E0;
    }
</style>

