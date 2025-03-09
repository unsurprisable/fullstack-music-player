<h1>Playlist Browser</h1>

<script lang="ts">
    import { goto } from "$app/navigation";
  import { apiFetch } from "$lib/api";
  import type { Playlist } from "$lib/types";
    import { onMount } from "svelte";
    import PlaylistDisplay from "./PlaylistDisplay.svelte";

  let nameInput: string = $state("");
  let nameError: string = $state("");

  let playlists = $state<Playlist[] | null>(null);

  function handleCreate() {
    if (nameInput.length == 0) {
      nameError = "Please enter a name for the playlist.";
      return;
    }
    createPlaylist(nameInput);
  }

  async function createPlaylist(name: string) {
    try {
      const data = await apiFetch("/playlists", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name }),
      });
      console.log(data);
      fetchAllPlaylists()
    } catch (error) {
      console.error(error);
    }
  }

  async function deletePlaylist(id: number) {
    try {
      const data = await apiFetch("/playlists/" + id, {
        method: 'DELETE'
      });
      console.log(data);
      fetchAllPlaylists()
    } catch (error) {
      console.error(error)
    }
  }

  async function fetchAllPlaylists() {
    try {
      const data: Playlist[] = await apiFetch<Playlist[]>("/playlists");
      console.log(data)
      playlists = data;
    } catch (error) {
      console.error(error)
    }
  }

  onMount(fetchAllPlaylists);
</script>

<input oninput={() => {nameError = ""}} bind:value={nameInput} />
<button onclick={handleCreate}>Create</button>

<div style="color: red;">
  {#if nameError.length != 0}
    <span>{nameError}</span>
  {/if}
</div>

{#if playlists}
  <div>
    {#each playlists.toReversed() as playlist}
      <PlaylistDisplay 
        {playlist} 
        editAction={() => { goto(`/playlist/${playlist.id}`) }} 
        deleteAction={() => { deletePlaylist(playlist.id) }} 
      />
    {/each}
  </div>
{/if}