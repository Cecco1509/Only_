<script lang="ts">
    import { getContext } from "svelte";
    import * as base64 from "byte-base64";
	import type { EncFile } from "$lib/EncFile/EncFile";

    let fileInput: HTMLInputElement;
    let previewUrl= $state<string>('');
    let password= $state<string>('');
    let filename= $state<string>('');
    
    let { isOpen = $bindable(), onSuccessfullUpload } : {isOpen : boolean; onSuccessfullUpload : (file : EncFile) => void } = $props()
    const user : {username : string; email: string, id: string, token:string} = getContext('user');

    const deriveAESCryptoKeyFromPassword = async (passWordCryptoKey : CryptoKey, salt : BufferSource) => {
				const cryptoKey = await window.crypto.subtle.deriveKey(
					{
						name: 'PBKDF2',
						salt,
						iterations: 100000,
						hash: 'SHA-256'
					},
					passWordCryptoKey,
					{ name: 'AES-GCM', length: 128 },
					false,
					['encrypt', 'decrypt']
				);
				return cryptoKey;
			};

    const encode = (text:string) => {
				return new TextEncoder().encode(text);
			};

    const generatePBKDF2CryptoKey = async (password : string) => {
				const passWordAsBytes = encode(password);
				const passWordCryptoKey = await window.crypto.subtle.importKey(
					'raw',
					passWordAsBytes,
					'PBKDF2',
					false,
					['deriveKey']
				);
				return passWordCryptoKey;
			};

    async function encrypt(buffer : ArrayBuffer, password : string) : Promise<Uint8Array> {
				const passWordCryptoKey : CryptoKey = await generatePBKDF2CryptoKey(password);
				const salt : Uint8Array = crypto.getRandomValues(new Uint8Array(16));
				const cryptoKey : CryptoKey = await deriveAESCryptoKeyFromPassword(passWordCryptoKey, salt);
				const iv : Uint8Array = crypto.getRandomValues(new Uint8Array(16));
				const encryptedBuffer = await window.crypto.subtle.encrypt(
					{
						name: 'AES-GCM',
						iv
					},
					cryptoKey,
					buffer
				);
				const encryptedData : Uint8Array = new Uint8Array(
					encryptedBuffer.byteLength + iv.byteLength + salt.byteLength
				);
				encryptedData.set(salt, 0);
				encryptedData.set(iv, salt.byteLength);
				encryptedData.set(new Uint8Array(encryptedBuffer), salt.byteLength + iv.byteLength);

				return encryptedData;
			}
    
    async function handleFileSelect(event: Event) {
      const file = (event.target as HTMLInputElement).files?.[0];
      if (!file) return;
      console.log(file)
      previewUrl = URL.createObjectURL(file);
    }
    
    function handleImport() {
      fileInput?.click();
    }
    
    async function handleSave() {
      const file = fileInput.files?.[0];
      if (!file || password.length < 6) return;
      console.log(file.name, (await file.arrayBuffer()).byteLength)
      const compressedEncrypted : Uint8Array = await encrypt(await file.arrayBuffer(), password);
      console.log("compressed: "+compressedEncrypted.byteLength)
      const b64encoded : string = base64.bytesToBase64(compressedEncrypted);
      console.log(b64encoded)

      try {
        //https://archiveservice:5000

        const res = await fetch('http://localhost/api/archive/efile', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + user.token
          },
          body: JSON.stringify({
            filename: filename.length > 0 ? filename : "IMG_"+new Date().getTime(),
            extension: file.type,
            enc_data: b64encoded
          }),
        });
        console.log('Response:', res);
        if (res.ok) {
          const body = await res.json();
          onSuccessfullUpload(body.Data as EncFile);
        }
      } catch (error) {
        console.log('Error:', error);
      }finally{
        handleClose();
      }
      
    }
    
    function handleClose() {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
      isOpen = false;
    }
  </script>
  
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-gray-900 border-2 border-gray-700 rounded-xl w-full max-w-2xl p-6">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-semibold text-gray-100">Upload new photo</h2>
        <button 
          class="bg-gray-800 border-2 text-gray-100 border-gray-700 rounded-lg px-4 py-2 hover:bg-gray-700"
          onclick={handleImport}
        >
          Import
        </button>
      </div>
  
      <input
        type="file"
        accept="image/*"
        class="hidden"
        bind:this={fileInput}
        onchange={handleFileSelect}
      />
      
      <div class="border-2 border-gray-700 rounded-xl aspect-video mb-6 overflow-hidden">
        {#if previewUrl}
          <img 
            src={previewUrl} 
            alt="Preview" 
            class="w-full h-full object-contain bg-gray-800"
          />
        {:else}
          <div class="w-full h-full flex items-center justify-center text-gray-400 bg-gray-800">
            Photo preview
          </div>
        {/if}
      </div>
      
      <div class="flex gap-4">
        <input
          type="text"
          placeholder="Filename"
          class="w-[50%] bg-gray-800 border-2 text-gray-100 border-gray-700 rounded-lg px-4 py-2 mb-6 focus:outline-none focus:border-gray-600"
          bind:value={filename}
        />

        <input
          type="password"
          placeholder="Password"
          class="w-[50%] bg-gray-800 border-2 text-gray-100 border-gray-700 rounded-lg px-4 py-2 mb-6 focus:outline-none focus:border-gray-600"
          bind:value={password}
        />
      </div>
      
      
      <div class="flex justify-end gap-4">
        <button 
          class="bg-gray-800 text-gray-100 border-2 border-gray-700 rounded-lg px-6 py-2 hover:bg-gray-700"
          onclick={handleClose}
        >
          Cancel
        </button>
        <button 
          class="bg-gray-800 border-2 text-gray-100 border-gray-700 rounded-lg px-6 py-2 hover:bg-gray-700"
          onclick={handleSave}
        >
          Save
        </button>
      </div>
    </div>
  </div>