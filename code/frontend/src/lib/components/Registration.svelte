
<script lang="ts">
	import { goto } from "$app/navigation";
	import Button from "./ui/button/button.svelte";

    let username = $state('');
    let email = $state('');
    let password = $state('');
    let confirmPassword =$state('');
    let error = $state('');

    let { registerFlag = $bindable() } : {registerFlag: boolean} = $props();
  
    async function handleSubmit(event: SubmitEvent) {
      event.preventDefault();
      
      if (password !== confirmPassword) {
        error = 'Passwords do not match';
        return;
      }
  
      try {
        const response = await fetch('http://localhost/api/auth/register', {
            method: 'POST',
            headers: {
            'Content-Type': 'application/json',
            },
            body: JSON.stringify({ "username" : username, "password" : password , "email" : email }),
        });

        if (response.ok) {
            console.log('Registration successful');
            registerFlag = false;
        } else console.log('Login failed');

        // let date = new Date();
        // date.setTime(jsonData.Data.expiresAt * 1000);
        // let expires = "; expires=" + date.toUTCString();
        // let newCookie = `token=${jsonData.Data.token}${expires};path=/`;
        // console.log(newCookie);
        // document.cookie = `token=${jsonData.Data.token}${expires};path=/`;
        // goto("/archive");

      } catch (error) {
        console.error('Login failed:', error);
        return;
      }


      error = '';
      console.log('Form submitted:', { username, email, password });
    }
  </script>


<div class="flex min-h-screen flex-col items-center justify-center bg-gray-50">
    <div class="w-full max-w-[400px] rounded-lg bg-white p-6 shadow-sm">
      <div class="space-y-6">
        <div>
          <h1 class="text-xl font-semibold text-gray-900">Register</h1>
        </div>
  
        <form class="space-y-4" onsubmit={handleSubmit}>
          <div class="space-y-2">
            <label for="username" class="block text-sm font-medium text-gray-900">
              Username
            </label>
            <input
              id="username"
              bind:value={username}
              placeholder="Username"
              required
              type="text"
              class="w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
  
          <div class="space-y-2">
            <label for="email" class="block text-sm font-medium text-gray-900">
              Email
            </label>
            <input
              id="email"
              bind:value={email}
              placeholder="name@example.com"
              required
              type="email"
              class="w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
  
          <div class="space-y-2">
            <label for="password" class="block text-sm font-medium text-gray-900">
              Password
            </label>
            <input
              id="password"
              bind:value={password}
              required
              type="password"
              placeholder="••••••••••••••"
              class="w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
  
          <div class="space-y-2">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-900">
              Confirm Password
            </label>
            <input
              id="confirmPassword"
              bind:value={confirmPassword}
              required
              type="password"
              placeholder="••••••••••••••"
              class="w-full rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
  
          {#if error}
            <p class="text-sm text-red-500">{error}</p>
          {/if}
  
          <button
            type="submit"
            class="w-full rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            Create Account
          </button>
        </form>
  
        <div class="text-center text-sm text-gray-500">
          Already have an account?
          <Button onclick={() => registerFlag = false} class="text-blue-600 bg-transparent hover:bg-transparent hover:underline">Sign in</Button>
        </div>
      </div>
    </div>
    <p class="mt-8 text-center text-sm text-gray-500">
      ©2023 Your Company. All rights reserved.
    </p>
  </div>