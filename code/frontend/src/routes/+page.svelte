<script lang="ts">
	import { goto } from "$app/navigation";
	import Registration from "$lib/components/Registration.svelte";
	import Button from "$lib/components/ui/button/button.svelte";

    let username = $state<String>("");
    let password = $state<String>("");
    let registerFlag = $state<boolean>(false);
  
    async function handleSubmit() {
      // Handle login logic here
      console.log('Login attempted with:', { username, password });
      
      try {
        const response = await fetch('http://localhost/api/auth/login', {
            method: 'POST',
            headers: {
            'Content-Type': 'application/json',
            },
            body: JSON.stringify({ "username" : username, "password" : password }),
        });

        response.ok ? console.log('Login successful') : console.log('Login failed');
        const jsonData = await response.json();

        let date = new Date();
        date.setTime(jsonData.Data.expiresAt * 1000);
        let expires = "; expires=" + date.toUTCString();
        let newCookie = `token=${jsonData.Data.token}${expires};path=/`;
        console.log(newCookie);
        document.cookie = `token=${jsonData.Data.token}${expires};path=/`;
        goto("/archive");

      } catch (error) {
        console.error('Login failed:', error);
        return;
      }
    }
  </script>
  
  {#if registerFlag}
    <Registration bind:registerFlag={registerFlag}/>
  {:else}
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900">
    <div class="w-full max-w-md">
      <form onsubmit={handleSubmit} class="bg-white dark:bg-gray-800 shadow-md rounded-lg px-8 pt-6 pb-8 mb-4">
        <h2 class="text-2xl font-bold mb-6 text-gray-900 dark:text-gray-100">Login</h2>
        
        <div class="mb-4">
          <label class="block text-gray-700 dark:text-gray-300 text-sm font-bold mb-2" for="username">
            Username
          </label>
          <input 
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight focus:outline-none focus:shadow-outline dark:bg-gray-700 dark:border-gray-600"
            id="username"
            type="text"
            placeholder="Username"
            bind:value={username}
            required
          />
        </div>
        
        <div class="mb-6">
          <label class="block text-gray-700 dark:text-gray-300 text-sm font-bold mb-2" for="password">
            Password
          </label>
          <input 
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 mb-3 leading-tight focus:outline-none focus:shadow-outline dark:bg-gray-700 dark:border-gray-600"
            id="password"
            type="password"
            placeholder="******************"
            bind:value={password}
            required
          />
        </div>
        
        <div class="flex items-center justify-between">
          <button 
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline transition-colors duration-200"
            type="submit"
          >
            Sign In
          </button>
          <Button 
            class="inline-block align-baseline font-bold text-sm text-blue-500 bg-transparent hover:bg-transparent hover:underline"
            onclick={() => registerFlag = true}
          >
            Register
        </Button>
        </div>
      </form>
      <p class="text-center text-gray-500 text-xs dark:text-gray-400">
        &copy;2023 Your Company. All rights reserved.
      </p>
    </div>
  </div>
  {/if}
  