<script lang="ts">
	import { decrypt } from "$lib/crypt";
	import type { EncFile } from "$lib/EncFile/EncFile";
	import { getContext } from "svelte";
	import Button from "./ui/button/button.svelte";
	import { X } from "lucide-svelte";
    
    const user : {username : string; email: string, id: string, token:string} = getContext('user');
    let {isOpen = $bindable(), photo} : {isOpen: boolean, photo : EncFile} = $props();
    let isImageDecrypted = $state<boolean>(false);
    let url = $state<string>("");
    
    let password = $state<string>("");
    
    async function handleDecrypt() {
      // Handle decryption logic here
      console.log("Decrypting with password:", password);

      try {
        const res = await fetch('http://localhost/api/archive/storage/'+photo.filename, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + user.token
            },
          });  
        console.log('Response:', res); 
        const body = await res.text(); 
        console.log('Body:', body);

        const decryptedPhoto = await decrypt(body, password)

        let blob = new Blob([decryptedPhoto], {type: photo.extension});
        url = URL.createObjectURL(blob);
        isImageDecrypted = true;

        } catch (error) {
            console.log('Error:', error); 
        }

      // Decrypt the photo
      //const decryptedPhoto = await decrypt(filepassword);
    }
    
    function handleClose() {
      isOpen = false;
    }

  </script>
  
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
    <div class="bg-[#1A1F2C] border border-grey-800 rounded-3xl p-8 w-full max-w-3xl h-[80%] relative">
      <!-- Close button -->

      <Button class="absolute top-4 right-4 text-white hover:text-white/80 bg-transparent" onclick={handleClose}>
        <X />
      </Button>

      <!-- Modal content -->
      <div class="space-y-4 h-[100%]">
        <!-- Photo name -->
        <h2 class="text-white text-xl font-light">{photo.filename}</h2>
        
        <!-- Photo preview area -->
        <div class="border-2 border-gray-700 h-[80%] rounded-xl aspect-3/4 mb-6 overflow-hidden">
            {#if isImageDecrypted}
              <img class="w-full h-full object-contain bg-gray-800 rounded-xl" src={url} alt={photo.filename} />
            {:else}
              <div class="w-full h-full flex items-center justify-center text-gray-400 bg-gray-800">
                Photo preview
              </div>
            {/if}
        </div>
        
        <!-- Password and decrypt section -->
        <div class="flex gap-4">
          <input
            type="password"
            placeholder="insert password"
            bind:value={password}
            class="flex-1 bg-transparent border border-white rounded-lg px-4 py-2 text-white placeholder-white/70 font-light focus:outline-none focus:ring-2 focus:ring-white/50"
          />
          <button
            onclick={handleDecrypt}
            class="bg-transparent border border-white rounded-lg px-6 py-2 text-white font-light hover:bg-white/10 transition-colors"
          >
            Decrypt
          </button>
        </div>
      </div>
    </div>
  </div>