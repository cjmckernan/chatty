<script>
  import { sessionId, setSessionId } from '$lib/stores/auth.js';  

  export let closeModal;
  let username = '';
  let password = '';
  let passwordConfirm = '';
  let errorMessage = '';

  const handleSubmit = async () => {
    try {
      const response = await fetch('http://localhost:9009/user/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password, passwordConfirm }),
      });

      if (response.ok) {
        const data = await response.json();
        const sessionId = response.headers.get('X-Session-ID');
        setSessionId(sessionId); 
        closeModal(); 
      } else {
        const errorData = await response.json();
        errorMessage = errorData.error || 'Account creation failed';
      }
    } catch (error) {
      console.error('Error:', error);
      errorMessage = 'An error occurred. Please try again.';
    }
  };
</script>

<div class="modal">
  <div class="modal-content">
    <h2>Sign Up</h2>
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="username">Username</label>
        <input id="username" type="text" bind:value={username} required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input id="password" type="password" bind:value={password} required />
      </div>
      <div class="form-group">
        <label for="password-Confirm">Confirm Password</label>
        <input id="password-Confirm" type="password" bind:value={passwordConfirm} required />
      </div>
      <button type="submit">Sign Up</button>
    </form>
    <button on:click={closeModal} class="close-button">Close</button>
  </div>
</div>

<style>
  .modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .modal-content {
    background: #1E1E1E;
    padding: 20px;
    border-radius: 8px;
    color: #E0E0E0;
    width: 90%;
    max-width: 400px; /* Set the maximum width of the modal */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  .form-group {
    margin-bottom: 15px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
  }

  .form-group input {
    width: 100%;
    padding: 8px;
    border: none;
    border-radius: 4px;
  }

  button {
    background-color: #82AAFF;
    color: white;
    padding: 10px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    width: 100%;
  }

  .close-button {
    background: transparent;
    color: #E0E0E0;
    margin-top: 10px;
    cursor: pointer;
    width: 100%;
    text-align: center;
  }
</style>

