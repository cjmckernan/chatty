<script>
  import Signin from '$lib/components/forms/auth/Signin.svelte';
  import Signup from '$lib/components/forms/auth/Signup.svelte';
  import { sessionId, destroySessionId } from '$lib/stores/auth.js';
  import { username, setUsername } from '$lib/stores/user.js';  

  import { onDestroy } from 'svelte'
  let showSignin = false;
  let showSignup = false;

  let currentSessionId;
  let pageUsername;

  $: pageUsername = $username
  $: currentSessionId = $sessionId;
  
  $: { 
  }

  const closeModal = () => {
    showSignin = false;
    showSignup = false;
  };

  const destroySession = () => {
    destroySessionId()
  }


</script>

<nav class="navbar navbar-expand-lg">
  <div class="container-fluid">
    <a class="navbar-brand" href="#">Chatty</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
      <ul class="navbar-nav ms-auto">

        {#if !currentSessionId}
        <li class="nav-item">
          <a class="nav-link" href="#" on:click={() => (showSignin = true)}>Sign In</a>

        </li>
        <li class="nav-item">
          <a class="nav-link active" aria-current="page" href="#" on:click={() => (showSignup = true)}>Sign Up</a>
        </li>
        {:else}
          <li>
            <a class="nav-link active" aria-current="page" href="#" >{pageUsername}</a>
          </li>
          <li class="nav-item">
            <a class="nav-link active" aria-current="page" href="#" on:click={destroySession}>Logout</a>
          </li>

        {/if}
      </ul>
    </div>
  </div>
</nav>



{#if showSignin}
  <Signin {closeModal} />
{/if}

{#if showSignup}
  <Signup {closeModal} />
{/if}
<style>
  .navbar {
    background-color: #1E1E1E; 
    color: #E0E0E0; 
  }

  .navbar-brand, .nav-link {
    color: #E0E0E0 !important; 
  }

  .navbar-brand:hover, .nav-link:hover {
    color: #B3E5FC !important; 
  }

  .navbar-toggler {
    border-color: #E0E0E0; /* Make the toggler border match the text color */
  }

  .navbar-toggler-icon {
    background-image: url("data:image/svg+xml;charset=UTF8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 30 30'%3E%3Cpath stroke='%23E0E0E0' stroke-width='2' stroke-linecap='round' stroke-miterlimit='10' d='M4 7h22M4 15h22M4 23h22'/%3E%3C/svg%3E");
  }

  .navbar-collapse {
    background-color: #1E1E1E;  
  }

  .navbar-nav .nav-link {
    padding: 10px 15px;
  }
</style>


