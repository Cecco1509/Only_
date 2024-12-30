<script lang="ts">

	import { goto } from '$app/navigation';
	import PhotoDisplayModal from '$lib/components/PhotoDisplayModal.svelte';
	import PhotoUploadModal from '$lib/components/PhotoUploadModal.svelte';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
	import type { EncFile } from '$lib/EncFile/EncFile.js';
	import type { UserInformation } from '$lib/User/user';
	import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Eye, Minus, MinusIcon, PenSquareIcon, Plus, ScanEye, Share2, Trash2, UserRoundPen } from 'lucide-svelte';
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';

  const { data } = $props();
  const ownedFiles : EncFile[]  = data.ownedFiles;
  const user : UserInformation = data.user;
  const sharingUsers = data.sharingUsers;
  let users = $state<{username : string; id : number}[]>([]);
  console.log("OH ALLORA", ownedFiles, sharingUsers);

  let isNewFileModalOpen = $state<boolean>(false);
  let isDisplayPhotoModalOpen = $state<boolean>(false);
  let selectedPhoto = $state<EncFile | null >(null);

  let files = $state<{ [key : number] : EncFile[]}>({[user.id] : ownedFiles});
  let currentUser = $state<number>(user.id);

  let searchQuery = $state<string>("");
  let filterQuery = $state<string>("");
  

	function handleNewPhoto(event: MouseEvent & { currentTarget: EventTarget & HTMLButtonElement; }) {
		isNewFileModalOpen = true;
	}


	async function handleLogout(event: MouseEvent & { currentTarget: EventTarget & HTMLButtonElement; }) {
		
    try {
      await fetch('http://localhost/api/auth/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        },
      });
      document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
    } catch (error) {
      console.log(error);
    }finally{
      goto("/");
    }

	}

	function displayPhoto(file: EncFile): any {
		selectedPhoto = file;
    isDisplayPhotoModalOpen = true;
	}

  $effect(() => {
    console.log("searchQuery", searchQuery);

    if(searchQuery.length < 3) {
      users = [];
      return;
    } 

    fetch('http://localhost/api/auth/users?search='+searchQuery, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + user.token
      },
    }).then(response => response.json().then(data => {
      console.log(data);
      users = data.Data
    })).catch(error => {
      console.log(error);
    });
  });

  type SharedWith = {shared_with_user_id: number, shared_with_username: string, shared_by_user_id: number, shared_by_username: string};

  const sharedWith : { [key : string] : SharedWith[] } = $state<{ [key : string] : SharedWith[] }>({});
 
  const fetchShare = async (fileId: number) => {
    console.log("fetching share");

    if (sharedWith[fileId]) return;

    try {
      const response = await fetch('http://localhost/api/archive/efile/'+fileId+"/shared", {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      });

      if (!response.ok) {
        console.log(await response.text());
        sharedWith[fileId.toString()] = [];
        return;
      }

      const data = await response.json();
      console.log(data);
      sharedWith[fileId.toString()] = data.Data;
    } catch (error) {
      console.log(error);
    }
  }

  const handleShare = async (fileId: number, userId: number, username : string) => {
    console.log("sharing file", fileId, userId);

    try {
      const response = await fetch('http://localhost/api/archive/share/efile', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        },
        body: JSON.stringify({
          username : username,
          file_id : fileId
        })
      });
      const data = await response.json();
      console.log(data);
      if (response.ok)
        sharedWith[fileId.toString()].push({shared_with_username : username, shared_with_user_id : userId, shared_by_user_id : -1, shared_by_username : ''});
    } catch (error) {
      console.log(error);
    }
  }



	async function checkShared(user: { username: string; id: number; }, file: EncFile) : Promise<boolean> {
    if (!sharedWith[file.id.toString()]){
      await fetchShare(file.id);
    }
		return sharedWith[file.id.toString()].some((sharedUser) => sharedUser.shared_with_user_id == user.id);
	}


	async function deleteSharing(file: EncFile, userW : SharedWith) {
		try {
      const response = await fetch(`http://localhost/api/archive/share/efile/${file.id}?userid=${userW.shared_with_user_id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      });
      const data = await response.json();
      console.log(data);
      if (response.ok)
        sharedWith[file.id.toString()] = sharedWith[file.id.toString()].filter((shared) => shared.shared_with_user_id != userW.shared_with_user_id);
    } catch (error) {
      console.log(error);
    }
	}


	async function fetchSharedFiles(u: { shared_with_user_id: number; shared_with_username: string; shared_by_user_id: number; shared_by_username: string; }) {
		
    try {
      const res = await fetch('http://localhost/api/archive/shared/efile/'+u.shared_by_username, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      });

      if (!res.ok) {
        console.log(await res.text());
        return;
      }

      const data : { Data : EncFile[] } = await res.json();
      console.log(data.Data);
      files[u.shared_by_user_id] = data.Data;
      currentUser = u.shared_by_user_id;
    } catch (error) {
      console.log(error);
    }
	}


	async function handleFileDelete(file: EncFile) {
		try {
      const response = await fetch('http://localhost/api/archive/storage/'+file.filename, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + user.token
        }
      });

      if (!response.ok){
        console.log(await response.text());
        return;
      }

      const data = await response.json();
      console.log(data);
      files[currentUser] = files[currentUser].filter((f) => f.id != file.id); 
    } catch (error) {
      console.log(error);
    }
	}
</script>
  {#if isNewFileModalOpen}
    <PhotoUploadModal bind:isOpen={isNewFileModalOpen} onSuccessfullUpload={(file : EncFile) => files[user.id].push(file)}/>
  {/if}

  {#if isDisplayPhotoModalOpen && selectedPhoto}
    <PhotoDisplayModal bind:isOpen={isDisplayPhotoModalOpen} photo={selectedPhoto} />
  {/if}
  <div class="min-h-screen bg-gray-900 text-gray-100 p-4">
    <!-- Header -->
    <header class="flex items-center justify-between mb-8">
      <div class="w-32 h-12 border-2 border-gray-700 rounded-full flex items-center justify-center">
        Logo
      </div>
      
      <div class="flex-1 mx-8">
        <input 
          type="search"
          placeholder="Search file"
          class="w-full bg-gray-800 border-2 border-gray-700 rounded-full px-6 py-2 focus:outline-none focus:border-gray-600"
          bind:value={searchQuery}
        />
      </div>
      
      <div class="flex gap-4">
        <button class="hover:text-gray-300">search user</button>
        <button class="hover:text-gray-300">{data.user.username}</button>
        <button class="hover:text-gray-300" onclick={handleLogout}>logout</button>
      </div>
    </header>
  
    <div class="flex gap-8">
      <!-- Sidebar -->
      <aside class="w-64 h-[calc(100vh-8rem)] border-2 border-gray-700 rounded-xl p-4 z-[0]">
        <h2 class="text-xl mb-4">Shared with me</h2>
        <div class="w-full h-full">
          <Button onclick={() => currentUser = user.id} class="w-full border-grey-400 border flex justify-between align-center rounded-lg mb-2">
            Owned files
            {#if currentUser == user.id}
              <ScanEye/>
            {/if}
          </Button>
          <ScrollArea class="h-full w-full m-0">
            {#each sharingUsers as u}
              <Button onclick={async () => await fetchSharedFiles(u)} class="w-full border-grey-400 border flex justify-between align-center rounded-lg mb-2">
                {u.shared_by_username}
                {#if u.shared_by_user_id == currentUser}
                  <ScanEye/>
                {/if}
              </Button>
            {/each}
          </ScrollArea>
        </div>
      </aside>
  
      <!-- Main Content -->
      <main class="flex-1">
        <div class="flex gap-4 mb-6">
          <button 
            class="bg-gray-800 border-2 border-gray-700 rounded-lg px-6 py-2 hover:bg-gray-700"
            onclick={handleNewPhoto}
          >
            New photo
          </button>
          
          <input 
            type="text"
            placeholder="Filters"
            class="flex-1 bg-gray-800 border-2 border-gray-700 rounded-lg px-6 py-2 focus:outline-none focus:border-gray-600"
            bind:value={filterQuery}
          />
        </div>
  
        <!-- Files Table -->
        <div class="border-2 border-gray-700 rounded-xl p-4">
          <!-- Table Header -->
          <div class="grid grid-cols-4 gap-4 mb-4 px-4 text-gray-400">
            <span class="text-center">Name</span>
            <span class="text-center">Update date</span>
            <span class="text-center">Shared With</span>
            <span class="text-center">Actions</span>
          </div>
  
          <!-- File List -->
          <div class="grid grid-cols-4 gap-4 px-4 py-2 border-2 border-gray-800 rounded-lg">
          {#each files[currentUser] as file}
            
              <div class="flex items-center justify-center">{file.filename}.{file.extension}</div>
              <div class="flex items-center justify-center">{(new Date(file.creation_dt)).toDateString()}</div>
              <!-- svelte-ignore a11y_click_events_have_key_events -->
              <!-- svelte-ignore a11y_no_static_element_interactions -->
                <DropdownMenu.Root>
                    <DropdownMenu.Trigger on:click={() => fetchShare(file.id)}>
                      <Button disabled={currentUser != user.id } class="w-full h-full flex justify-center align-center bg-transparent hover:bg-transparent">
                        <UserRoundPen size=20/>
                      </Button>
                    </DropdownMenu.Trigger>
                  <DropdownMenu.Content>
                    <DropdownMenu.Group>
                      <!-- <DropdownMenu.GroupHeading>Shared with</DropdownMenu.GroupHeading>
                      <DropdownMenu.Separator /> -->
                      {#await fetchShare(file.id) then }
                        {#each sharedWith[file.id.toString()] as sharedUser}
                          <DropdownMenu.Item>
                            <div class="flex justify-between w-full items-center">
                              <div>{sharedUser.shared_with_username}</div>
                              <Button onclick={() => deleteSharing(file, sharedUser)} size="sm" class="bg-transparent hover:bg-transparent text-red-300 hover:text-red-600">
                                <MinusIcon class="flex-float-right"/>
                              </Button>
                            </div>
                          </DropdownMenu.Item>
                        {/each}
                      {/await}
                      
                      <!-- <DropdownMenu.Item>Profile</DropdownMenu.Item>
                      <DropdownMenu.Item>Billing</DropdownMenu.Item>
                      <DropdownMenu.Item>Team</DropdownMenu.Item>
                      <DropdownMenu.Item>Subscription</DropdownMenu.Item> -->
                    </DropdownMenu.Group>
                  </DropdownMenu.Content>
                </DropdownMenu.Root>

              <div class="flex gap-4 justify-center align-center">
                <Button variant="ghost" size="sm" onclick={async () => await displayPhoto(file)}>
                  <Eye/>
                </Button>
                <Dialog.Root>
                  <Button disabled={currentUser != user.id} class="bg-transparent hover:bg-transparent">
                    <Dialog.Trigger class={buttonVariants({ variant: "ghost", size: "sm"})}>
                      <Share2 color="orange"/>
                    </Dialog.Trigger>
                  </Button>
                  <Dialog.Content class="sm:max-w-[425px] bg-gray-900">
                    <Dialog.Header>
                      <Dialog.Title><span style="color: white;">Share photo</span></Dialog.Title>
                      <Dialog.Description>
                        Share this photo with whoever you like.
                      </Dialog.Description>
                    </Dialog.Header>
                    <div class="grid gap-4 py-4">
                      <Input id="name" placeholder="Inser username" bind:value={searchQuery} class="col-span-3"/>
                    </div>
                    <ScrollArea class="h-60 w-full rounded-lg border">
                      {#each users as user}
                        {#if user.id != file.user_id}
                          {#await checkShared(user, file) then isShared}
                            {#if !isShared}
                              <div class="border-grey-400 border flex justify-between align-center rounded-lg m-1">
                                <div class="p-1 h-max">
                                  <span class="text-lg text-white">{user.username}</span>
                                </div>
                                <Button variant="ghost" size="sm" onclick={async () => await handleShare(file.id, user.id, user.username)}>
                                  <Plus color="green"/>
                                </Button>
                              </div>
                            {/if}
                          {/await}
                        {/if}
                      {/each}
                    </ScrollArea>
                  </Dialog.Content>
                </Dialog.Root>
                <Button disabled={currentUser != user.id} variant="ghost" size="sm" onclick={async () => await handleFileDelete(file)}>
                  <Trash2 color="red"/>
                </Button>
              </div>
              
            
          {/each}
        </div>
        </div>
      </main>
    </div>
  </div>